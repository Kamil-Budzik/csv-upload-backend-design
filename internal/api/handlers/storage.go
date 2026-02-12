package handlers

import (
	"context"
	"io"
)

type FileStorage interface {
	UploadCSV(ctx context.Context, fileName string, fileSize int64, file io.Reader) (string, error)
}
