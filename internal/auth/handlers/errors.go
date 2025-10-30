package handlers

type ErrResponse struct {
	Error   error
	Message string
}

type MessageResponse struct {
	Message string
}

func ReturnError(err error, message string) ErrResponse {
	return ErrResponse{
		Error:   err,
		Message: message,
	}
}

func ReturnMessage(message string) MessageResponse {
	return MessageResponse{
		Message: message,
	}
}
