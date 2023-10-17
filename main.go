package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	http.HandleFunc("/", handler)
	r := chi.NewRouter()
	r.Get("/{cep}", handler)
	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apicep := make(chan string)
	viacep := make(chan string)

	go func() {
		getApicep(apicep, cep)
	}()

	go func() {
		getViacep(viacep, cep)
	}()

	select {
	case cepResult := <-viacep:
		fmt.Printf("Resultado da VIACEP: %s\n", cepResult)

	case cepResult := <-apicep:
		fmt.Printf("Resultado da APICEP: %s\n", cepResult)

	case <-time.After(1 * time.Second):
		fmt.Println("TIMEOUT")
	}

	w.WriteHeader(http.StatusOK)
}

func getApicep(apicep chan string, cep string) {
	apicepResponse, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
	if err != nil {
		return
	}
	defer apicepResponse.Body.Close()

	apicepResult, err := io.ReadAll(apicepResponse.Body)
	if err != nil {
		return
	}

	if apicepResponse.StatusCode >= 400 {
		return
	}

	apicep <- string(apicepResult)
}

func getViacep(viacep chan string, cep string) {
	viacepResponse, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return
	}
	defer viacepResponse.Body.Close()

	viacepResult, err := io.ReadAll(viacepResponse.Body)
	if err != nil {
		return
	}

	if viacepResponse.StatusCode >= 400 {
		return
	}

	viacep <- string(viacepResult)
}
