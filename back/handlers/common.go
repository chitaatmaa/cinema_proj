<<<<<<< HEAD
package handlers

import (
	"cinema_proj/back/dbx"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := SearchUser(dbx.DB, request.Login)
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
		Login      string `json:"login"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		MiddleName string `json:"middle_name"`
		BirthDate  string `json:"birth_date"`
		RoleName   string `json:"role_name"`
	}{
		Login:      user.Login,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		BirthDate:  user.BirthDate,
		RoleName:   user.RoleName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserPhoto(login string) ([]byte, error) {
	var photo []byte
	err := dbx.DB.QueryRow(`
        SELECT p.photo 
        FROM cinema.persons p
        JOIN cinema.users u ON u.person_id = p.id
        WHERE u.login = $1`, login).Scan(&photo)
	return photo, err
}
=======
package handlers

import (
	"cinema_proj/back/dbx"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := SearchUser(dbx.DB, request.Login)
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
		Login      string `json:"login"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		MiddleName string `json:"middle_name"`
		BirthDate  string `json:"birth_date"`
		RoleName   string `json:"role_name"`
	}{
		Login:      user.Login,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		BirthDate:  user.BirthDate,
		RoleName:   user.RoleName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserPhoto(login string) ([]byte, error) {
	var photo []byte
	err := dbx.DB.QueryRow(`
        SELECT p.photo 
        FROM cinema.persons p
        JOIN cinema.users u ON u.person_id = p.id
        WHERE u.login = $1`, login).Scan(&photo)
	return photo, err
}
>>>>>>> 741fb8c1d90e4e1b14d660659e9dfa19713f6128
