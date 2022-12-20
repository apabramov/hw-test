package main

import (
	"embed"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	from   string
	offset int64
	limit  int64
	cmp    string
}

//go:embed testdata/*
var td embed.FS

func TestCopy(t *testing.T) {
	tests := []test{
		{from: "./testdata/input.txt", cmp: "testdata/out_offset0_limit0.txt"},
		{from: "./testdata/input.txt", limit: 10, cmp: "testdata/out_offset0_limit10.txt"},
		{from: "./testdata/input.txt", limit: 1000, cmp: "testdata/out_offset0_limit1000.txt"},
		{from: "./testdata/input.txt", limit: 10000, cmp: "testdata/out_offset0_limit10000.txt"},
		{from: "./testdata/input.txt", limit: 1000, offset: 100, cmp: "testdata/out_offset100_limit1000.txt"},
		{from: "./testdata/input.txt", limit: 1000, offset: 6000, cmp: "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("check %s", tc.cmp), func(t *testing.T) {
			f, err := os.CreateTemp("", "tmp")
			require.NoError(t, err)
			defer f.Close()

			err = Copy(tc.from, f.Name(), tc.offset, tc.limit)
			require.NoError(t, err)

			require.FileExists(t, f.Name())
			out, err := f.Stat()
			require.NoError(t, err)

			cmp, err := td.Open(tc.cmp)
			require.NoError(t, err)
			defer cmp.Close()

			fc, err := cmp.Stat()
			require.NoError(t, err)

			b, err := os.ReadFile(f.Name())
			require.NoError(t, err)

			c, err := os.ReadFile(tc.cmp)
			require.NoError(t, err)

			os.Remove(f.Name())
			require.Equal(t, fc.Size(), out.Size(), "File size not equal")
			require.True(t, reflect.DeepEqual(b, c), "DeepEqual file not equal")
		})
	}
}
