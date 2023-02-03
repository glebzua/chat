package app

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/database"
	"chatprjkt/internal/infra/filesystem"
	"log"
)

type ImageService interface {
	Save(image domain.Image, content []byte) (domain.Image, error)
	FindById(id int64) (domain.Image, error)
	FindByName(name string) (domain.Image, error)
	FindAll(objId int64, objType string) ([]domain.Image, error)
	Delete(id int64) error
	DeleteAll(objId int64, objType string) error
	Find(int64) (interface{}, error)
	Sync(objId int64, objType string, objImages, newImages []domain.Image) error
}

type imageService struct {
	repo            database.ImageRepository
	filesys         filesystem.ImageStorageService
	messagesService MessagesService
}

func NewImageService(r database.ImageRepository, s filesystem.ImageStorageService, ms MessagesService) ImageService {
	return &imageService{
		repo:            r,
		filesys:         s,
		messagesService: ms,
	}
}

func (s imageService) Save(image domain.Image, content []byte) (domain.Image, error) {
	err := s.filesys.SaveImage(image.Name, content)
	if err != nil {
		log.Print(err)
		return domain.Image{}, err
	}

	var img domain.Image
	img, err = s.repo.Save(image)
	if err != nil {
		log.Print(err)
		return domain.Image{}, err
	}
	msg := domain.Message{}

	msg.SenderId = img.SenderId
	msg.RecipientId = img.RecipientId
	msg.ChatId = img.ChatId
	msg.FileLoc = img.Name
	s.messagesService.Save(msg)
	return img, nil
}

func (s imageService) FindAll(objId int64, objType string) ([]domain.Image, error) {
	return s.repo.FindAll(objId, objType)
}

func (s imageService) Sync(objId int64, objType string, objImages, newImages []domain.Image) error {
	imgsNew := make(map[int64]struct{})
	for _, newImage := range newImages {
		isNewImage := true
		imgsNew[newImage.Id] = struct{}{}
		for _, objImage := range objImages {
			if objImage.Id == newImage.Id {
				isNewImage = false
				break
			}
		}
		if isNewImage {
			i, err := s.repo.FindById(newImage.Id)
			if err != nil {
				return err
			}
			i.ObjId = objId
			i.ObjType = objType
			_, err = s.repo.Update(i)
			if err != nil {
				return err
			}
		}
	}
	for _, objImage := range objImages {
		if _, exist := imgsNew[objImage.Id]; !exist {
			err := s.repo.Delete(objImage.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s imageService) FindById(id int64) (domain.Image, error) {
	return s.repo.FindById(id)
}

func (s imageService) FindByName(name string) (domain.Image, error) {
	return s.repo.FindByName(name)
}

func (s imageService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s imageService) DeleteAll(objId int64, objType string) error {
	return s.repo.DeleteAll(objId, objType)
}

func (s imageService) Find(id int64) (interface{}, error) {
	return s.repo.FindById(id)
}
