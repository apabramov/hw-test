package main

import (
	"dp-git/github.com/stretchr/testify/require"
	"embed"
	_ "embed"
	"os"
	"testing"
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
	t.Run("simple case", func(t *testing.T) {

		tests := []test{
			{from: "./testdata/input.txt", cmp: "testdata/out_offset0_limit0.txt"},
			{from: "./testdata/input.txt", limit: 10, cmp: "testdata/out_offset0_limit10.txt"},
			{from: "./testdata/input.txt", limit: 1000, cmp: "testdata/out_offset0_limit1000.txt"},
			{from: "./testdata/input.txt", limit: 10000, cmp: "testdata/out_offset0_limit10000.txt"},
			{from: "./testdata/input.txt", limit: 1000, offset: 100, cmp: "testdata/out_offset100_limit1000.txt"},
			//{from: "./testdata/input.txt", limit: 1000, offset: 6000, cmp: "testdata/out_offset6000_limit1000.txt"},
		}

		for _, tc := range tests {
			f, _ := os.CreateTemp("", "tmp")
			Copy(tc.from, f.Name(), tc.offset, tc.limit)
			out, _ := f.Stat()
			cmp, _ := td.Open(tc.cmp)
			fc, _ := cmp.Stat()
			f.Close()
			cmp.Close()
			os.Remove(f.Name())
			require.Equal(t, fc.Size(), out.Size(), "OK")
		}

	})
}
