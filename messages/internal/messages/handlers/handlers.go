package handlers

import "messages/internal/messages/service"

type MessagesHandlers struct {
	service *service.MessagesService
}

func (h *MessagesHandlers) NewMessagesService(service *service.MessagesService) *MessagesHandlers {
	return &MessagesHandlers{
		service: service,
	}
}
