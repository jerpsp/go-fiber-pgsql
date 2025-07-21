package storage

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

type S3Repository interface {
	UploadPublicFile(objFile *multipart.FileHeader) (string, error)
	DeletePublicFile(objKey string) error
	GetURLFile(objKey string) (string, error)
}

type s3Repository struct {
	s3 Storage
}

func NewS3Repo(s3 Storage) S3Repository {
	return &s3Repository{s3: s3}
}

func (s *s3Repository) UploadPublicFile(objFile *multipart.FileHeader) (string, error) {
	file, err := objFile.Open()
	if err != nil {
		return "", err
	}

	nameList := strings.Split(objFile.Filename, ".")

	// generate uniq key for file name in object storage
	key := fmt.Sprintf("%s%s%s", uuid.New().String(), ".", nameList[len(nameList)-1])

	_, err = s.s3.PutPublicFile(file, key)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *s3Repository) DeletePublicFile(objKey string) error {
	_, err := s.s3.DeletePublicFile(objKey)

	if err != nil {
		return err
	}

	return nil
}

func (s *s3Repository) GetURLFile(objKey string) (string, error) {
	url, err := s.s3.GetPresignURL(objKey)

	if err != nil {
		return "", err
	}

	return url, nil
}
