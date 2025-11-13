package service

import "messages/internal/messages/store"

type MessagesService struct {
	store store.Repository
}

func NewMessagesService(repo store.Repository) *MessagesService {
	return &MessagesService{
		store: repo,
	}
}
