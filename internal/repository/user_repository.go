package repository

import (
	"github.com/rahulp18/bookeasy-backend/internal/db"
	"github.com/rahulp18/bookeasy-backend/internal/models"
)

func GetAllUsers() ([]models.User, error) {
	rows, err := db.DB.Query("SELECT id, email, password, name, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
