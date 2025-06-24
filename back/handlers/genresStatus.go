package handlers

import (
	"cinema_proj/back/dbx"
)

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetGenres() ([]Genre, error) {
	rows, err := dbx.DB.Query("SELECT id, name FROM cinema.genres ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var g Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, rows.Err()
}

func GetStatuses() ([]Status, error) {
	rows, err := dbx.DB.Query("SELECT id, name FROM cinema.statuses ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []Status
	for rows.Next() {
		var g Status
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, g)
	}
	return statuses, rows.Err()
}
