package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	appToken    = "d28721be-fd2d-4b45-869e-9f253b554e50"
	promoID     = "43e35910-c168-4634-ad4f-52fd764a843f"
	eventsDelay = 20000 // in milliseconds
	keyCount    = 4     // Example key count
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

func generateClientID() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	randomNumbers := randSeq(19)
	return fmt.Sprintf("%d-%s", timestamp, randomNumbers)
}

func randSeq(n int) string {
	const letters = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func login(clientID string) (string, error) {
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

func emulateProgress(clientToken string) (bool, error) {
	url := "https://api.gamepromo.io/promo/register-event"
	body := map[string]interface{}{
		"promoId":     promoID,
		"eventId":     randSeq(36), // Assuming UUID generation is replaced by random sequence
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

func generateKey(clientToken string) (string, error) {
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

func main() {
	var wg sync.WaitGroup

	for i := 0; i < keyCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			startTime := time.Now()

			clientID := generateClientID()
			clientToken, err := login(clientID)
			if err != nil {
				log.Fatalf("Login failed: %v", err)
			}

			hasCode := false
			count := 0

			for !hasCode {
				count++
				time.Sleep(2 * time.Second)

				hasCode, err := emulateProgress(clientToken)
				fmt.Println("hasCode: ", hasCode)
				if err != nil {
					log.Fatalf("Emulate progress failed: %v", err)
				}
				if hasCode {
					break
				}
			}
			fmt.Println(count)

			promoCode, err := generateKey(clientToken)
			if err != nil {
				log.Fatalf("Generate key failed: %v", err)
				return
			}
			fmt.Printf("Generated key: %s\n", promoCode)
			duration := time.Since(startTime) // Вычисляем время выполнения
			fmt.Printf("Time taken for key %d: %v\n", i+1, duration)

		}(i)

	}

	wg.Wait()
}
