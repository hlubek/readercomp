package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
)

const sizeMb = 256
const n = 10

func main() {
	const sizeB = sizeMb * 1 << 20

	// Construct a bunch of larger files that are pessimal for comparison (size is equal, content is equal up to the last byte)
	for i := 0; i < n; i++ {
		err := prepareFile(filename(sizeMb, i), sizeB)
		checkErr(err)
	}
}

func filename(sizeMb int, idx int) string {
	return fmt.Sprintf("file-%d-%002d", sizeMb, idx)
}

func prepareFile(filename string, sizeB int64) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.CopyN(f, nullReader{}, sizeB-1)
	if err != nil {
		return err
	}

	b := byte(rand.Int() % 256)
	x, err := f.Write([]byte{b})
	if err != nil {
		return err
	}
	fmt.Println(x, b)

	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type nullReader struct{}

func (r nullReader) Read(buf []byte) (n int, err error) {
	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}
	return len(buf), nil
}
