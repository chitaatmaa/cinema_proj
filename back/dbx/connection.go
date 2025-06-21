package dbx

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnToDB() {
	connStr := "host=localhost port=5432 user=postgres password=7258 dbname=postgres sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
	log.Println("Connected to DB")
}

func SaveUser(roleID int, firstName, lastName, middleName, birthDate string, photo []byte, expYears int, login, pass string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Парсим дату рождения
	parsedBirthDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return fmt.Errorf("invalid birth date format: %w", err)
	}

	var personID int
	err = tx.QueryRow(`
		INSERT INTO cinema.persons (first_name, last_name, middle_name, birth_date, photo, experience_years) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		firstName, lastName, middleName, parsedBirthDate, photo, expYears,
	).Scan(&personID)
	if err != nil {
		return fmt.Errorf("insert person: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO cinema.users (login, pass, role_id, person_id) 
		VALUES ($1, $2, $3, $4)`,
		login, pass, roleID, personID,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return tx.Commit()
}

// Новая функция для получения фото по логину
func GetUserPhoto(login string) ([]byte, error) {
	var photo []byte
	err := DB.QueryRow(`
		SELECT p.photo
		FROM cinema.persons p
		JOIN cinema.users u ON u.person_id = p.id
		WHERE u.login = $1
	`, login).Scan(&photo)
	return photo, err
}

func GetRoles() (*sql.Rows, error) {
	return DB.Query("SELECT id, name FROM cinema.user_roles WHERE id != 0")
}

func GetRoles_admin() (*sql.Rows, error) {
	return DB.Query("SELECT id, name FROM cinema.user_roles")
}
