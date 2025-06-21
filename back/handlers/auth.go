package handlers

import (
	"cinema_proj/back/dbx"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	RoleID int    `json:"role_id"`
	Error  string `json:"error,omitempty"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "front/templates/auth.html")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Method not allowed"})
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Invalid request"})
		return
	}

	if req.Login == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Login and password required"})
		return
	}

	var (
		storedPass string
		roleID     int
	)

	db := dbx.DB
	err := db.QueryRow(
		"SELECT pass, role_id FROM cinema.users WHERE login = $1",
		req.Login,
	).Scan(&storedPass, &roleID)

	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{Error: "User not found"})
		return
	case err != nil:
		log.Printf("Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Database error"})
		return
	case storedPass != req.Password:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Invalid password"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{RoleID: roleID})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rows, err := dbx.GetRoles()
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

		tmpl := template.Must(template.ParseFiles("front/templates/register.html"))
		if err := tmpl.Execute(w, roles); err != nil {
			log.Println("Template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	} else if r.Method == "POST" {
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
		defaultPhotoPath := filepath.Join("front", "static", "images", "default_user_photo.jpg")

		file, _, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()
			photo, err = io.ReadAll(file)
			if err != nil {
				log.Println("Error reading photo:", err)
				photo, _ = os.ReadFile(defaultPhotoPath)
			}
		} else {
			log.Println("No photo uploaded, using default")
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
	}
}
