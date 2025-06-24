<<<<<<< HEAD
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

	// Общие роуты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/auth", http.StatusFound)
	})
	http.HandleFunc("/auth", handlers.AuthHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Админские роуты
	http.HandleFunc("/admin", handlers.AdminPanel)
	http.HandleFunc("/admin/user_photo", handlers.GetUserPhotoHandler)
	http.HandleFunc("/admin/user_data", handlers.GetUserDataHandler)
	http.HandleFunc("/admin/search", handlers.SearchUsersHandler)
	http.HandleFunc("/admin/delete", handlers.DeleteUserHandler)
	http.HandleFunc("/admin/regis_data", handlers.GetRegisDataHandler)
	http.HandleFunc("/admin/prod_data", handlers.GetProdDataHandler)
	http.HandleFunc("/admin/create_film", handlers.MovieHandler)

	// Режиссерские роуты
	http.HandleFunc("/regisser", handlers.RegisserMainHandler)
	http.HandleFunc("/regisser/add_group", handlers.AddGroupHandler)
	http.HandleFunc("/regisser/add_actor", handlers.AddActorHandler)
	http.HandleFunc("/regisser/start_film", handlers.StartFilmHandler)
	http.HandleFunc("/regisser/film_details", handlers.FilmDetailsHandler)
	http.HandleFunc("/regisser/update_film", handlers.UpdateFilmHandler)

	// Продюсерские роуты
	http.HandleFunc("/producer", handlers.ProducerHandler)
	http.HandleFunc("/producer/movies", handlers.MoviesHandler)
	http.HandleFunc("/producer/movie/budget", handlers.BudgetHandler)
	http.HandleFunc("/producer/report/groups", handlers.GroupsReportHandler)
	http.HandleFunc("/producer/report/detailed", handlers.DetailedReportHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
=======
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

	// Общие роуты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/auth", http.StatusFound)
	})
	http.HandleFunc("/auth", handlers.AuthHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Админские роуты
	http.HandleFunc("/admin", handlers.AdminPanel)
	http.HandleFunc("/admin/user_photo", handlers.GetUserPhotoHandler)
	http.HandleFunc("/admin/user_data", handlers.GetUserDataHandler)
	http.HandleFunc("/admin/search", handlers.SearchUsersHandler)
	http.HandleFunc("/admin/delete", handlers.DeleteUserHandler)
	http.HandleFunc("/admin/regis_data", handlers.GetRegisDataHandler)
	http.HandleFunc("/admin/prod_data", handlers.GetProdDataHandler)
	http.HandleFunc("/admin/create_film", handlers.MovieHandler)

	// Режиссерские роуты
	http.HandleFunc("/regisser", handlers.RegisserMainHandler)
	http.HandleFunc("/regisser/add_group", handlers.AddGroupHandler)
	http.HandleFunc("/regisser/add_actor", handlers.AddActorHandler)
	http.HandleFunc("/regisser/start_film", handlers.StartFilmHandler)
	http.HandleFunc("/regisser/film_details", handlers.FilmDetailsHandler)
	http.HandleFunc("/regisser/update_film", handlers.UpdateFilmHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
>>>>>>> 741fb8c1d90e4e1b14d660659e9dfa19713f6128
