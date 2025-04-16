package main

import (
	"errors"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		offset  int64
		limit   int64
		want    string
		wantErr error
	}{
		{"all file", "testdata/input.txt", "out.txt", 0, 0, "testdata/out_offset0_limit0.txt", nil},
		{"limit 10", "testdata/input.txt", "out.txt", 0, 10, "testdata/out_offset0_limit10.txt", nil},
		{"limit 1000", "testdata/input.txt", "out.txt", 0, 1000, "testdata/out_offset0_limit1000.txt", nil},
		{"limit 10000", "testdata/input.txt", "out.txt", 0, 10000, "testdata/out_offset0_limit10000.txt", nil},
		{"offset 100 limit 1000", "testdata/input.txt", "out.txt", 100, 1000, "testdata/out_offset100_limit1000.txt", nil},
		{"offset 6000 limit 1000", "testdata/input.txt", "out.txt", 6000, 1000, "testdata/out_offset6000_limit1000.txt", nil},
		{"offset > file", "testdata/input.txt", "out.txt", 100000, 10, "", ErrOffsetExceedsFileSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Remove(tt.to)
			err := Copy(tt.from, tt.to, tt.offset, tt.limit)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ожидалась ошибка %v, получено %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			wantData, _ := os.ReadFile(tt.want)
			gotData, _ := os.ReadFile(tt.to)
			if string(wantData) != string(gotData) {
				t.Fatalf("файлы не совпадают")
			}
			os.Remove(tt.to)
		})
	}
}
