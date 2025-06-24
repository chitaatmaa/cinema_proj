package handlers

import (
	"cinema_proj/back/dbx"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Usero struct {
	Login           string
	FirstName       string
	LastName        string
	MiddleName      string
	BirthDate       string
	RoleUser        string
	ExperienceYears int
}

func SearchUser1(db *sql.DB, login string) (*Usero, error) {
	query := `
    SELECT 
        u.login,
        p.first_name,
        p.last_name,
        p.middle_name,
        TO_CHAR(p.birth_date, 'DD.MM.YYYY') AS birth_date,
		r.name AS role_name,
        p.experience_years
    FROM cinema.users u
    JOIN cinema.persons p ON u.person_id = p.id
	JOIN cinema.user_roles r ON u.role_id = r.id
    WHERE u.login = $1 AND r.id = 1
`

	user1 := &Usero{}
	err := db.QueryRow(query, login).Scan(
		&user1.Login,
		&user1.FirstName,
		&user1.LastName,
		&user1.MiddleName,
		&user1.BirthDate,
		&user1.RoleUser,
		&user1.ExperienceYears,
	)

	if err != nil {
		return nil, err
	}

	return user1, err
}

func SearchUser2(db *sql.DB, login string) (*Usero, error) {
	query := `
    SELECT 
        u.login,
        p.first_name,
        p.last_name,
        p.middle_name,
        TO_CHAR(p.birth_date, 'DD.MM.YYYY') AS birth_date,
		r.name AS role_name,
        p.experience_years
    FROM cinema.users u
    JOIN cinema.persons p ON u.person_id = p.id
	JOIN cinema.user_roles r ON u.role_id = r.id
    WHERE u.login = $1 AND r.id = 2
`

	user2 := &Usero{}
	err := db.QueryRow(query, login).Scan(
		&user2.Login,
		&user2.FirstName,
		&user2.LastName,
		&user2.MiddleName,
		&user2.BirthDate,
		&user2.RoleUser,
		&user2.ExperienceYears,
	)

	if err != nil {
		return nil, err
	}

	return user2, err
}

func GetProdDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Login string `json:"login"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user1, err := SearchUser1(dbx.DB, request.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Search user error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Login           string `json:"login"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		MiddleName      string `json:"middle_name"`
		BirthDate       string `json:"birth_date"`
		RoleUser        string `json:"role_name"`
		ExperienceYears int    `json:"experience_years"`
	}{
		Login:           user1.Login,
		FirstName:       user1.FirstName,
		LastName:        user1.LastName,
		MiddleName:      user1.MiddleName,
		BirthDate:       user1.BirthDate,
		RoleUser:        user1.RoleUser,
		ExperienceYears: user1.ExperienceYears,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetRegisDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Login string `json:"login"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user2, err := SearchUser2(dbx.DB, request.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Search user error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Login           string `json:"login"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		MiddleName      string `json:"middle_name"`
		BirthDate       string `json:"birth_date"`
		RoleUser        string `json:"role_name"`
		ExperienceYears int    `json:"experience_years"`
	}{
		Login:           user2.Login,
		FirstName:       user2.FirstName,
		LastName:        user2.LastName,
		MiddleName:      user2.MiddleName,
		BirthDate:       user2.BirthDate,
		RoleUser:        user2.RoleUser,
		ExperienceYears: user2.ExperienceYears,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		GenreID  int    `json:"genre_id"`
		StatusID int    `json:"status_id"`
		Producer string `json:"producer"`
		Regisser string `json:"regisser"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var movieID int
	err := dbx.DB.QueryRow(`
        INSERT INTO cinema.movies (title, genre_id, status_id, producer, regisser)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, req.Title, req.GenreID, req.StatusID, req.Producer, req.Regisser).Scan(&movieID)

	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Failed to add movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"movieID": movieID,
	})
}
