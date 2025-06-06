package dbx

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filename))
}

func ServeIndexHTML(w http.ResponseWriter, r *http.Request) {
	root := getProjectRoot()
	filePath := filepath.Join(root, "..", "front", "index.html")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Файл не найден: %s", filePath)
		http.Error(w, "Index file not found", http.StatusNotFound)
		return
	}

	log.Printf("Открываем файл: %s", filePath)
	http.ServeFile(w, r, filePath)
}
