package service

const (
	ErrInvalidCredentials = "invalid credentials"
	ErrUserNotFound       = "user not found"
	ErrGenerateTokens     = "failed to generate tokens"
	ErrHashToken          = "failed to hash token"
	ErrUpdateToken        = "failed to update refresh token"
	ErrInvalidToken       = "invalid token"
	ErrTypeToken          = "wrong token type"
	ErrParseToken         = "invalid token after parse"
)
