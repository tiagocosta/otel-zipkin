package zipcodeapi

import "github.com/stretchr/testify/mock"

type ZipAPIMocked struct {
	mock.Mock
}

func (m *ZipAPIMocked) FindCity(zipCode string) (string, error) {
	args := m.Called(zipCode)
	return args.String(0), args.Error(1)
}
