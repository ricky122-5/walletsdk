package service

import (
	"fmt"
	"regexp"
	"strings"
)

var addressPattern = regexp.MustCompile(`^(0x)?[0-9a-fA-F]{40}$`)

func ValidateTransaction(tx *Transaction) error {
	if tx == nil {
		return fmt.Errorf("%w: missing transaction", ErrValidation)
	}
	if !addressPattern.MatchString(strings.TrimSpace(tx.To)) {
		return fmt.Errorf("%w: invalid to address", ErrValidation)
	}
	if strings.TrimSpace(tx.Value) == "" {
		return fmt.Errorf("%w: value required", ErrValidation)
	}
	if tx.GasLimit == 0 {
		return fmt.Errorf("%w: gasLimit required", ErrValidation)
	}
	if strings.TrimSpace(tx.GasPrice) == "" {
		return fmt.Errorf("%w: gasPrice required", ErrValidation)
	}
	if tx.Nonce == 0 {
		return fmt.Errorf("%w: nonce required", ErrValidation)
	}
	return nil
}
