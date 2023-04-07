package testhelpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertEqualErrors(t *testing.T) {
	type args struct {
		want error
		got  error
	}
	tests := []struct {
		name string
		args args
		pass bool
	}{
		{
			name: "Should fail if expected is nil, but got an error",
			args: args{
				want: nil,
				got:  errors.New("some-unexpected-error"),
			},
			pass: false,
		},
		{
			name: "Should pass if expected and actual errors are nil",
			args: args{
				want: nil,
				got:  nil,
			},
			pass: true,
		},
		{
			name: "Should fail if expected is not nil, but actual error is nil",
			args: args{
				want: errors.New("some-expected-error"),
				got:  nil,
			},
			pass: false,
		},
		{
			name: "Should fail if expected and actual error are different",
			args: args{
				want: errors.New("some-expected-error"),
				got:  errors.New("some-actual-error"),
			},
			pass: false,
		},
		{
			name: "Should pass if expected and actual error are equal",
			args: args{
				want: errors.New("some-error"),
				got:  errors.New("some-error"),
			},
			pass: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &testing.T{}
			AssertEqualErrors(mockT, tt.args.want, tt.args.got)

			assert.Equal(t, tt.pass, !mockT.Failed())
		})
	}
}
