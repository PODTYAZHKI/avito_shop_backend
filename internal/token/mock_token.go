package token

import "github.com/stretchr/testify/mock"

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) Generate(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}
