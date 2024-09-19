// Package readercomp provides comparison functions for readers or more specific for comparing file contents
package readercomp

import (
	"bytes"
	"io"
	"os"
)

// Equal compares two io.Reader instances efficiently
func Equal(r1, r2 io.Reader, bufSize int) (bool, error) {
	if bufSize <= 0 {
		bufSize = 4096 // Default buffer size if not specified or invalid
	}
	buf1 := make([]byte, bufSize)
	buf2 := make([]byte, bufSize)

	for {
		n1, err1 := readFull(r1, buf1)
		n2, err2 := readFull(r2, buf2)

		// Check for non-EOF errors first
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}

		// Compare the read data
		if n1 != n2 || !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}

		// Handle EOF cases
		if err1 == io.EOF && err2 == io.EOF {
			return true, nil // Both readers reached EOF simultaneously
		}
		if err1 == io.EOF || err2 == io.EOF {
			return false, nil // One reader reached EOF before the other
		}
	}
}

// readFull reads from r into buf, handling short reads
func readFull(r io.Reader, buf []byte) (int, error) {
	n, err := io.ReadFull(r, buf)
	if err == io.ErrUnexpectedEOF {
		return n, io.EOF
	}
	return n, err
}

// FilesEqual tests two files for equality (same content)
func FilesEqual(name1, name2 string) (bool, error) {
	f1Info, err := os.Stat(name1)
	if err != nil {
		return false, err
	}
	f2Info, err := os.Stat(name2)
	if err != nil {
		return false, err
	}

	if f1Info.Size() != f2Info.Size() {
		return false, nil
	}

	f1, err := os.Open(name1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(name2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	// 64 KiB buffers seem to be most performant with larger files
	return Equal(f1, f2, 64*1024)
}
