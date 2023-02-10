package filesystem

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type ImageStorageService interface {
	SaveImage(filename string, content []byte) error
	ReadFileFromStorage(filePath string) (string, error)
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

func (s imageStorageService) ReadFileFromStorage(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile("file_storage/" + filePath)
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	// Print the full base64 representation of the image
	return base64Encoding, err
}
func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
