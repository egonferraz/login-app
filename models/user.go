package models

import (
	"login-app/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Nickname string
	Username string
}

func CreateUser(nickname, username, password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(
		"INSERT INTO users(nickname,username,password) VALUES(?,?,?)",
		nickname,
		username,
		string(hash),
	)

	return err
}

func CheckUser(username, password string) bool {

	var hash string

	err := database.DB.QueryRow(
		"SELECT password FROM users WHERE username = ?",
		username,
	).Scan(&hash)

	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GetUsers() ([]User, error) {
	rows, err := database.DB.Query("SELECT id, nickname, username FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Nickname, &u.Username)

		users = append(users, u)
	}

	return users, nil
}

func DeleteUser(id string) error {
	_, err := database.DB.Exec(
		"DELETE FROM users WHERE id=?",
		id,
	)
	return err
}

func UpdateUser(id, nickname, username string) error {

	_, err := database.DB.Exec(
		"UPDATE users SET nickname = ?, username = ? WHERE id = ?",
		nickname,
		username,
		id,
	)

	return err
}

func GetUserByID(id string) (User, error) {

	var user User

	err := database.DB.QueryRow(
		"SELECT id, nickname, username FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Nickname, &user.Username)

	return user, err
}
