package handlers

import (
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Сбрасываем данные аутентификации
	session.Values["authenticated"] = false
	delete(session.Values, "user_id")
	delete(session.Values, "username")

	session.Options.MaxAge = -1
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteLaxMode

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Дополнительно очищаем куку на клиенте
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}
