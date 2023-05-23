package pkg

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
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

func (s *StorageGCP) GetPdf(filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*25)
	defer cancel()
	rc, err := s.ClG.Bucket(s.BucketName).Object(filename).NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
