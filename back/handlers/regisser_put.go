package handlers

import (
	"cinema_proj/back/dbx"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type FilmDetails struct {
	ID     int         `json:"id"`
	Title  string      `json:"title"`
	Groups []FilmGroup `json:"groups"`
	Actors []FilmActor `json:"actors"`
}

type FilmGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cost int    `json:"cost"`
}

type FilmActor struct {
	ID            int    `json:"id"`
	Login         string `json:"login"`
	CharacterName string `json:"scenic"`
	Cost          int    `json:"cost1"`
}

type UpdateFilmRequest struct {
	FilmID int         `json:"film_id"`
	Groups []FilmGroup `json:"groups"`
	Actors []FilmActor `json:"actors"`
}

func GetFilmDetails(filmID int) (FilmDetails, error) {
	var details FilmDetails

	// Получаем основную информацию о фильме
	err := dbx.DB.QueryRow(`
        SELECT id, title FROM cinema.movies 
        WHERE id = $1`, filmID).Scan(&details.ID, &details.Title)
	if err != nil {
		return details, err
	}

	// Получаем группы фильма
	groupsQuery := `
        SELECT g.id, g.name, mg.cost 
        FROM cinema.movie_groups mg
        JOIN cinema.groups g ON g.id = mg.group_id
        WHERE mg.movie_id = $1`
	groupsRows, err := dbx.DB.Query(groupsQuery, filmID)
	if err != nil {
		return details, err
	}
	defer groupsRows.Close()

	for groupsRows.Next() {
		var group FilmGroup
		if err := groupsRows.Scan(&group.ID, &group.Name, &group.Cost); err != nil {
			return details, err
		}
		details.Groups = append(details.Groups, group)
	}

	// Получаем актеров фильма
	actorsQuery := `
        SELECT a.id, a.login, ma.character_name, ma.cost 
        FROM cinema.movie_actors ma
        JOIN cinema.actors a ON a.id = ma.actor_id
        WHERE ma.movie_id = $1`
	actorsRows, err := dbx.DB.Query(actorsQuery, filmID)
	if err != nil {
		return details, err
	}
	defer actorsRows.Close()

	for actorsRows.Next() {
		var actor FilmActor
		if err := actorsRows.Scan(&actor.ID, &actor.Login, &actor.CharacterName, &actor.Cost); err != nil {
			return details, err
		}
		details.Actors = append(details.Actors, actor)
	}

	return details, nil
}

func FilmDetailsHandler(w http.ResponseWriter, r *http.Request) {
	filmIDStr := r.URL.Query().Get("film_id")
	if filmIDStr == "" {
		http.Error(w, "film_id is required", http.StatusBadRequest)
		return
	}

	filmID, err := strconv.Atoi(filmIDStr)
	if err != nil {
		http.Error(w, "Invalid film ID", http.StatusBadRequest)
		return
	}

	details, err := GetFilmDetails(filmID)
	if err != nil {
		log.Printf("Get film details error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func UpdateFilmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateFilmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tx, err := dbx.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Удаляем существующие связи
	_, err = tx.Exec(`DELETE FROM cinema.movie_groups WHERE movie_id = $1`, req.FilmID)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete groups error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM cinema.movie_actors WHERE movie_id = $1`, req.FilmID)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete actors error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Вставляем обновленные группы
	for _, group := range req.Groups {
		_, err := tx.Exec(`
            INSERT INTO cinema.movie_groups (movie_id, group_id, cost)
            VALUES ($1, $2, $3)`,
			req.FilmID, group.ID, group.Cost)

		if err != nil {
			tx.Rollback()
			log.Printf("Group insert error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// Вставляем обновленных актеров
	for _, actor := range req.Actors {
		_, err := tx.Exec(`
            INSERT INTO cinema.movie_actors (movie_id, actor_id, cost, character_name)
            VALUES ($1, $2, $3, $4)`,
			req.FilmID, actor.ID, actor.Cost, actor.CharacterName)

		if err != nil {
			tx.Rollback()
			log.Printf("Actor insert error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
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
		"message": "Film updated successfully",
	})
}
