package service

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/yaninyzwitty/messaging-service/models"
	"github.com/yaninyzwitty/messaging-service/repository"
)

type MessagesService interface {
	CreateMessage(ctx context.Context, message models.Message) (models.Message, error)
	GetMessages(ctx context.Context) ([]models.Message, error)
	GetMessage(ctx context.Context, messageId gocql.UUID) (models.Message, error)
	DeleteMessage(ctx context.Context, messageId gocql.UUID) error
	UpdateMessage(ctx context.Context, messageId gocql.UUID, message models.Message) (models.Message, error)
	GetMessagesByPagingState(ctx context.Context, pageSize int, pagingState []byte) ([]models.Message, []byte, error)
}

type messageService struct {
	repo repository.MessagesRepository
}

func NewMessagesService(repo repository.MessagesRepository) MessagesService {
	return &messageService{repo: repo}
}

func (s *messageService) CreateMessage(ctx context.Context, message models.Message) (models.Message, error) {
	return s.repo.CreateMessage(ctx, message)
}

func (s *messageService) GetMessages(ctx context.Context) ([]models.Message, error) {
	return s.repo.GetMessages(ctx)
}

func (s *messageService) GetMessage(ctx context.Context, messageId gocql.UUID) (models.Message, error) {
	return s.repo.GetMessage(ctx, messageId)
}

func (s *messageService) DeleteMessage(ctx context.Context, messageId gocql.UUID) error {
	return s.repo.DeleteMessage(ctx, messageId)
}

func (s *messageService) UpdateMessage(ctx context.Context, messageId gocql.UUID, message models.Message) (models.Message, error) {
	return s.repo.UpdateMessage(ctx, messageId, message)
}

func (s *messageService) GetMessagesByPagingState(ctx context.Context, pageSize int, pagingState []byte) ([]models.Message, []byte, error) {
	return s.repo.GetMessagesByPagingState(ctx, pageSize, pagingState)
}
