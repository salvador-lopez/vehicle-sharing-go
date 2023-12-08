package domain

import (
	"errors"
	"fmt"
	"regexp"
)

type VIN struct {
	number string
}

func NewVIN(number string) (*VIN, error) {
	err := guardVIN(number)
	if err != nil {
		return nil, err
	}
	return &VIN{number: number}, nil
}

var ErrInvalidVin = errors.New("invalid vin provided")

func guardVIN(number string) error {
	matches, _ := regexp.Match("^[A-HJ-NPR-Z\\d]{8}[\\dX][A-HJ-NPR-Z\\d]{8}$", []byte(number))
	if !matches {
		return fmt.Errorf("%v: %s", ErrInvalidVin, number)
	}

	return nil
}
