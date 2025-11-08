package service

import "messages/internal/messages/store"

type MessagesService struct {
	store *store.MessageRepository
}

func NewMessageService(store *store.MessageRepository) *MessagesService {
	return &MessagesService{
		store: store,
	}
}
