package api

import (
	"bytes"
	"encoding/json"
	"ham/internal/useCase"
	"io"
	"net/http"
)

type loginResponse struct {
	ClientToken string `json:"clientToken"`
}

type registerEventResponse struct {
	HasCode bool `json:"hasCode"`
}

type generateKeyResponse struct {
	PromoCode string `json:"promoCode"`
}

func Login(clientID, appToken string) (string, error) {
	url := "https://api.gamepromo.io/promo/login-client"
	body := map[string]interface{}{
		"appToken":     appToken,
		"clientId":     clientID,
		"clientOrigin": "deviceid",
	}

	response, err := postRequest(url, body)
	if err != nil {
		return "", err
	}

	var loginResp loginResponse
	if err := json.Unmarshal(response, &loginResp); err != nil {
		return "", err
	}
	return loginResp.ClientToken, nil
}

func EmulateProgress(clientToken, promoID string) (bool, error) {
	url := "https://api.gamepromo.io/promo/register-event"
	body := map[string]interface{}{
		"promoId":     promoID,
		"eventId":     useCase.RandSeq(36),
		"eventOrigin": "undefined",
	}

	response, err := postRequestWithAuth(url, body, clientToken)
	if err != nil {
		return false, err
	}

	var eventResp registerEventResponse
	if err := json.Unmarshal(response, &eventResp); err != nil {
		return false, err
	}
	return eventResp.HasCode, nil
}

func GenerateKey(clientToken, promoID string) (string, error) {
	url := "https://api.gamepromo.io/promo/create-code"
	body := map[string]interface{}{
		"promoId": promoID,
	}

	response, err := postRequestWithAuth(url, body, clientToken)
	if err != nil {
		return "", err
	}

	var keyResp generateKeyResponse
	if err := json.Unmarshal(response, &keyResp); err != nil {
		return "", err
	}
	return keyResp.PromoCode, nil
}

func postRequest(url string, body interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func postRequestWithAuth(url string, body interface{}, token string) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
