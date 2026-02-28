package utils

import (
	"errors"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	ErrInvalidSyntax      = errors.New("invalid email syntax")
	ErrNoMXRecords        = errors.New("domain has no MX records")
	ErrMailboxNotExist    = errors.New("mailbox does not exist")
)

// Create verifier once (reuse it)
var verifier = emailverifier.NewVerifier()
// Only enable SMTP if you REALLY need strict validation
// var verifier = emailverifier.NewVerifier().EnableSMTPCheck()

func IsEmailValid(email string) error {
	ret, err := verifier.Verify(email)
	if err != nil {
		return err
	}

	if !ret.Syntax.Valid {
		return ErrInvalidSyntax
	}

	if !ret.HasMxRecords {
		return ErrNoMXRecords
	}

	// If SMTP check enabled:
	if ret.Reachable == "no" {
		return ErrMailboxNotExist
	}

	// Accept "yes" AND "unknown"
	return nil
}