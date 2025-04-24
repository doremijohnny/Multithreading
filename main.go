package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

func viaCEP(cep string, channel chan<- string) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%v/json/", cep)
	req, err := http.Get(url)
	if err != nil {
		response := fmt.Sprintf("erro na criacao de chamada viaCEP: %v", err)
		channel <- response
		return
	}
	defer req.Body.Close()

	response, err := io.ReadAll(req.Body)
	if err != nil {
		response := fmt.Sprintf("erro na criacao de chamada viaCEP: %v", err)
		channel <- response
		return

	}
	channel <- string(response)
	return
}

func brasilAPI(cep string, channel chan<- string) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%v", cep)
	req, err := http.Get(url)
	if err != nil {
		response := fmt.Sprintf("erro na criacao de chamada brasilAPI: %v", err)
		channel <- response
		return
	}
	defer req.Body.Close()

	response, err := io.ReadAll(req.Body)
	if err != nil {
		response := fmt.Sprintf("erro na leitura de response brasilAPI: %v", err)
		channel <- response
		return
	}
	channel <- string(response)
	return
}

func main() {
	cep := flag.String("cep", "09895070", "CEP a ser utilizado na consulta")
	c1 := make(chan string)
	c2 := make(chan string)
	go viaCEP(*cep, c1)
	go brasilAPI(*cep, c2)
	select {
	case msg1 := <-c1:
		fmt.Printf("API Utilizada: ViaCEP \nResponse: %v\n", msg1)

	case msg2 := <-c2:
		fmt.Printf("API Utilizada: BrasilAPI \nResponse: %v\n", msg2)

	case <-time.After(time.Second * 1):
		println("Timeout na consulta")
	}
}
