package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func init() {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		// Генерация временного ключа для разработки
		log.Println("WARNING: SESSION_KEY not set, using temporary key")
		key = generateTempKey()
	}
	Store = sessions.NewCookieStore([]byte(key))

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   os.Getenv("GO_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	}
}

func generateTempKey() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("Failed to generate session key: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}
