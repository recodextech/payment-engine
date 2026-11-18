package domain

const (
	AccountTable       string = `"fix.account"`
	PaymentTable       string = `"fix.payments"`
	TransactionTable   string = `"fix.payment.transactions"`
	PlatformFeeTable   string = `"fix.platform.fees"`
	AccountWalletTable string = `"fix.account.wallets"`
	JobTable           string = `"fix.job"`
	ContractorTable    string = `"fix.contractor"`

	SubscriptionPlanTable    string = `"fix.configuration.subscription.plan"`
	SubscriptionTable        string = `"fix.subscriptions"`
	SubscriptionPaymentTable string = `"fix.subscription.payments"`
	QuotaLedgerTable         string = `"fix.subscription.quota.ledger"`
)
