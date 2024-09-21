package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/yaninyzwitty/messaging-service/models"
)

// MessagesRepository defines the interface for message-related operations.
type MessagesRepository interface {
	CreateMessage(ctx context.Context, message models.Message) (models.Message, error)
	UpdateMessage(ctx context.Context, messageId gocql.UUID, message models.Message) (models.Message, error)
	DeleteMessage(ctx context.Context, messageId gocql.UUID) error
	GetMessages(ctx context.Context) ([]models.Message, error)
	GetMessage(ctx context.Context, id gocql.UUID) (models.Message, error)
}

// messagesRepository is the concrete implementation of MessagesRepository.
type messagesRepository struct {
	session *gocqlx.Session
}

// NewMessagesRepository creates a new instance of messagesRepository.
func NewMessagesRepository(session *gocqlx.Session) MessagesRepository {
	return &messagesRepository{session: session}
}

// CreateMessage inserts a new message into the database.
func (r *messagesRepository) CreateMessage(ctx context.Context, message models.Message) (models.Message, error) {

	q := r.session.Query(models.MessageTable.Insert()).BindStruct(message)

	// query := qb.Insert(models.MessageTable.Name()).
	// 	Columns(models.MessageTable.Metadata().Columns...).
	// 	Query(*r.session)
	// if err := query.BindStruct(message).ExecRelease(); err != nil {
	// 	return models.Message{}, err
	// }

	if err := q.ExecRelease(); err != nil {
		return models.Message{}, err
	}

	return message, nil
}

// UpdateMessage updates an existing message in the database.
func (r *messagesRepository) UpdateMessage(ctx context.Context, id gocql.UUID, message models.Message) (models.Message, error) {

	query := qb.Update(models.MessageTable.Name()).
		Set("conversation_id", "sender_id", "body", "updated_at", "is_soft_deleted").
		Where(qb.Eq("id")).
		Query(*r.session)

	err := query.BindStruct(message).ExecRelease()
	if err != nil {
		return models.Message{}, err
	}
	return message, nil

}

func (r *messagesRepository) DeleteMessage(ctx context.Context, id gocql.UUID) error {
	query := qb.Delete(models.MessageTable.Name()).Where(qb.Eq("id")).Query(*r.session)
	err := query.BindMap(qb.M{"id": id}).ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

// GetMessages retrieves all messages from the database.
func (r *messagesRepository) GetMessages(ctx context.Context) ([]models.Message, error) {
	var messages []models.Message

	query := qb.Select(models.MessageTable.Name()).
		Columns("id", "conversation_id", "sender_id", "created_at", "updated_at", "body", "is_soft_deleted").
		Query(*r.session)

	iter := query.Iter()
	defer iter.Close()

	var message models.Message
	for iter.Scan(&message.ID, &message.ConversationID, &message.SenderId, &message.CreatedAt, &message.UpdatedAt, &message.Body, &message.IsSoftDeleted) {
		messages = append(messages, message)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil

}

// GetMessage retrieves a single message by its ID from the database.
func (r *messagesRepository) GetMessage(ctx context.Context, id gocql.UUID) (models.Message, error) {
	query := qb.Select(models.MessageTable.Name()).
		Columns("id", "conversation_id", "sender_id", "created_at", "updated_at", "body", "is_soft_deleted").
		Where(qb.Eq("id")).
		Query(*r.session)

	var message models.Message
	if err := query.BindMap(qb.M{"id": id}).GetRelease(&message); err != nil {
		return models.Message{}, err
	}
	return message, nil

}
