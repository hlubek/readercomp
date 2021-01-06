package main

import (
	"testing"

	"github.com/hlubek/readercomp"
)

func BenchmarkFilesEqual(b *testing.B) {
	for m := 0; m < b.N; m++ {
		for i := 0; i < n; i += 2 {
			j := i + 1
			res, err := readercomp.FilesEqual(filename(sizeMb, i), filename(sizeMb, j))
			if err != nil {
				b.Fatal(err)
			}
			b.Log(i, j, res)
			b.SetBytes(sizeMb * (2 << 20))
		}
	}
}
