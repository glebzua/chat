package requests

import "chatprjkt/internal/domain"

type UserRequest struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,alphanum,gte=6"`
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phoneNumber" `
	//Nickname    string `json:"nickname" validate:"required"`
}

type AuthRequest struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,gte=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,alphanum,gte=6"`
	NewPassword string `json:"newPassword" validate:"required,alphanum,gte=6"`
}

func (r UserRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Email:       r.Email,
		Name:        r.Name,
		PhoneNumber: r.PhoneNumber,
	}, nil
}

func (r RegisterRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Email:       r.Email,
		Password:    r.Password,
		Name:        r.Name,
		PhoneNumber: r.PhoneNumber,
	}, nil
}

func (r AuthRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Email:    r.Email,
		Password: r.Password,
	}, nil
}

func (r ChangePasswordRequest) ToDomainModel() (interface{}, error) {
	return domain.ChangePassword{
		OldPassword: r.OldPassword,
		NewPassword: r.NewPassword,
	}, nil
}
