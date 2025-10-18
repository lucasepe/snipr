package io

import (
	"bytes"
	"errors"
	"io"
	"os"
)

var ErrNoInputDetected = errors.New("no input detected")

// FileOrStdin returns an io.Reader based on the given filename.
// If the filename is non-empty, it attempts to open the file.
// If the filename is empty, it checks if data is available on stdin.
// The returned cleanup function should be deferred by the caller
// to close the file if necessary.
// If no valid input source is found, it returns an error.
func FileOrStdin(filename string) (io.Reader, func(), error) {
	var (
		src     *os.File
		cleanup func() = func() {} // no-op di default
		err     error
	)

	if filename != "" {
		src, err = os.Open(filename)
		if err != nil {
			return nil, cleanup, err
		}
		cleanup = func() {
			src.Close()
		}

		return src, cleanup, nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, cleanup, err
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		src = os.Stdin
	}

	if src == nil {
		return nil, cleanup, ErrNoInputDetected
	}

	return src, cleanup, nil
}

func ReaderToString(r io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	return buf.String(), err
}
