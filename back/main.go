package main

import (
	"cinema_proj/back/dbx"
	"cinema_proj/back/handlers"
	"log"
	"net/http"
)

func main() {
	dbx.ConnToDB()
	defer dbx.DB.Close()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("front/static"))))

	// Роуты для регистрации/авторизации (Create, Read)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/auth", http.StatusFound)
	})
	http.HandleFunc("/auth", handlers.AuthHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Роуты для работы администратора, продюссера, режиссёра
	//Роуты для работы с пользователями (регистрация, удаление)
	http.HandleFunc("/admin", handlers.AdminPanel)
	http.HandleFunc("/admin/user_photo", handlers.GetUserPhotoHandler)
	http.HandleFunc("/admin/user_data", handlers.GetUserDataHandler)
	http.HandleFunc("/admin/search", handlers.SearchUsersHandler)
	http.HandleFunc("/admin/delete", handlers.DeleteUserHandler)
	//Роусты для добавления фильмов и привязки к ним продюссера и режиссера
	http.HandleFunc("/admin/regis_data", handlers.GetRegisDataHandler)
	http.HandleFunc("/admin/prod_data", handlers.GetProdDataHandler)
	http.HandleFunc("/admin/create_film", handlers.AddMovieHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
