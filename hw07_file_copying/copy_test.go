package main

import (
	"embed"
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
	t.Run("copy check size", func(t *testing.T) {
		tests := []test{
			{from: "./testdata/input.txt", cmp: "testdata/out_offset0_limit0.txt"},
			{from: "./testdata/input.txt", limit: 10, cmp: "testdata/out_offset0_limit10.txt"},
			{from: "./testdata/input.txt", limit: 1000, cmp: "testdata/out_offset0_limit1000.txt"},
			{from: "./testdata/input.txt", limit: 10000, cmp: "testdata/out_offset0_limit10000.txt"},
			{from: "./testdata/input.txt", limit: 1000, offset: 100, cmp: "testdata/out_offset100_limit1000.txt"},
			{from: "./testdata/input.txt", limit: 1000, offset: 6000, cmp: "testdata/out_offset6000_limit1000.txt"},
		}

		for _, tc := range tests {
			f, err := os.CreateTemp("", "tmp")
			if err != nil {
				t.Error(err)
			}
			err = Copy(tc.from, f.Name(), tc.offset, tc.limit)
			if err != nil {
				t.Error(err)
			}
			out, err := f.Stat()
			if err != nil {
				t.Error(err)
			}
			cmp, err := td.Open(tc.cmp)
			if err != nil {
				t.Error(err)
			}
			fc, err := cmp.Stat()
			if err != nil {
				t.Error(err)
			}
			b, err := os.ReadFile(f.Name())
			if err != nil {
				t.Error(err)
			}
			c, err := os.ReadFile(tc.cmp)
			if err != nil {
				t.Error(err)
			}

			f.Close()
			cmp.Close()
			os.Remove(f.Name())
			require.Equal(t, fc.Size(), out.Size(), "File size not equal")
			require.True(t, reflect.DeepEqual(b, c), "DeepEqual file not equal")
		}
	})
}
