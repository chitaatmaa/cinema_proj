package handlers

import (
	"cinema_proj/back/dbx"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func AdminPanel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Println("ParseMultipartForm error:", err)
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		roleID, _ := strconv.Atoi(r.FormValue("role_id"))
		expYears, _ := strconv.Atoi(r.FormValue("experience_years"))
		login := r.FormValue("login")
		pass := r.FormValue("pass")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		middleName := r.FormValue("middle_name")
		birthDate := r.FormValue("birth_date")

		var photo []byte
		defaultPhotoPath := "/static/images/default_user_photo.jpg"

		file, header, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()

			buff := make([]byte, 512)
			if _, err = file.Read(buff); err != nil {
				log.Println("Photo header read error:", err)
				photo, _ = os.ReadFile(defaultPhotoPath)
			} else {
				contentType := http.DetectContentType(buff)
				if !strings.HasPrefix(contentType, "image/") {
					log.Println("Invalid file type:", contentType)
					http.Error(w, "Only images are allowed", http.StatusBadRequest)
					return
				}

				if header.Size > 5<<20 {
					log.Println("File too large:", header.Size)
					http.Error(w, "File too large (max 5MB)", http.StatusBadRequest)
					return
				}
				if _, err = file.Seek(0, 0); err != nil {
					log.Println("File seek error:", err)
					photo, _ = os.ReadFile(defaultPhotoPath)
				}
			}
		} else {
			log.Println("No photo uploaded:", err)
			photo, _ = os.ReadFile(defaultPhotoPath)
		}

		err = dbx.SaveUser(
			roleID,
			firstName,
			lastName,
			middleName,
			birthDate,
			photo,
			expYears,
			login,
			pass,
		)

		if err != nil {
			log.Println("SaveUser error:", err)
			http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Registration successful!"))
	} else if r.Method == "GET" {
		rows, err := dbx.GetRoles_admin()
		if err != nil {
			log.Println("GetRoles error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Role struct {
			ID   int
			Name string
		}
		var roles []Role
		for rows.Next() {
			var r Role
			if err := rows.Scan(&r.ID, &r.Name); err != nil {
				log.Println("Role scan error:", err)
				continue
			}
			roles = append(roles, r)
		}

		tmpl := template.Must(template.ParseFiles("front/templates/admin.html"))
		if err := tmpl.Execute(w, roles); err != nil {
			log.Println("Template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

type User struct {
	Login      string `json:"login"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"`
	RoleName   string `json:"role_name"`
}

type SearchRequest struct {
	Login string `json:"login"`
}

type Response struct {
	User  *User  `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

type DeleteRequest struct {
	Login string `json:"login"`
}

func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Error: "Method not allowed"})
		return
	}

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "Invalid request"})
		return
	}

	if req.Login == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "Login is required"})
		return
	}
	user, err := SearchUser(dbx.DB, req.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Response{})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Database error"})
		log.Printf("DB query error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{User: user})
}

func SearchUser(db *sql.DB, login string) (*User, error) {
	query := `
        SELECT 
            u.login,
            p.first_name,
            p.last_name,
            p.middle_name,
            TO_CHAR(p.birth_date, 'DD.MM.YYYY') AS birth_date,
            r.name AS role_name
        FROM cinema.users u
        JOIN cinema.persons p ON u.person_id = p.id
        JOIN cinema.user_roles r ON u.role_id = r.id
        WHERE u.login = $1
    `

	user := &User{}
	err := db.QueryRow(query, login).Scan(
		&user.Login,
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
		&user.BirthDate,
		&user.RoleName,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Error: "Method not allowed"})
		return
	}

	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "Invalid request"})
		return
	}

	if req.Login == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "Login is required"})
		return
	}

	if err := deleteUser(dbx.DB, req.Login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Delete failed: " + err.Error()})
		log.Printf("Delete user error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{})
}

func deleteUser(db *sql.DB, login string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var userID, personID int
	err = tx.QueryRow(`
		SELECT u.id, u.person_id 
		FROM cinema.users u 
		WHERE u.login = $1
	`, login).Scan(&userID, &personID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return err
	}

	_, err = tx.Exec("DELETE FROM cinema.users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	_, err = tx.Exec("DELETE FROM cinema.persons WHERE id = $1", personID)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

var (
	photoCache = struct {
		sync.RWMutex
		m map[string][]byte
	}{m: make(map[string][]byte)}
)

func GetUserPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	login := r.URL.Query().Get("login")
	if login == "" {
		http.Error(w, "Login parameter is required", http.StatusBadRequest)
		return
	}

	// Проверка кэша
	photoCache.RLock()
	photo, cached := photoCache.m[login]
	photoCache.RUnlock()

	if !cached {
		var err error
		photo, err = dbx.GetUserPhoto(login)
		if err != nil {
			log.Printf("Photo fetch error for %s: %v", login, err)
		}

		// Обновление кэша
		photoCache.Lock()
		photoCache.m[login] = photo
		photoCache.Unlock()
	}

	// Определение Content-Type
	contentType := http.DetectContentType(photo)
	if !strings.HasPrefix(contentType, "image/") {
		contentType = "image/jpeg"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(photo)
}
