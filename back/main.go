package main

import (
	"cinema_proj/back/dbx"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/term"
)

func getPassword() (string, error) {
	fmt.Print("Введите пароль от БД: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

func main() {
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		var err error
		pass, err = getPassword()
		if err != nil {
			log.Fatal("Ошибка чтения пароля:", err)
		}
	}

	// Инициализация БД
	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres dbname=cinema password=%s sslmode=disable", pass)
	if err := dbx.InitDB(connStr); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	log.Println("Успешное подключение к БД")

	mux := http.NewServeMux()
	mux.HandleFunc("/main", dbx.ServeIndexHTML)
	handler := corsMiddleware(mux)

	serverAddr := ":8080"
	log.Printf("Сервер запущен на %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
