package filesystem

import (
	"log"
	"os"
	"path"
)

type ImageStorageService interface {
	SaveImage(filename string, content []byte) error
}

type imageStorageService struct {
	loc string
}

func NewImageStorageService(location string) ImageStorageService {
	return imageStorageService{
		loc: location,
	}
}

func (s imageStorageService) SaveImage(filename string, content []byte) error {
	location := path.Join(s.loc, filename)
	err := writeFileToStorage(location, content)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func writeFileToStorage(location string, file []byte) error {
	dirLocation := path.Dir(location)
	err := os.MkdirAll(dirLocation, os.ModePerm)
	if err != nil {
		log.Print(err)
		return err
	}

	err = os.WriteFile(location, file, os.ModePerm)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
