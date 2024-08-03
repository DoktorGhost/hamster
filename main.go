package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	appToken    = "d28721be-fd2d-4b45-869e-9f253b554e50"
	promoID     = "43e35910-c168-4634-ad4f-52fd764a843f"
	eventsDelay = 20000 // in milliseconds
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func generateKeysHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	keyCountStr := r.FormValue("keyCount")
	keyCount, err := strconv.Atoi(keyCountStr)
	if err != nil || keyCount <= 0 {
		http.Error(w, "Invalid keyCount parameter", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	keys := make([]string, 0, keyCount)
	durations := make([]string, 0, keyCount)

	for i := 0; i < keyCount; i++ {
		wg.Add(1)

		log.Printf("Старт приложения\n")
		log.Printf("Запуск горутин. Количество: %d\n", keyCount)

		go func(i int) {
			defer wg.Done()

			startTime := time.Now()

			clientID := generateClientID()
			clientToken, err := login(clientID)
			if err != nil {
				log.Printf("Login failed: %v", err)
				return
			}

			hasCode := false
			count := 0

			for !hasCode {
				count++
				time.Sleep(2 * time.Second)

				hasCode, err := emulateProgress(clientToken)
				if err != nil {
					log.Printf("Emulate progress failed: %v", err)
					return
				}
				if hasCode {
					break
				}
				duration := time.Since(startTime)
				log.Printf("горутина %d работает %s\n", i+1, duration)
			}

			promoCode, err := generateKey(clientToken)
			if err != nil {
				log.Printf("Generate key failed: %v", err)
				return
			}

			duration := time.Since(startTime)
			keys = append(keys, fmt.Sprintf("%s", promoCode))
			durations = append(durations, fmt.Sprintf("Key %d: %v", i+1, duration))

		}(i)

	}

	wg.Wait()

	pageData := struct {
		Keys      []string
		Durations []string
	}{
		Keys:      keys,
		Durations: durations,
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>Генерация кодов</title>
	    <link rel="stylesheet" href="/static/css/styles.css">
	</head>
	<body>

	{{if .Keys}}
	<div id="keyContainer">
	    <h2>Сгенерированные ключи:</h2>
	    <ul>
	        {{range .Keys}}
	        <li>{{.}}</li>
	        {{end}}
	    </ul>
	</div>
	{{end}}

	{{if .Durations}}
	    <h3>Время генерации:</h3>
	    <ul>
	        {{range $index, $duration := .Durations}}
	        <li>{{$duration}}</li>
	        {{end}}
	    </ul>
	{{end}}
	</body>
	</html>`

	t := template.Must(template.New("result").Parse(tmpl))
	t.Execute(w, pageData)
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate_keys", generateKeysHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Serve static files

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
