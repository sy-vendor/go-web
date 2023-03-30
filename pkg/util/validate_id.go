package util

import (
	"errors"
)

func ValidateID(id uint64) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}
	return nil
}
