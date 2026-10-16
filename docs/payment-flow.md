# Payments, Transactions and Money Movement

This document explains how the payment engine records a payment for a job,
how the money is held while the job runs, and how wallet balances change when
the job completes or is cancelled.

## Core concepts

| Concept | Table | What it is |
|---|---|---|
| **Wallet** | `fix.account.wallets` | A balance owned by an account. Each account has at most one wallet per type: `CASH`, `POINTS` or `CARD`. The `balance` column is the source of truth for the current balance, and `sequence_number` in the payload goes up by one each time the balance changes. |
| **Payment** | `fix.payments` | One money movement tied to a job. It holds the total `amount`, a `status` and, in its payload, the list of cr/dr **transaction entries** it will post. |
| **Transaction** | `fix.payment.transactions` | One ledger line on one wallet: a **debit** (negative amount) or a **credit** (positive amount). Transactions are written only when a payment succeeds, and they are never updated or deleted. |

### Double entry

Every payment is balanced: the payer wallet is **debited** the full amount and
the payee wallets are **credited** shares that add up to the same amount.

```
sum(credits) == payment amount == -(debit)
```

`PaymentService.CreatePayment` rejects a payment that doesn't add up (within a
float tolerance of `0.000001`). It also rejects a payment where:

- the amount or any entry amount is `<= 0`,
- `job_id` or `payer_id` is missing,
- a payee wallet is the payer wallet,
- a payee wallet's type differs from the payer wallet's type (money never
  moves between `CASH`, `POINTS` and `CARD`),
- an entry has neither `payee_id` nor `payee_account_id`, or the payee
  account has no wallet of the payer's type.

A payee can be given as a wallet (`payee_id`), or as an account
(`payee_account_id`), in which case that account's wallet of the payer's type
is used.

## Payment lifecycle

```
                  POST /jobs/{jobId}/payment/create
                                 │
                                 ▼
                         ┌──────────────┐   amount is HELD on the payer wallet
                         │ IN_PROGRESS  │   balances unchanged
                         └──────┬───────┘
          PATCH …/cancelled     │      PATCH …/success
       ┌────────────────────────┴────────────────────────┐
       ▼                                                 ▼
┌──────────────┐                                 ┌──────────────┐
│  CANCELLED   │  hold released                  │   SUCCESS    │  ledger rows written
│              │  no ledger rows, no balance     │              │  balances updated
│              │  change                         │              │  hold released
└──────────────┘                                 └──────────────┘
```

`CANCELLED` and `SUCCESS` are final: a payment that isn't `IN_PROGRESS`
can't change state again.

### 1. Create: the hold

`POST /jobs/{jobId}/payment/create` with a body like:

```json
{
  "amount": 1000,
  "type": "PAYMENT",
  "payer_id": "<contractor CASH wallet>",
  "reason_code": "JOB_COMPLETED",
  "transaction_entries": [
    { "amount": 900, "payee_account_id": "<worker account>" },
    { "amount": 100, "payee_id": "<platform CASH wallet>" }
  ],
  "payment_information": { "method": "CASH", "reference": {} }
}
```

1. The service validates the request and resolves every wallet (see
   [Double entry](#double-entry)). It builds the entry list: one debit of
   `-amount` on the payer, then one credit per payee.
2. `PaymentRepository.CreatePayment` opens **one DB transaction**
   (`InTransaction`), and inside it:
   - locks every wallet in the entries with `SELECT … FOR UPDATE`, in wallet
     ID order,
   - copies each locked wallet's current balance and sequence number into its
     entry, as a snapshot of the wallet when the payment was placed,
   - inserts the payment with status `IN_PROGRESS`.
3. The transaction commits and the locks are released.

**No balance changes at this point.** The `IN_PROGRESS` payment *is* the hold:
its debit entries are money that is promised but not yet moved. There is no
separate hold record to keep in sync.

### 2. Cancel: release the hold

`PATCH /jobs/{jobId}/payment/{paymentId}/cancelled`

In one DB transaction:

1. Load the job. A payment can be cancelled only while the job is `ACCEPTED`
   or `STARTED`, and never after the job has `ENDED`.
2. Lock the payment row (`FOR UPDATE`) and check that it is `IN_PROGRESS`.
3. Set the status to `CANCELLED`.

That's all. No transactions are written and no balance changes, so the hold
disappears because the payment is no longer `IN_PROGRESS`.

### 3. Success: move the money

`PATCH /jobs/{jobId}/payment/{paymentId}/success`

In one DB transaction:

1. Lock the payment row (`FOR UPDATE`) and check that it is `IN_PROGRESS`.
2. Set the status to `SUCCESS`.
3. Lock every wallet in the entries (`FOR UPDATE`, wallet ID order) and read
   their current balances.
4. For each entry, in order:
   - `balance_after = wallet.balance + entry.amount` (debits are negative),
   - insert a row into `fix.payment.transactions` with the amount and
     `payee.current_balance = balance_after`,
   - update the wallet's `balance` to `balance_after` and add one to its
     `sequence_number`.
5. Commit. Every ledger row and balance change is applied together, or none
   of them are.

## Worked example

The contractor's CASH wallet has a balance of 5,000, and the worker and
platform wallets have 0. A 1,000 payment is created for a job.

| Step | Contractor | Worker | Platform | Held (contractor) | Available (contractor) |
|---|---|---|---|---|---|
| Before | 5,000 | 0 | 0 | 0 | 5,000 |
| Create (`IN_PROGRESS`) | 5,000 | 0 | 0 | 1,000 | 4,000 |
| **Success** | 4,000 | 900 | 100 | 0 | 4,000 |
| *or* **Cancel** | 5,000 | 0 | 0 | 0 | 5,000 |

On success, three transaction rows are written:

| Wallet | Amount | Balance after |
|---|---|---|
| Contractor | -1,000 | 4,000 |
| Worker | +900 | 900 |
| Platform | +100 | 100 |

## Balances shown on a wallet

`GET /wallets` returns these for each wallet:

| Field | Source |
|---|---|
| `balance` | `fix.account.wallets.balance`: money that has actually moved |
| `totalCredited` | sum of positive transaction rows on the wallet |
| `totalDebited` | sum of negative transaction rows on the wallet, as a positive number |
| `heldAmount` | sum of the debit entries of the wallet's `IN_PROGRESS` payments |
| `availableBalance` | `balance - heldAmount` |

`availableBalance` **can be negative.** Creating a payment doesn't check
funds, so a contractor can commit to more than their balance and the wallet
shows the shortfall.

`GET /wallets/{walletId}/transactions` returns the latest 100 ledger rows of
a wallet, newest first, each with its job ID and the balance after it.

## Concurrency

Two jobs can create or complete payments on the same wallet at the same
moment. Without locking, both would read the same balance and one update would
overwrite the other. The engine prevents this with row locks:

- **Wallet rows** are locked (`FOR UPDATE`) during create and success. A
  second request touching the same wallet waits until the first transaction
  commits, then reads the updated balance. Requests on different wallets
  don't wait for each other.
- **Payment rows** are locked during success and cancel, so the same payment
  can't be completed twice, or completed and cancelled at the same time.
- **Lock order** is always the payment row first, then wallets sorted by ID.
  Two requests can't each hold a lock the other is waiting for, so they can't
  deadlock.

These locks last only for one DB transaction, usually a few milliseconds. The
hold that lasts for the whole job is the `IN_PROGRESS` status, which locks
nothing, so other jobs on the same wallet keep working while it is in place.

Things to keep in mind:

- Postgres waits for a row lock indefinitely by default. Setting
  `lock_timeout` makes a stuck request fail fast instead of making others wait.
- A wallet that takes part in nearly every payment, such as a platform fee
  wallet, makes those payments run one at a time.

## Code map

| Concern | Where |
|---|---|
| HTTP routes | `internal/http/routes.go` |
| Request shape and validation tags | `internal/http/request/payment.go` |
| Validation, wallet resolution, building the cr/dr entries | `internal/services/payment_service.go` |
| Create (hold), cancel, success (posting), lock helpers | `internal/adaptors/repositories/payment_repository.go` |
| Wallet reads, row locks, balance update, totals, held amount | `internal/adaptors/repositories/account_wallet_repository.go` |
| Statuses, wallet types, table names | `internal/domain/definitions.go`, `internal/domain/tables.go` |
