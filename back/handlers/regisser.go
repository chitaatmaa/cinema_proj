package handlers

import (
	"bytes"
	"cinema_proj/back/dbx"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Film struct {
	ID    int    `json:"ID"`
	Title string `json:"Title"`
}
type FGroup struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type CGroup struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type CActor struct {
	Login      string `json:"login"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"`
	Experience int    `json:"experience"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type FActor struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type FTemplateData struct {
	Films  []Film
	Groups []FGroup
	Actors []FActor
}

func GetAllFilms(login string) ([]Film, error) {
	query := `SELECT f.id, f.title FROM cinema.movies f WHERE f.regisser = $1`
	rows, err := dbx.DB.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var f Film
		if err := rows.Scan(&f.ID, &f.Title); err != nil {
			return nil, err
		}
		films = append(films, f)
	}

	return films, nil
}

func GetAllGroups() ([]FGroup, error) {
	query := `SELECT f.id, f.name, f.count FROM cinema.groups f`
	rows, err := dbx.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allGroups []FGroup
	for rows.Next() {
		var fg FGroup
		if err := rows.Scan(&fg.ID, &fg.Name, &fg.Count); err != nil {
			return nil, err
		}
		allGroups = append(allGroups, fg)
	}

	return allGroups, nil
}

func GetAllActors() ([]FActor, error) {
	query := `SELECT f.id, f.login FROM cinema.actors f`
	rows, err := dbx.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allActors []FActor
	for rows.Next() {
		var fa FActor
		if err := rows.Scan(&fa.ID, &fa.Login); err != nil {
			return nil, err
		}
		allActors = append(allActors, fa)
	}

	return allActors, nil
}

func RegisserMainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		login := r.URL.Query().Get("login")
		if login == "" {
			http.Error(w, "Login parameter is required", http.StatusBadRequest)
			return
		}

		films, err := GetAllFilms(login)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		groups, err := GetAllGroups()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		actors, err := GetAllActors()
		if err != nil {
			log.Printf("Database insert error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}

		data := FTemplateData{
			Films:  films,
			Groups: groups,
			Actors: actors,
		}
		// Загружаем и выполняем шаблон
		tmpl2, err := template.ParseFiles("front/templates/regisser.html")
		if err != nil {
			log.Printf("Ошибка загрузки шаблона: %v", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		// БУФЕРИЗАЦИЯ ДАННЫХ
		buf := new(bytes.Buffer)
		err = tmpl2.Execute(buf, data)
		if err != nil {
			log.Printf("Ошибка выполнения шаблона: %v", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, err = buf.WriteTo(w)
		if err != nil {
			log.Printf("Ошибка отправки ответа: %v", err)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AddGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Декодируем JSON
	var req CGroup
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if req.Name == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}
	if req.Count <= 0 {
		http.Error(w, "Group count must be positive", http.StatusBadRequest)
		return
	}

	// Вставляем данные в БД
	var groupID int
	err := dbx.DB.QueryRow(`
		INSERT INTO cinema.groups (name, count)
		VALUES ($1, $2)
		RETURNING id
	`, req.Name, req.Count).Scan(&groupID)

	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"status":  "success",
		"message": "Group added successfully",
		"group": map[string]interface{}{
			"id":    groupID,
			"name":  req.Name,
			"count": req.Count,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func AddActorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Декодируем JSON
	var req CActor
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var actorID int
	err := dbx.DB.QueryRow(`
		INSERT INTO cinema.actors (login, first_name, last_name, middle_name, birth_date, experience, email, phone)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, req.Login, req.FirstName, req.LastName, req.MiddleName, req.BirthDate, req.Experience, req.Email, req.Phone).Scan(&actorID)

	if err != nil {
		log.Printf("Database error: %v", err)
		// Отправляем JSON-ошибку вместо текста
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"status": "success",
		"actor":  map[string]interface{}{"id": actorID},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

type FilmGroupRequest struct {
	GroupID int `json:"group_id"`
	Cost    int `json:"cost"`
}

type FilmActorRequest struct {
	ActorID       int    `json:"actor_id"`
	Cost          int    `json:"cost1"`
	CharacterName string `json:"scenic"`
}

type StartFilmRequest struct {
	FilmID int                `json:"film_id"`
	Groups []FilmGroupRequest `json:"groups"`
	Actors []FilmActorRequest `json:"actors"`
}

func StartFilmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartFilmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Начинаем транзакцию
	tx, err := dbx.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Вставляем группы
	for _, group := range req.Groups {
		_, err := tx.Exec(`
            INSERT INTO cinema.movie_groups (movie_id, group_id, cost)
            VALUES ($1, $2, $3)
            ON CONFLICT (movie_id, group_id) 
            DO UPDATE SET cost = EXCLUDED.cost`,
			req.FilmID, group.GroupID, group.Cost)

		if err != nil {
			tx.Rollback()
			log.Printf("Group insert error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}
	}

	// Вставляем актёров
	for _, actor := range req.Actors {
		_, err := tx.Exec(`
            INSERT INTO cinema.movie_actors (movie_id, actor_id, cost, character_name)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (movie_id, actor_id)
            DO UPDATE SET cost = EXCLUDED.cost, character_name = EXCLUDED.character_name`,
			req.FilmID, actor.ActorID, actor.Cost, actor.CharacterName)

		if err != nil {
			tx.Rollback()
			log.Printf("Actor insert error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Transaction commit error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Фильм успешно запущен!",
	})
}
