package resources

import (
	"chatprjkt/internal/domain"
)

type UserDto struct {
	Id          int64   `json:"id"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phoneNumber"`
	Avatar      *string `json:"avatar,omitempty"`
	Role        string  `json:"role"`
}
type UserInfoDto struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar,omitempty"`
	test   string
}

type UsersDto struct {
	Items []UserDto `json:"items"`
	Pages uint64    `json:"pages"`
	Total uint64    `json:"total"`
}

type AuthDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

func (d UserDto) DomainToDto(user domain.User) UserDto {
	return UserDto{
		Id:          user.Id,
		Email:       user.Email,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Avatar:      user.Avatar,
		Role:        string(user.Role),
	}
}
func (d UserInfoDto) DomainToDto(user domain.User) UserInfoDto {
	return UserInfoDto{
		Id:     user.Id,
		Name:   user.Name,
		Avatar: user.Avatar,
	}
}

func (d UserDto) DomainToDtoCollection(u domain.Users) []UserDto {
	var result []UserDto

	for _, c := range u.Items {
		var dto UserDto
		userDto := dto.DomainToDto(c)
		result = append(result, userDto)
	}

	return result
}
func (d UserInfoDto) DomainToDtoCollection(u domain.Users) []UserInfoDto {
	var result []UserInfoDto

	for _, c := range u.Items {
		var dto UserInfoDto
		userDto := dto.DomainToDto(c)
		result = append(result, userDto)
	}

	return result
}
func (d AuthDto) DomainToDto(token string, user domain.User) AuthDto {
	var userDto UserDto
	return AuthDto{
		Token: token,
		User:  userDto.DomainToDto(user),
	}
}
