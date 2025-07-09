package service

// TokenService is an interface for token generation
type TokenService interface {
	Generate(userID int64) (string, error)
	Validate(token string) (int64, error)
}
