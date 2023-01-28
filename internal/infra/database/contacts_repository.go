package database

import (
	"chatprjkt/internal/domain"
	"fmt"
	"github.com/upper/db/v4"
	"time"
)

const ContactsTableName = "contacts"

type contacts struct {
	Id          int64      `db:"id,omitempty"`
	UserId      int64      `db:"userid,omitempty"`
	ContactId   int64      `db:"contactid,omitempty"`
	Activated   bool       `db:"activated"`
	ChatId      string     `db:"chatid"`
	Nickname    string     `db:"nickname"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type ContactsRepository interface {
	Save(domainItem domain.Contact) (domain.Contact, error)
	//FindAll() (domain.Contacts, error)
	//FindAll(pageSize, page uint) (domain.Contacts, error)
	FindAllForId(id int64) (domain.Contacts, error)
	FindById(id int64) (domain.Contact, error)
	Delete(id int64) error
	Activation(id int64, chatid string) error
}

type contactsRepository struct {
	coll db.Collection
}

func NewContactsRepository(dbSession db.Session) ContactsRepository {
	return contactsRepository{
		coll: dbSession.Collection(ContactsTableName),
	}
}
func (r contactsRepository) Save(domainContact domain.Contact) (domain.Contact, error) {
	s := r.mapDomainToModel(domainContact)
	s.CreatedDate = time.Now()
	s.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&s)
	if err != nil {
		return domain.Contact{}, fmt.Errorf("contacts repository Save: %w", err)
	}

	return r.mapModelToDomain(s), nil
}

//func (r contactsRepository) FindAll() (domain.Contacts, error) {
//	var contact []contacts
//
//	err := r.coll.Find().All(&contact)
//	if err != nil {
//		return domain.Contacts{}, err
//	}
//
//	return r.mapModelToDomainCollection(contact), nil
//}

//	func (r contactsRepository) FindAll(pageSize, page uint) (domain.Contacts, error) {
//		var contact []contacts
//
//		err := r.coll.Find().Paginate(pageSize).Page(page).All(&contact)
//		if err != nil {
//			return domain.Contacts{}, err
//		}
//
//		return r.mapModelToDomainCollection(contact), nil
//	}
func (r contactsRepository) FindAllForId(userId int64) (domain.Contacts, error) {
	var contact []contacts
	contactCond := db.Cond{"userid": userId}
	err := r.coll.Find(contactCond).All(&contact)
	if err != nil {
		return domain.Contacts{}, err
	}
	return r.mapModelToDomainCollection(contact), nil
}
func (r contactsRepository) Activation(id int64, chatId string) error {

	err := r.coll.Find(db.Cond{"id": id}).Update(map[string]interface{}{"activated": true, "chatid": chatId})
	return err
}
func (r contactsRepository) FindById(id int64) (domain.Contact, error) {
	var u contacts

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&u)
	if err != nil {
		return domain.Contact{}, err
	}

	return r.mapModelToDomain(u), nil
}

func (r contactsRepository) Delete(id int64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r contactsRepository) mapDomainToModel(d domain.Contact) contacts {
	return contacts{
		Id:          d.Id,
		UserId:      d.UserId,
		ContactId:   d.ContactId,
		Activated:   d.Activated,
		ChatId:      d.ChatId,
		Nickname:    d.Nickname,
		CreatedDate: d.CreatedDate,
	}
}

func (r contactsRepository) mapModelToDomain(m contacts) domain.Contact {
	return domain.Contact{
		Id:          m.Id,
		UserId:      m.UserId,
		ContactId:   m.ContactId,
		Activated:   m.Activated,
		ChatId:      m.ChatId,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
		Nickname:    m.Nickname,
	}
}

func (r contactsRepository) mapModelToDomainCollection(contact []contacts) domain.Contacts {
	var result []domain.Contact

	for _, u := range contact {
		result = append(result, r.mapModelToDomain(u))
	}

	res := domain.Contacts{
		Items: result,
	}

	return res
}
