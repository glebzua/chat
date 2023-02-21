package requests

import "chatprjkt/internal/domain"

type ContactRequest struct {
	ContactId int64  `json:"contactId"`
	Nickname  string `json:"nickname"`
}

func (r ContactRequest) ToDomainModel() (interface{}, error) {
	return domain.Contact{
		ContactId: r.ContactId,
		Nickname:  r.Nickname,
	}, nil
}
