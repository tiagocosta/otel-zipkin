package zipcodeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

func FindCity(zipCode string) (string, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + zipCode + "/json/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return "", err
	}

	return c.Localidade, nil
}
