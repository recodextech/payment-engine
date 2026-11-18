package billing

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"payment-engine/pkg/env"
)

// Config holds the bank account contractors pay plans and add-on packs into. Contractors transfer the
// price with their payment reference (prefix + account number, e.g. AA1111111) and an admin records it.
type Config struct {
	BankName          string `env:"BANK_NAME"`
	BankAccountName   string `env:"BANK_ACCOUNT_NAME"`
	BankAccountNumber string `env:"BANK_ACCOUNT_NUMBER"`
	BankBranch        string `env:"BANK_BRANCH"`
	// ReferencePrefix starts every contractor payment reference, e.g. AA1111111
	ReferencePrefix string `env:"PAYMENT_REFERENCE_PREFIX" envDefault:"AA"`
}

var prefixPattern = regexp.MustCompile(`^[A-Z]{1,6}$`)

func (c *Config) Register() error {
	return env.Parse(c)
}

func (c *Config) Validate() error {
	if !prefixPattern.MatchString(c.ReferencePrefix) {
		return fmt.Errorf("PAYMENT_REFERENCE_PREFIX must be 1-6 uppercase letters")
	}
	return nil
}

func (c *Config) Print() interface{} {
	defer log.Println("billing configs loaded")
	return *c
}

// Reference is the payment reference of an account number, e.g. AA1111111.
func (c *Config) Reference(accountNumber int64) string {
	return c.ReferencePrefix + strconv.FormatInt(accountNumber, 10)
}

// ParseReference returns the account number in a payment reference. Spaces, dashes and case are ignored.
func (c *Config) ParseReference(ref string) (int64, bool) {
	clean := make([]byte, 0, len(ref))
	for i := 0; i < len(ref); i++ {
		b := ref[i]
		switch {
		case b == ' ' || b == '-':
			continue
		case b >= 'a' && b <= 'z':
			b -= 'a' - 'A'
		}
		clean = append(clean, b)
	}
	s := string(clean)
	if len(s) <= len(c.ReferencePrefix) || s[:len(c.ReferencePrefix)] != c.ReferencePrefix {
		return 0, false
	}
	var no int64
	for _, d := range s[len(c.ReferencePrefix):] {
		if d < '0' || d > '9' || no > 1e15 {
			return 0, false
		}
		no = no*10 + int64(d-'0')
	}
	return no, no > 0
}
