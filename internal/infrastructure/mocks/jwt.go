package mocks

type JWTManagerMock struct {
	GenerateFn func(userID string) (string, error)
	ValidateFn func(token string) (string, error)
}

func (m *JWTManagerMock) Generate(userID string) (string, error) {
	return m.GenerateFn(userID)
}

func (m *JWTManagerMock) Validate(token string) (string, error) {
	return m.ValidateFn(token)
}
