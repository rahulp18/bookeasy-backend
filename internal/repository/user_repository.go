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

func CreateUser(user *models.User) error {
	query := `INSERT INTO users(id,name,email,password) VALUES ($1,$2,$3,$4)`
	_, err := db.DB.Exec(query, user.ID, user.Name, user.Email, user.Password)
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id,name,email,password,created_at FROM users WHERE email=$1`
	err := db.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
