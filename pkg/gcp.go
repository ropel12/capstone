package pkg

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/education-hub/BE/errorr"
)

type StorageGCP struct {
	ClG        *storage.Client
	ProjectID  string
	BucketName string
	Path       string
}

func (s *StorageGCP) UploadFile(file multipart.File, fileName string) error {
	if !strings.Contains(strings.ToLower(fileName), ".jpg") && !strings.Contains(strings.ToLower(fileName), ".png") && !strings.Contains(strings.ToLower(fileName), ".jpeg") && !strings.Contains(strings.ToLower(fileName), ".pdf") {
		fmt.Println(strings.Contains(strings.ToLower(fileName), ".jpg"))
		return errorr.NewBad("File type not allowed")
	}
	return nil
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	wc := s.ClG.Bucket(s.BucketName).Object(s.Path + fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return errorr.NewInternal(err.Error())
	}

	if err := wc.Close(); err != nil {
		return errorr.NewInternal(err.Error())
	}
	return nil
}
