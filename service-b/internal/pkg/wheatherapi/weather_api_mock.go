package wheatherapi

import "github.com/stretchr/testify/mock"

type WeatherAPIMocked struct {
	mock.Mock
}

func (m *WeatherAPIMocked) GetWeather(city string, apiKey string) (float64, error) {
	args := m.Called(city, apiKey)
	return args.Get(0).(float64), args.Error(1)
}
