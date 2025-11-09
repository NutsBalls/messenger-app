package domain

import "errors"

var (
	ErrMessageNotFound = errors.New("message not found")
	ErrChatNotFound    = errors.New("chat not found")
	ErrUserNotInChat   = errors.New("user is not in chat")
	ErrBadRequest      = errors.New("bad request")
)
