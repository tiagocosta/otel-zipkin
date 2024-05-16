package usecase

import "github.com/stretchr/testify/mock"

type GetWeatherUseCaseMocked struct {
	mock.Mock
}

func (m *GetWeatherUseCaseMocked) Execute(input GetWeatherInputDTO) (GetWeatherOutputDTO, error) {
	args := m.Called(input)
	return args.Get(0).(GetWeatherOutputDTO), args.Error(1)
}
