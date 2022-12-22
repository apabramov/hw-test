package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Check struct {
		Ms []string `validate:"len:1"`
		Mi []int    `validate:"max:1"`
	}
)

func TestValidate(t *testing.T) {
	var ve ValidationErrors
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          &User{ID: "123", Name: "Alex", Age: 3, meta: nil},
			expectedErr: ValidationErrors{ValidationError{Field: "Age", Err: ErrMin}},
		},
		{
			in:          &User{ID: "1", Name: "Alex", Age: 18},
			expectedErr: ve,
		},
		{
			in:          &User{ID: "123", Email: "mail@google"},
			expectedErr: ValidationErrors{ValidationError{Field: "Email", Err: ErrRegExp}},
		},
		{
			in:          &App{Version: "1.2.3.456"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: ErrLen}},
		},
		{
			in:          &Response{Code: 307},
			expectedErr: ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}},
		},
		{
			in:          &Token{Header: []byte{0, 1, 2}},
			expectedErr: ve,
		},
		{
			in:          &Check{Ms: []string{"1", "22"}},
			expectedErr: ValidationErrors{ValidationError{Field: "Ms", Err: ErrLen}},
		},
		{
			in:          &Check{Mi: []int{1, 22}},
			expectedErr: ValidationErrors{ValidationError{Field: "Mi", Err: ErrMax}},
		},
		{
			in:          &[]int{1, 2},
			expectedErr: ErrStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			if err := Validate(tt.in); err != nil {
				require.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
