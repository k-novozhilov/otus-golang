package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func printProgress(done, total int64) error {
	percent := float64(done) / float64(total) * 100
	if total == 0 {
		percent = 100
	}
	if _, err := os.Stdout.WriteString("\rProgress: "); err != nil {
		return err
	}
	if _, err := os.Stdout.WriteString(fmt.Sprintf("%6.2f%%", percent)); err != nil {
		return err
	}
	return nil
}

func openSourceFile(fromPath string, offset int64) (*os.File, int64, error) {
	from, err := os.Open(fromPath)
	if err != nil {
		return nil, 0, err
	}

	info, err := from.Stat()
	if err != nil {
		return from, 0, ErrUnsupportedFile
	}
	if !info.Mode().IsRegular() {
		return from, 0, ErrUnsupportedFile
	}

	fileSize := info.Size()
	if offset > fileSize {
		return from, 0, ErrOffsetExceedsFileSize
	}

	_, err = from.Seek(offset, io.SeekStart)
	if err != nil {
		return from, 0, err
	}

	return from, fileSize, nil
}

func createDestFile(toPath string) (*os.File, error) {
	to, err := os.Create(toPath)
	if err != nil {
		return nil, err
	}
	return to, nil
}

func writeData(to *os.File, copied *int64, n int, buf []byte) error {
	wn, werr := to.Write(buf[:n])
	if werr != nil {
		return werr
	}
	if wn != n {
		return io.ErrShortWrite
	}
	*copied += int64(n)
	return nil
}

func updateProgress(copied, copySize int64) error {
	if copied == copySize || (copySize > 100 && copied%(copySize/100) == 0) || copied%(1024*1024) == 0 {
		if err := printProgress(copied, copySize); err != nil {
			return err
		}
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	from, fileSize, err := openSourceFile(fromPath, offset)
	if err != nil {
		return err
	}
	defer func() {
		_ = from.Close()
	}()

	to, err := createDestFile(toPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = to.Close()
	}()

	copySize := fileSize - offset
	if limit > 0 && limit < copySize {
		copySize = limit
	}

	bufSize := int64(32 * 1024)
	buf := make([]byte, bufSize)
	var copied int64

	if err := printProgress(0, copySize); err != nil {
		return err
	}

	for copied < copySize {
		toRead := bufSize
		if copySize-copied < bufSize {
			toRead = copySize - copied
		}
		n, err := from.Read(buf[:toRead])

		if n > 0 {
			if err := writeData(to, &copied, n, buf); err != nil {
				return err
			}

			if err := updateProgress(copied, copySize); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	if _, err := os.Stdout.WriteString("\n"); err != nil {
		return err
	}
	return nil
}
