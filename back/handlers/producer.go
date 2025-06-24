package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"cinema_proj/back/dbx"
)

func ProducerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/templates/producer.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Ошибка рендеринга шаблона: %v", err)
	}
}

func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	rows, err := dbx.DB.Query(`
        SELECT 
            m.id AS movie_id,
            m.title AS movie_title,
            s.name AS status_name,
            COALESCE(m.budget, 0) AS movie_budget
        FROM cinema.movies m
        JOIN cinema.statuses s ON m.status_id = s.id
        WHERE m.producer = $1
    `, login)

	if err != nil {
		log.Printf("Ошибка SQL запроса: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Movie struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
		Budget int    `json:"budget"`
	}
	var movies []Movie

	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Status, &m.Budget); err != nil {
			log.Printf("Ошибка сканирования строк: %v", err)
			continue
		}
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при обработке результатов: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

func BudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		MovieID int `json:"movie_id"`
		Budget  int `json:"budget"`
	}

	// Читаем и логируем тело запроса для отладки
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if data.MovieID == 0 {
		http.Error(w, "Movie ID is required", http.StatusBadRequest)
		return
	}

	if data.Budget < 0 {
		http.Error(w, "Budget cannot be negative", http.StatusBadRequest)
		return
	}

	_, err := dbx.DB.Exec("UPDATE cinema.movies SET budget = $1 WHERE id = $2", data.Budget, data.MovieID)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Бюджет успешно обновлен"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

func GroupsReportHandler(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	rows, err := dbx.DB.Query(`
        SELECT m.title, g.name, g.count
        FROM cinema.movies m
        JOIN cinema.movie_groups mg ON m.id = mg.movie_id
        JOIN cinema.groups g ON mg.group_id = g.id
        WHERE m.producer = $1
    `, login)

	if err != nil {
		log.Printf("Ошибка SQL запроса: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	report := make(map[string][]string)
	for rows.Next() {
		var movieTitle, groupName string
		var count int
		if err := rows.Scan(&movieTitle, &groupName, &count); err != nil {
			log.Printf("Ошибка сканирования строк: %v", err)
			continue
		}
		groupInfo := fmt.Sprintf("%s (Участников: %d)", groupName, count)
		report[movieTitle] = append(report[movieTitle], groupInfo)
	}

	content := "ОТЧЕТ ПО СЪЕМОЧНЫМ ГРУППАМ\n\n"
	content += fmt.Sprintf("Сгенерирован: %s\n\n", time.Now().Format("2006-01-02"))
	for movie, groups := range report {
		content += fmt.Sprintf("Фильм: %s\n", movie)
		content += "Съемочные группы:\n"
		for i, group := range groups {
			content += fmt.Sprintf("  %d. %s\n", i+1, group)
		}
		content += "\n"
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=groups_report.txt")
	_, err = w.Write([]byte(content))
	if err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
	}
}

func DetailedReportHandler(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	rows, err := dbx.DB.Query(`
        SELECT 
            m.title, 
            COALESCE(m.budget, 0) AS budget,
            g.name, 
            COALESCE(mg.cost, 0) AS group_cost,
            a.first_name, 
            a.middle_name, 
            a.last_name,
            COALESCE(ma.cost, 0) AS actor_cost,
            COALESCE(ma.character_name, 'N/A') AS character_name
        FROM cinema.movies m
        LEFT JOIN cinema.movie_groups mg ON m.id = mg.movie_id
        LEFT JOIN cinema.groups g ON mg.group_id = g.id
        LEFT JOIN cinema.movie_actors ma ON m.id = ma.movie_id
        LEFT JOIN cinema.actors a ON ma.actor_id = a.id
        WHERE m.producer = $1
    `, login)

	if err != nil {
		log.Printf("Ошибка SQL запроса: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Group struct {
		Name string
		Cost int
	}
	type Actor struct {
		Name      string
		Character string
		Cost      int
	}

	movies := make(map[string]struct {
		Budget int
		Groups []Group
		Actors []Actor
	})

	for rows.Next() {
		var (
			movieTitle      string
			movieBudget     int
			groupName       sql.NullString
			groupCost       sql.NullInt64
			actorFirstName  sql.NullString
			actorMiddleName sql.NullString
			actorLastName   sql.NullString
			actorCost       sql.NullInt64
			actorCharacter  sql.NullString
		)

		if err := rows.Scan(
			&movieTitle, &movieBudget,
			&groupName, &groupCost,
			&actorFirstName, &actorMiddleName, &actorLastName,
			&actorCost, &actorCharacter,
		); err != nil {
			log.Printf("Ошибка сканирования строк: %v", err)
			continue
		}

		movieData := movies[movieTitle]
		movieData.Budget = movieBudget

		// Добавление группы
		if groupName.Valid {
			movieData.Groups = append(movieData.Groups, Group{
				Name: groupName.String,
				Cost: int(groupCost.Int64), // Преобразование
			})
		}

		// Добавление актера
		if actorLastName.Valid {
			nameParts := []string{}
			if actorLastName.String != "" {
				nameParts = append(nameParts, actorLastName.String)
			}
			if actorFirstName.String != "" {
				nameParts = append(nameParts, actorFirstName.String)
			}
			if actorMiddleName.String != "" {
				nameParts = append(nameParts, actorMiddleName.String)
			}
			fullName := strings.Join(nameParts, " ")

			movieData.Actors = append(movieData.Actors, Actor{
				Name:      fullName,
				Character: actorCharacter.String,
				Cost:      int(actorCost.Int64), // Преобразование
			})
		}

		movies[movieTitle] = movieData
	}

	content := "ДЕТАЛЬНЫЙ ОТЧЕТ ПО ПРОЕКТАМ\n\n"
	content += fmt.Sprintf("Сгенерирован: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	for movie, data := range movies {
		content += fmt.Sprintf("ФИЛЬМ: %s\n", movie)
		content += fmt.Sprintf("Бюджет проекта: $%d\n\n", data.Budget) // Изменено %f на %d

		content += "Съемочные группы:\n"
		if len(data.Groups) > 0 {
			for _, group := range data.Groups {
				content += fmt.Sprintf("  - %s: $%d\n", group.Name, group.Cost) // Изменено %f на %d
			}
		} else {
			content += "  Группы не назначены\n"
		}

		content += "\nАктеры:\n"
		if len(data.Actors) > 0 {
			for _, actor := range data.Actors {
				content += fmt.Sprintf("  - %s (%s): $%d\n", // Изменено %f на %d
					actor.Name, actor.Character, actor.Cost)
			}
		} else {
			content += "  Актеры не назначены\n"
		}

		content += "\n" + strings.Repeat("=", 50) + "\n\n"
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=detailed_report.txt")
	_, err = w.Write([]byte(content))
	if err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
	}
}
