package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/cep/{cep}", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	brasilAPI := make(chan string)
	viaCep := make(chan string)

	cep := r.PathValue("cep")

	go requestViaCep(cep, brasilAPI)
	go requestBrasilapi(cep, viaCep)

	select {
	case response := <-brasilAPI:
		responseHandler(w, "API: BrasilAPI: "+response)
	case response := <-viaCep:
		responseHandler(w, "API: ViaCep: "+response)
	case <-time.After(time.Second):
		println("timeout")
	}
}

func responseHandler(w http.ResponseWriter, response string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func requestViaCep(cep string, brasilAPI chan string) {
	response, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	brasilAPI <- string(body)
	defer response.Body.Close()
}
func requestBrasilapi(cep string, viaCep chan string) {
	response, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	viaCep <- string(body)
	defer response.Body.Close()
}
