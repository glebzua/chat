package requests

import "chatprjkt/internal/domain"

type ContactRequest struct {
	ContactId int64 `json:"contactId"`
}

//type ContactCreateRequest struct {
//	ContactId int64 `json:"contactId" validate:"required,contactId"`
//}

//
//type RegisterRequest struct {
//	Email       string `json:"email" validate:"required,email"`
//	Password    string `json:"password" validate:"required,alphanum,gte=6"`
//	Name        string `json:"name" validate:"required"`
//	PhoneNumber string `json:"phoneNumber" validate:"required"`
//}

func (r ContactRequest) ToDomainModel() (interface{}, error) {
	return domain.Contact{
		ContactId: r.ContactId,
	}, nil
}

//func (r ContactCreateRequest) ToDomainModel() (interface{}, error) {
//	return domain.Contact{
//		ContactId: r.ContactId,
//	}, nil
//}

//func (r RegisterRequest) ToDomainModel() (interface{}, error) {
//	return domain.User{
//		Email:       r.Email,
//		Password:    r.Password,
//		Name:        r.Name,
//		PhoneNumber: r.PhoneNumber,
//	}, nil
//}
//
//func (r AuthRequest) ToDomainModel() (interface{}, error) {
//	return domain.User{
//		Email:    r.Email,
//		Password: r.Password,
//	}, nil
//}
//
//func (r ChangePasswordRequest) ToDomainModel() (interface{}, error) {
//	return domain.ChangePassword{
//		OldPassword: r.OldPassword,
//		NewPassword: r.NewPassword,
//	}, nil
//}
