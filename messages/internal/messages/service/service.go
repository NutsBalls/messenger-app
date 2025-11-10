package service

import "messages/internal/messages/store"

type MessagesService struct {
	store *store.MessagesRepository
}

func (s *MessagesService) NewMessageService(store *store.MessagesRepository) *MessagesService {
	return &MessagesService{
		store: store,
	}
}
