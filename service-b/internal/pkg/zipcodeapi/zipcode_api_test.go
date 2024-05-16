package zipcodeapi

import (
	"errors"
	"testing"
)

func TestFindCity(t *testing.T) {
	mock := new(ZipAPIMocked)
	mock.On("FindCity", "1234567").Return("", errors.New("invalid zipcode"))
	mock.On("FindCity", "12345678").Return("", nil)
	mock.On("FindCity", "71218010").Return("Brasília", nil)

	type args struct {
		zipCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Invalid zipcode (length less than 8). Should return an error",
			args{
				zipCode: "1234567",
			},
			"",
			true,
		},
		{
			"Valid but non existing zipcode. Should return an empty city",
			args{
				zipCode: "12345678",
			},
			"",
			false,
		},
		{
			"Valid and existing zipcode. Should return a city",
			args{
				zipCode: "71218010",
			},
			"Brasília",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mock.FindCity(tt.args.zipCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindCity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindCity() = %v, want %v", got, tt.want)
			}
		})
	}
}
