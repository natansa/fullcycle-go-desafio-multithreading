package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	http.HandleFunc("/", handler)

	r := chi.NewRouter()
	r.Get("/{cep}", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctxHttp, cancelHttp := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancelHttp()

	cepRequest, err := http.NewRequestWithContext(ctxHttp, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println(err)
	}

	cepResponse, err := http.DefaultClient.Do(cepRequest)
	if err != nil {
		log.Println(err)
	}
	defer cepResponse.Body.Close()

	cepResult, err := io.ReadAll(cepResponse.Body)
	if err != nil {
		log.Println(err)
	}

	w.Write([]byte(cepResult))
}
