package handlers

import (
	"fmt"
	"ham/internal/api"
	"ham/internal/conf"
	"ham/internal/useCase"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func InitRout() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate_keys", generateKeysHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func generateKeysHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	gameNameStr := r.FormValue("gameName")
	gameName, err := strconv.Atoi(gameNameStr)
	if err != nil || gameName < 0 {
		http.Error(w, "Invalid gameName parameter", http.StatusBadRequest)
		return
	}

	keyCountStr := r.FormValue("keyCount")
	keyCount, err := strconv.Atoi(keyCountStr)
	if err != nil || keyCount <= 0 {
		http.Error(w, "Invalid keyCount parameter", http.StatusBadRequest)
		return
	}

	timeOutStr := r.FormValue("time_out")
	timeOut, err := strconv.Atoi(timeOutStr)
	if err != nil || timeOut <= 0 {
		http.Error(w, "Invalid timeout parameter", http.StatusBadRequest)
		return
	}

	if keyCount > 10 || timeOut < 5 || timeOut > 20 {
		return
	}

	var wg sync.WaitGroup

	var n int
	var appToken, promoID, gName string

	if gameName == 99 {
		n = conf.CountGames
	} else {
		n = 1
		appToken = conf.MapToken[gameName].AppToken
		promoID = conf.MapToken[gameName].PromoID
	}

	keySlices := make([][]string, n)

	for i := 0; i < n; i++ {
		keySlices[i] = make([]string, 0, keyCount)
	}

	startTime := time.Now()

	for i := 0; i < n; i++ {
		if n != 1 {
			appToken = conf.MapToken[i].AppToken
			promoID = conf.MapToken[i].PromoID
			gName = conf.MapToken[i].GameName
		}
		keyGen(keyCount, &wg, gName, appToken, promoID, timeOut, &keySlices[i])
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
		KeySets  [][]string
		MapToken map[int]conf.Token
	}{
		KeySets:  keySlices,
		MapToken: conf.MapToken,
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Unable to execute template", http.StatusInternalServerError)
		return
	}
}

func keyGen(keyCount int, wg *sync.WaitGroup, gameName string, appToken, promoID string, timeOut int, arr *[]string) {
	for i := 0; i < keyCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			startTimeGo := time.Now()

			clientID := useCase.GenerateClientID()
			clientToken, err := api.Login(clientID, appToken)
			if err != nil {
				log.Printf("Login failed: %v", err)
				return
			}

			for {
				hasCode, err := api.EmulateProgress(clientToken, promoID)
				if err != nil {
					log.Printf("Emulate progress failed: %v", err)
					return
				}
				if hasCode {
					break
				}

				duration := time.Since(startTimeGo)

				if duration > time.Duration(timeOut)*time.Minute {
					log.Printf("Горутина appToken-%s(%d) завершена из-за превышения времени выполнения", gameName, i+1)
					return
				}

				log.Printf("горутина appToken-%s(%d) работает %s\n", gameName, i+1, duration)

				randomDelay := rand.Intn(11) + 10

				time.Sleep(time.Duration(randomDelay) * time.Second)
			}

			promoCode, err := api.GenerateKey(clientToken, promoID)
			if err != nil {
				log.Printf("Generate key failed: %v", err)
				return
			}

			*arr = append(*arr, fmt.Sprintf("%s", promoCode))

		}(i)

	}
}
