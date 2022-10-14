package app

import (
	"chatprjkt/config"
	"chatprjkt/internal/domain"
	"chatprjkt/internal/infra/database"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService interface {
	Register(user domain.User) (domain.User, string, error)
	Login(user domain.User) (domain.User, string, error)
	ChangePassword(user domain.User, req domain.ChangePassword, sess domain.Session) error
	Logout(sess domain.Session) error
	Check(sess domain.Session) error
	GenerateJwt(user domain.User) (string, error)
}

type authService struct {
	authRepo    database.SessionRepository
	userService UserService
	config      config.Configuration
	tokenAuth   *jwtauth.JWTAuth
}

func NewAuthService(ar database.SessionRepository, us UserService, cf config.Configuration, ta *jwtauth.JWTAuth) AuthService {
	return authService{
		authRepo:    ar,
		userService: us,
		config:      cf,
		tokenAuth:   ta,
	}
}

func (s authService) Register(user domain.User) (domain.User, string, error) {
	_, err := s.userService.FindByEmail(user.Email)
	if err == nil {
		log.Printf("invalid credentials")
		return domain.User{}, "", errors.New("invalid credentials")
	} else if !errors.Is(err, db.ErrNoMoreRows) {
		log.Print(err)
		return domain.User{}, "", err
	}

	user.Activated = true
	user, err = s.userService.Save(user)
	if err != nil {
		log.Print(err)
		return domain.User{}, "", err
	}

	token, err := s.GenerateJwt(user)
	return user, token, err
}

func (s authService) Login(user domain.User) (domain.User, string, error) {
	u, err := s.userService.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			log.Printf("AuthService: failed to find user %s", err)
		}
		log.Printf("AuthService: login error %s", err)
		return domain.User{}, "", err
	}

	valid := s.checkPasswordHash(user.Password, u.Password)
	if !valid {
		return domain.User{}, "", errors.New("invalid credentials")
	}

	token, err := s.GenerateJwt(u)
	return u, token, err
}

func (s authService) Logout(sess domain.Session) error {
	return s.authRepo.Delete(sess)
}

func (s authService) ChangePassword(user domain.User, req domain.ChangePassword, sess domain.Session) error {
	if !s.checkPasswordHash(req.OldPassword, user.Password) {
		return errors.New("invalid credentials")
	}

	if s.checkPasswordHash(req.NewPassword, user.Password) {
		return errors.New("old password used")
	}

	user.Password = req.NewPassword

	_, err := s.userService.Update(user, domain.User{Password: req.NewPassword})
	if err != nil {
		return err
	}

	err = s.authRepo.Delete(sess)
	if err != nil {
		return err
	}

	return nil
}

func (s authService) GenerateJwt(user domain.User) (string, error) {
	sess := domain.Session{UserId: user.Id, UUID: uuid.New()}
	err := s.authRepo.Save(sess)
	if err != nil {
		log.Printf("AuthService: failed to save session %s", err)
		return "", err
	}

	claims := map[string]interface{}{
		"user_id": sess.UserId,
		"uuid":    sess.UUID,
	}
	jwtauth.SetExpiryIn(claims, s.config.JwtTTL)
	_, tokenString, err := s.tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s authService) Check(sess domain.Session) error {
	return s.authRepo.Exists(sess)
}

func (s authService) checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
