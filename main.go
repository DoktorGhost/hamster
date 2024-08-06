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
	appToken1 = "d28721be-fd2d-4b45-869e-9f253b554e50" //Riding Extreme 3D
	promoID1  = "43e35910-c168-4634-ad4f-52fd764a843f" //Riding Extreme 3D

	appToken2 = "d1690a07-3780-4068-810f-9b5bbf2931b2" // Chain Cube
	promoID2  = "b4170868-cef0-424f-8eb9-be0622e8e8e3" // Chain Cube

	appToken3 = "74ee0b5b-775e-4bee-974f-63e7f4d5bacb" //My Clone Army
	promoID3  = "fe693b26-b342-4159-8808-15e3ff7f8767" //My Clone Army

	appToken4 = "82647f43-3f87-402d-88dd-09a90025313f" //Train Miner
	promoID4  = "c4480ac7-e178-4973-8061-9ed5b2e17954" //Train Miner

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

// генерация
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

// апи
func login(clientID, appToken string) (string, error) {
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

func emulateProgress(clientToken, promoID string) (bool, error) {
	url := "https://api.gamepromo.io/promo/register-event"
	body := map[string]interface{}{
		"promoId":     promoID,
		"eventId":     randSeq(36),
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

// апи
func generateKey(clientToken, promoID string) (string, error) {
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

// апи
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

// апи
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

// хендлер
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
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

	if keyCount > 10 {
		return
	}

	var wg sync.WaitGroup
	keys1 := make([]string, 0, keyCount)
	keys2 := make([]string, 0, keyCount)
	keys3 := make([]string, 0, keyCount)
	keys4 := make([]string, 0, keyCount)
	startTime := time.Now()
	//log.Printf("Запуск горутин. Количество: %d\n", keyCount*4)

	for i := 0; i < keyCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			startTimeGo := time.Now()

			clientID := generateClientID()
			clientToken, err := login(clientID, appToken1)
			if err != nil {
				log.Printf("Login failed: %v", err)
				return
			}

			hasCode := false

			for !hasCode {
				time.Sleep(2 * time.Second)

				hasCode, err := emulateProgress(clientToken, promoID1)
				if err != nil {
					log.Printf("Emulate progress failed: %v", err)
					return
				}
				if hasCode {
					break
				}
				duration := time.Since(startTimeGo)
				log.Printf("горутина appToken1-%d работает %s\n", i+1, duration)
			}

			promoCode, err := generateKey(clientToken, promoID1)
			if err != nil {
				log.Printf("Generate key failed: %v", err)
				return
			}

			keys1 = append(keys1, fmt.Sprintf("%s", promoCode))

		}(i)

	}

	for i := 0; i < keyCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			startTimeGo := time.Now()

			clientID := generateClientID()
			clientToken, err := login(clientID, appToken2)
			if err != nil {
				log.Printf("Login failed: %v", err)
				return
			}

			hasCode := false

			for !hasCode {
				time.Sleep(10 * time.Second)

				hasCode, err := emulateProgress(clientToken, promoID2)
				if err != nil {
					log.Printf("Emulate progress failed: %v", err)
					return
				}
				if hasCode {
					break
				}
				duration := time.Since(startTimeGo)
				log.Printf("горутина appToken2-%d работает %s\n", i+1, duration)
			}

			promoCode, err := generateKey(clientToken, promoID2)
			if err != nil {
				log.Printf("Generate key failed: %v", err)
				return
			}

			keys2 = append(keys2, fmt.Sprintf("%s", promoCode))

		}(i)

	}

	for i := 0; i < keyCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			startTimeGo := time.Now()

			clientID := generateClientID()
			clientToken, err := login(clientID, appToken3)
			if err != nil {
				log.Printf("Login failed: %v", err)
				return
			}

			hasCode := false

			for !hasCode {
				time.Sleep(2 * time.Second)

				hasCode, err := emulateProgress(clientToken, promoID3)
				if err != nil {
					log.Printf("Emulate progress failed: %v", err)
					return
				}
				if hasCode {
					break
				}
				duration := time.Since(startTimeGo)
				log.Printf("горутина appToken3-%d работает %s\n", i+1, duration)
			}

			promoCode, err := generateKey(clientToken, promoID3)
			if err != nil {
				log.Printf("Generate key failed: %v", err)
				return
			}

			keys3 = append(keys3, fmt.Sprintf("%s", promoCode))

		}(i)

	}

	for i := 0; i < keyCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			startTimeGo := time.Now()

			clientID := generateClientID()
			clientToken, err := login(clientID, appToken4)
			if err != nil {
				log.Printf("Login failed: %v", err)
				return
			}

			hasCode := false

			for !hasCode {
				time.Sleep(2 * time.Second)

				hasCode, err := emulateProgress(clientToken, promoID4)
				if err != nil {
					log.Printf("Emulate progress failed: %v", err)
					return
				}
				if hasCode {
					break
				}
				duration := time.Since(startTimeGo)
				log.Printf("горутина appToken4-%d работает %s\n", i+1, duration)
			}

			promoCode, err := generateKey(clientToken, promoID4)
			if err != nil {
				log.Printf("Generate key failed: %v", err)
				return
			}

			keys4 = append(keys4, fmt.Sprintf("%s", promoCode))

		}(i)

	}

	wg.Wait()
	duration := time.Since(startTime)
	log.Printf("Горутины отработали за %s \n", duration)

	// Читаем HTML-шаблон из файла
	tmplPath := "static/results.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Unable to parse template", http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон и передаем данные
	pageData := struct {
		Keys1 []string
		Keys2 []string
		Keys3 []string
		Keys4 []string
	}{
		Keys1: keys1,
		Keys2: keys2,
		Keys3: keys3,
		Keys4: keys4,
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Unable to execute template", http.StatusInternalServerError)
		return
	}
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate_keys", generateKeysHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Serve static files

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
