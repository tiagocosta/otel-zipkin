package usecase

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetWeatherUseCase_Execute(t *testing.T) {
	mock := new(GetWeatherUseCaseMocked)
	mock.On("Execute", GetWeatherInputDTO{
		ZipCode: "71218010",
	}).Return(GetWeatherOutputDTO{
		Celsius:    "25.60",
		Fahrenheit: "78.08",
		Kelvin:     "298.60",
	}, nil)
	mock.On("Execute", GetWeatherInputDTO{
		ZipCode: "12345678",
	}).Return(GetWeatherOutputDTO{}, errors.New("can not find zipcode"))

	type args struct {
		input GetWeatherInputDTO
	}
	tests := []struct {
		name    string
		uc      *GetWeatherUseCase
		args    args
		want    GetWeatherOutputDTO
		wantErr bool
	}{
		{
			"Valid inputDTO. Should return all temperatures in outputDTO",
			&GetWeatherUseCase{},
			args{
				GetWeatherInputDTO{
					ZipCode: "71218010",
				},
			},
			GetWeatherOutputDTO{
				Celsius:    "25.60",
				Fahrenheit: "78.08",
				Kelvin:     "298.60",
			},
			false,
		},
		{
			"Non existing zipcode. Should return an error",
			&GetWeatherUseCase{},
			args{
				GetWeatherInputDTO{
					ZipCode: "12345678",
				},
			},
			GetWeatherOutputDTO{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mock.Execute(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeatherUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWeatherUseCase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
