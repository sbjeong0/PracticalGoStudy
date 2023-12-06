package main

import (
	"errors"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		config
		err error
	}{
		{
			config: config{
				numTimes: -1,
			},
			err: errors.New("must specify a number greater than 0"),
		},
		{
			config: config{},
			err:    errors.New("must specify a number greater than 0"),
		},
		{
			config: config{
				numTimes: 10,
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		err := validateArgs(tc.config)
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Errorf("expected error %v, got %v", tc.err, err)
		}
		if tc.err == nil && err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}
}
