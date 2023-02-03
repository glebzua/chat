package app

import (
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/database"
	"chatprjkt/internal/infra/http/resources"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Save(user domain.User) (domain.User, error)
	FindAll() (domain.Users, error)
	//FindAll(pageSize, page uint) (domain.Users, error)
	Find(id int64) (interface{}, error)
	FindById(id int64) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	Update(user domain.User, req domain.User) (domain.User, error)
	Delete(id int64, sess domain.Session) error
	ForcedDelete(id int64) error
}

type userService struct {
	userRepo      database.UserRepository
	pusherService PusherService
}

func NewUserService(ur database.UserRepository, ps PusherService) UserService {
	return userService{
		userRepo:      ur,
		pusherService: ps,
	}
}

func (s userService) Save(user domain.User) (domain.User, error) {
	var err error

	user.Password, err = s.generatePasswordHash(user.Password)
	if err != nil {
		log.Printf("UserService generatePasswordHash: %s", err)
		return domain.User{}, err
	}
	user.Activated = true
	user.Role = domain.ROLE_USER

	u, err := s.userRepo.Save(user)
	if err != nil {
		log.Printf("UserService Save: %s", err)
		return domain.User{}, err
	}
	var userDto resources.UserDto
	s.pusherService.NewUser(userDto.DomainToDto(user))
	return u, err
}

func (s userService) FindAll() (domain.Users, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		log.Printf("UserService FindAll: %s", err)
		return domain.Users{}, err
	}

	return users, nil
}

//	func (s userService) FindAll(pageSize, page uint) (domain.Users, error) {
//		users, err := s.userRepo.FindAll(pageSize, page)
//		if err != nil {
//			log.Printf("UserService FindAll: %s", err)
//			return domain.Users{}, err
//		}
//
//		return users, nil
//	}
func (s userService) Find(id int64) (interface{}, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		log.Printf("UserService Find: %s", err)
		return nil, err
	}

	return user, nil
}

func (s userService) FindById(id int64) (domain.User, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		log.Printf("UserService FindById: %s", err)
		return domain.User{}, err
	}

	return user, err
}

func (s userService) FindByEmail(email string) (domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("UserService FindByEmail: %s", err)
		return domain.User{}, err
	}

	return user, err
}

func (s userService) Update(user domain.User, req domain.User) (domain.User, error) {
	var err error
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	user, err = s.userRepo.Update(user)
	if err != nil {
		log.Printf("UserService Update: %s", err)
		return domain.User{}, err
	}

	return user, err
}

func (s userService) Delete(id int64, sess domain.Session) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		log.Printf("UserService: %s", err)
		return err
	}

	return nil
}

func (s userService) ForcedDelete(id int64) error {
	err := s.userRepo.ForcedDelete(id)
	if err != nil {
		log.Printf("UserService: %s", err)
		return err
	}

	return nil
}

func (s userService) generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
