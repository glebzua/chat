package database

import (
	"chatprjkt/internal/domain"
	"fmt"
	"github.com/upper/db/v4"
	"time"
)

const MessagesTableName = "messages"

type messages struct {
	Id          int64      `db:"id,omitempty"`
	ChatId      string     `db:"chatid,"`
	SenderId    int64      `db:"senderid,"`
	RecipientId int64      `db:"recipientid"`
	Message     string     `db:"message"`
	Sended      bool       `db:"sended"`
	Received    bool       `db:"received"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type MessagesRepository interface {
	Save(domainItem domain.Message) (domain.Message, error)
	FindAll(pageSize, page uint) (domain.Messages, error)
	FindAllForId(id string) (domain.Messages, error)
	FindAllMessagesInChat(id int64, chatId string) (domain.Messages, error)
	FindById(id int64) (domain.Message, error)
	Delete(id int64) error
}

type messagesRepository struct {
	coll db.Collection
}

func NewMessagesRepository(dbSession db.Session) MessagesRepository {
	return messagesRepository{
		coll: dbSession.Collection(MessagesTableName),
	}
}
func (r messagesRepository) Save(domainMessages domain.Message) (domain.Message, error) {
	s := r.mapDomainToModel(domainMessages)
	s.CreatedDate = time.Now()
	s.UpdatedDate = time.Now()
	s.Sended = true
	err := r.coll.InsertReturning(&s)
	if err != nil {
		return domain.Message{}, fmt.Errorf("Messages repository Save: %w", err)
	}

	return r.mapModelToDomain(s), nil
}
func (r messagesRepository) FindAll(pageSize, page uint) (domain.Messages, error) {
	var message []messages

	err := r.coll.Find().Paginate(pageSize).Page(page).All(&message)
	if err != nil {
		return domain.Messages{}, err
	}

	return r.mapModelToDomainCollection(message), nil
}
func (r messagesRepository) FindAllForId(chatId string) (domain.Messages, error) {
	var message []messages
	messageCond := db.Cond{"chatid": chatId}
	err := r.coll.Find(messageCond).All(&message)
	if err != nil {
		return domain.Messages{}, err
	}
	return r.mapModelToDomainCollection(message), nil
}

func (r messagesRepository) FindAllMessagesInChat(_ int64, chatId string) (domain.Messages, error) {
	var message []messages
	messageCond := db.Cond{"chatid": chatId}
	err := r.coll.Find(messageCond).All(&message)
	if err != nil {
		return domain.Messages{}, err
	}

	return r.mapModelToDomainCollection(message), nil
}
func (r messagesRepository) FindById(id int64) (domain.Message, error) {
	var u messages

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&u)
	if err != nil {
		return domain.Message{}, err
	}

	return r.mapModelToDomain(u), nil
}

func (r messagesRepository) Delete(id int64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r messagesRepository) mapDomainToModel(d domain.Message) messages {
	return messages{
		Id:          d.Id,
		ChatId:      d.ChatId,
		SenderId:    d.SenderId,
		RecipientId: d.RecipientId,
		Message:     d.Message,
		Sended:      d.Sended,
		Received:    d.Received,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r messagesRepository) mapModelToDomain(m messages) domain.Message {
	return domain.Message{
		Id:          m.Id,
		ChatId:      m.ChatId,
		SenderId:    m.SenderId,
		RecipientId: m.RecipientId,
		Message:     m.Message,
		Sended:      m.Sended,
		Received:    m.Received,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r messagesRepository) mapModelToDomainCollection(message []messages) domain.Messages {
	var result []domain.Message

	for _, u := range message {
		result = append(result, r.mapModelToDomain(u))
	}

	res := domain.Messages{
		Items: result,
	}

	return res
}
