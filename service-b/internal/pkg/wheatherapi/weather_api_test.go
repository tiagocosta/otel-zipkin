package wheatherapi

import (
	"errors"
	"testing"
)

func TestGetWeather(t *testing.T) {
	mock := new(WeatherAPIMocked)
	mock.On("GetWeather", "Brasília", "123456").Return(29.0, nil)
	mock.On("GetWeather", "non existing city", "123456").Return(0.0, errors.New("can not find zipcode"))
	type args struct {
		city   string
		apiKey string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			"Valid City. Should return the correct temperature",
			args{
				city:   "Brasília",
				apiKey: "123456",
			},
			29.0,
			false,
		},
		{
			"Invalid City. Should return error",
			args{
				city:   "non existing city",
				apiKey: "123456",
			},
			0.0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mock.GetWeather(tt.args.city, tt.args.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetWeather() = %v, want %v", got, tt.want)
			}
		})
	}
}
