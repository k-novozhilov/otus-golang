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

func Copy(fromPath, toPath string, offset, limit int64) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = from.Close()
	}()

	info, err := from.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	if !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	fileSize := info.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = from.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	to, err := os.Create(toPath)
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
	for copied < copySize {
		toRead := bufSize
		if copySize-copied < bufSize {
			toRead = copySize - copied
		}
		n, err := from.Read(buf[:toRead])
		if n > 0 {
			wn, werr := to.Write(buf[:n])
			if werr != nil {
				return werr
			}
			if wn != n {
				return io.ErrShortWrite
			}
			copied += int64(n)
			if err := printProgress(copied, copySize); err != nil {
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
