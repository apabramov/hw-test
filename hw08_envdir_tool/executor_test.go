package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	e, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	code := RunCmd([]string{"/bin/bash", "arg=1"}, e)
	require.Equal(t, code, 0, "Exit code should  be 0")
}
