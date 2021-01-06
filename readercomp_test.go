package readercomp

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type readResult struct {
	data []byte
	err  error
}

var errTest = errors.New("test IO error")

func TestEqual(t *testing.T) {
	tests := []struct {
		name    string
		r1      []readResult
		r2      []readResult
		bufSize int
		want    bool
		wantErr bool
	}{
		{
			name:    "single read, eq",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 100), io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 100), io.EOF}},
			bufSize: 8192,
			want:    true,
			wantErr: false,
		},
		{
			name:    "single read, r1 err",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 100), errTest}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 100), io.EOF}},
			bufSize: 8192,
			want:    false,
			wantErr: true,
		},
		{
			name:    "single read, r2 err",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 100), io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 100), errTest}},
			bufSize: 8192,
			want:    false,
			wantErr: true,
		},
		{
			name:    "single read, same size, neq",
			r1:      []readResult{{bytes.Repeat([]byte{'x'}, 100), io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 100), io.EOF}},
			bufSize: 8192,
			want:    false,
			wantErr: false,
		},
		{
			name:    "single read, len(r1) < len(r2), neq",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 80), io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 100), io.EOF}},
			bufSize: 8192,
			want:    false,
			wantErr: false,
		},
		{
			name:    "single read, len(r1) < len(r2), neq",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 100), io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 80), io.EOF}},
			bufSize: 8192,
			want:    false,
			wantErr: false,
		},
		{
			name:    "multi read, eq",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 128), io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 128), io.EOF}},
			bufSize: 64,
			want:    true,
			wantErr: false,
		},
		{
			name:    "multi read, r1 with late EOF, eq",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 128), nil}, {nil, io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 128), io.EOF}},
			bufSize: 64,
			want:    true,
			wantErr: false,
		},
		{
			name:    "multi read, r2 with late EOF, eq",
			r1:      []readResult{{bytes.Repeat([]byte{'v'}, 128), io.EOF}},
			r2:      []readResult{{bytes.Repeat([]byte{'v'}, 128), nil}, {nil, io.EOF}},
			bufSize: 64,
			want:    true,
			wantErr: false,
		},
		{
			name: "multi read, diff short, eq",
			r1: []readResult{
				{bytes.Repeat([]byte{'v'}, 45), nil},
				{bytes.Repeat([]byte{'v'}, 100), nil},
				{bytes.Repeat([]byte{'v'}, 51), io.EOF},
			},
			r2: []readResult{
				{bytes.Repeat([]byte{'v'}, 31), nil},
				{bytes.Repeat([]byte{'v'}, 27), nil},
				{bytes.Repeat([]byte{'v'}, 138), nil},
				{nil, io.EOF},
			},
			bufSize: 64,
			want:    true,
			wantErr: false,
		},
		{
			name: "multi read, diff short, len(r1) < len(r2), neq",
			r1: []readResult{
				{bytes.Repeat([]byte{'v'}, 45), nil},
				{bytes.Repeat([]byte{'v'}, 100), nil},
				{bytes.Repeat([]byte{'v'}, 51), io.EOF},
			},
			r2: []readResult{
				{bytes.Repeat([]byte{'v'}, 31), nil},
				{bytes.Repeat([]byte{'v'}, 27), nil},
				{bytes.Repeat([]byte{'v'}, 139), nil},
				{nil, io.EOF},
			},
			bufSize: 64,
			want:    false,
			wantErr: false,
		},
		{
			name: "multi read, diff short, r1 err",
			r1: []readResult{
				{bytes.Repeat([]byte{'v'}, 45), nil},
				{bytes.Repeat([]byte{'v'}, 100), nil},
				{bytes.Repeat([]byte{'v'}, 51), errTest},
			},
			r2: []readResult{
				{bytes.Repeat([]byte{'v'}, 31), nil},
				{bytes.Repeat([]byte{'v'}, 27), nil},
				{bytes.Repeat([]byte{'v'}, 138), nil},
				{nil, io.EOF},
			},
			bufSize: 64,
			want:    false,
			wantErr: true,
		},
		{
			name: "multi read, diff short, r2 err",
			r1: []readResult{
				{bytes.Repeat([]byte{'v'}, 45), nil},
				{bytes.Repeat([]byte{'v'}, 100), nil},
				{bytes.Repeat([]byte{'v'}, 51), io.EOF},
			},
			r2: []readResult{
				{bytes.Repeat([]byte{'v'}, 31), nil},
				{bytes.Repeat([]byte{'v'}, 27), nil},
				{bytes.Repeat([]byte{'v'}, 138), nil},
				{nil, errTest},
			},
			bufSize: 64,
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr1 := fixtureReader(tt.r1)
			fr2 := fixtureReader(tt.r2)
			got, err := Equal(&fr1, &fr2, tt.bufSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("Equal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Equal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilesEqual(t *testing.T) {
	result, err := FilesEqual("go.mod", "go.mod")
	if err != nil {
		t.Fatalf("FilesEqual() returned error: %v", err)
	}
	if !result {
		t.Fatalf("FilesEqual() expected to return true")
	}

	result, err = FilesEqual("go.mod", "readercomp.go")
	if err != nil {
		t.Fatalf("FilesEqual() returned error: %v", err)
	}
	if result {
		t.Fatalf("FilesEqual() expected to return false")
	}
}

type fixtureReader []readResult

func (f *fixtureReader) Read(buf []byte) (n int, err error) {
	if f == nil || len(*f) == 0 {
		return 0, io.EOF
	}
	res := (*f)[0]
	m := copy(buf, res.data)
	if len(res.data) == m {
		// Fully copied result data, update results slice for next result
		*f = (*f)[1:]
		return m, res.err
	}

	// Partially copied, update current result slice, do not yet return any error!
	(*f)[0].data = (*f)[0].data[m:]
	return m, nil
}
