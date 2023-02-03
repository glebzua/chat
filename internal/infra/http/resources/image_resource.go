package resources

import (
	"chatprjkt/internal/domain"
	"fmt"
)

type ImageDto struct {
	Id          int64  `json:"id"`
	Link        string `json:"link"`
	SenderId    int64  `json:"SenderId"`
	RecipientId int64  `json:"RecipientId"`
	ChatId      string `json:"ChatId"`
}

type ImagesResource struct {
	images []domain.Image
}

func NewImageResource(images []domain.Image) ImagesResource {
	return ImagesResource{images: images}
}

func (ir ImagesResource) Serialize() []ImageDto {
	result := make([]ImageDto, 0, len(ir.images))

	for _, i := range ir.images {
		var dto ImageDto
		result = append(result, dto.DomainToDto(i))
	}

	return result
}

func (d *ImageDto) DomainToDto(i domain.Image) ImageDto {
	link := fmt.Sprintf("/static/%s", i.Name)
	d.Id = i.Id
	d.Link = link
	d.SenderId = i.SenderId
	d.RecipientId = i.RecipientId
	d.ChatId = i.ChatId
	return *d
}
