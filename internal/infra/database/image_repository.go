package database

import (
	"chatprjkt/internal/domain"
	"fmt"
	"time"

	"github.com/upper/db/v4"
)

const ImagesTableName = "images"

type image struct {
	Id          int64      `db:"id,omitempty"`
	SenderId    int64      `db:"senderid"`
	ObjId       *int64     `db:"objid"`
	ChatId      string     `db:"chatid"`
	ObjType     *string    `db:"obj_type"`
	Name        string     `db:"name"`
	RecipientId int64      `db:"recipientid"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type ImageRepository interface {
	Save(img domain.Image) (domain.Image, error)
	Update(img domain.Image) (domain.Image, error)
	FindById(id int64) (domain.Image, error)
	FindByName(name string) (domain.Image, error)
	FindAll(objId int64, objType string) ([]domain.Image, error)
	Delete(id int64) error
	DeleteAll(objId int64, objType string) error
}

type imageRepository struct {
	coll db.Collection
}

func NewImageRepository(dbSession db.Session) ImageRepository {
	return imageRepository{
		coll: dbSession.Collection(ImagesTableName),
	}
}

func (r imageRepository) Save(i domain.Image) (domain.Image, error) {
	img := image{}
	fmt.Println("imageRepository  i", i)
	img.FromDomainModel(i)
	img.CreatedDate = time.Now()
	img.UpdatedDate = img.CreatedDate

	err := r.coll.InsertReturning(&img)
	if err != nil {
		return domain.Image{}, err
	}

	return img.ToDomainModel(), nil

}

func (r imageRepository) Update(i domain.Image) (domain.Image, error) {
	img := image{}
	img.UpdatedDate = time.Now()
	img.FromDomainModel(i)

	err := r.coll.Find(db.Cond{"id": i.Id}).Update(&img)
	if err != nil {
		return domain.Image{}, err
	}

	return img.ToDomainModel(), nil
}

func (r imageRepository) FindById(id int64) (domain.Image, error) {
	var i image

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&i)
	if err != nil {
		return domain.Image{}, err
	}

	return i.ToDomainModel(), nil
}
func (r imageRepository) FindByName(name string) (domain.Image, error) {
	var i image

	err := r.coll.Find(db.Cond{"name": name, "deleted_date": nil}).One(&i)
	if err != nil {
		return domain.Image{}, err
	}

	return i.ToDomainModel(), nil
}

func (r imageRepository) FindAll(objId int64, objType string) ([]domain.Image, error) {
	var images []image

	err := r.coll.Find(db.Cond{"deleted_date": nil, "obj_id": objId, "obj_type": objType}).OrderBy("id").All(&images)
	if err != nil {
		return nil, err
	}

	return mapImagesToDomainCollection(images), nil
}

func (r imageRepository) Delete(id int64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r imageRepository) DeleteAll(objId int64, objType string) error {
	return r.coll.Find(db.Cond{"deleted_date": nil, "obj_id": objId, "obj_type": objType}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (i *image) FromDomainModel(img domain.Image) {
	dd := img.DeletedDate
	if dd == nil {
		dd = &time.Time{}
	}
	if img.ObjId != 0 {
		i.ObjId = &img.ObjId
		i.ObjType = &img.ObjType
	}
	i.Id = img.Id
	i.ChatId = img.ChatId
	i.SenderId = img.SenderId
	i.RecipientId = img.RecipientId
	i.Name = img.Name
	i.CreatedDate = img.CreatedDate
	i.UpdatedDate = img.UpdatedDate
	i.DeletedDate = dd
}

func (i image) ToDomainModel() domain.Image {
	var objId int64
	var objType string
	if i.ObjId != nil {
		objId = *i.ObjId
		objType = *i.ObjType
	}
	return domain.Image{
		Id:          i.Id,
		SenderId:    i.SenderId,
		ObjId:       objId,
		ChatId:      i.ChatId,
		ObjType:     objType,
		Name:        i.Name,
		RecipientId: i.RecipientId,
		CreatedDate: i.CreatedDate,
		UpdatedDate: i.UpdatedDate,
		DeletedDate: i.DeletedDate,
	}
}

func mapImagesToDomainCollection(images []image) []domain.Image {
	result := make([]domain.Image, 0, len(images))

	for _, i := range images {
		result = append(result, i.ToDomainModel())
	}

	return result
}
