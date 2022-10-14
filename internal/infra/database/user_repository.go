package database

import (
	"chatprjkt/internal/domain"
	"github.com/upper/db/v4"
	"time"
)

const UsersTableName = "users"

type user struct {
	Id          int64      `db:"id,omitempty"`
	Email       string     `db:"email"`
	Password    string     `db:"password,omitempty"`
	Name        string     `db:"name"`
	PhoneNumber string     `db:"phone_number"`
	Avatar      *string    `db:"avatar,omitempty"`
	Activated   bool       `db:"activated"`
	Role        string     `db:"role"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type UserRepository interface {
	Save(user domain.User) (domain.User, error)
	FindAll(pageSize, page uint) (domain.Users, error)
	FindById(id int64) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	Update(user domain.User) (domain.User, error)
	Delete(id int64) error
	ForcedDelete(id int64) error
}

type userRepository struct {
	coll db.Collection
}

func NewUserRepository(dbSession db.Session) UserRepository {
	return userRepository{
		coll: dbSession.Collection(UsersTableName),
	}
}

func (r userRepository) Save(user domain.User) (domain.User, error) {
	u := r.mapDomainToModel(user)

	u.CreatedDate = time.Now()
	u.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&u)
	if err != nil {
		return domain.User{}, err
	}

	return r.mapModelToDomain(u), nil
}

func (r userRepository) FindAll(pageSize, page uint) (domain.Users, error) {
	var users []user

	err := r.coll.Find().Paginate(pageSize).Page(page).All(&users)
	if err != nil {
		return domain.Users{}, err
	}

	return r.mapModelToDomainCollection(users), nil
}

func (r userRepository) FindById(id int64) (domain.User, error) {
	var u user

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&u)
	if err != nil {
		return domain.User{}, err
	}

	return r.mapModelToDomain(u), nil
}

func (r userRepository) FindByEmail(email string) (domain.User, error) {
	var u user

	err := r.coll.Find(db.Cond{"email": email, "deleted_date": nil}).One(&u)
	if err != nil {
		return domain.User{}, err
	}

	return r.mapModelToDomain(u), nil
}

func (r userRepository) Update(user domain.User) (domain.User, error) {
	u := r.mapDomainToModel(user)

	u.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": u.Id}).Update(&u)
	if err != nil {
		return domain.User{}, err
	}

	return r.mapModelToDomain(u), nil
}

func (r userRepository) Delete(id int64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r userRepository) ForcedDelete(id int64) error {
	return r.coll.Find(db.Cond{"id": id}).Delete()
}

func (r userRepository) mapDomainToModel(d domain.User) user {
	return user{
		Id:          d.Id,
		Email:       d.Email,
		Password:    d.Password,
		Name:        d.Name,
		PhoneNumber: d.PhoneNumber,
		Avatar:      d.Avatar,
		Activated:   d.Activated,
		Role:        string(d.Role),
	}
}

func (r userRepository) mapModelToDomain(m user) domain.User {
	return domain.User{
		Id:          m.Id,
		Email:       m.Email,
		Name:        m.Name,
		Password:    m.Password,
		PhoneNumber: m.PhoneNumber,
		Avatar:      m.Avatar,
		Activated:   m.Activated,
		Role:        domain.Role(m.Role),
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r userRepository) mapModelToDomainCollection(users []user) domain.Users {
	var result []domain.User

	for _, u := range users {
		result = append(result, r.mapModelToDomain(u))
	}

	res := domain.Users{
		Items: result,
	}

	return res
}
