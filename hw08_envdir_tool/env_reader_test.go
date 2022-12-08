package main

import (
	"dp-git/github.com/stretchr/testify/require"
	"errors"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {

	t.Run("ENV set unset", func(t *testing.T) {

		e, err := ReadDir("./testdata/env")
		require.NoError(t, err)

		for i, v := range e {
			if val, ok := os.LookupEnv(i); ok {
				require.Equal(t, val, v.Value, "Env ")
			} else {
				if v.Value != "" {
					require.Error(t, errors.New("Env should deleted"))
				}
			}
		}
	})
}
