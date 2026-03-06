package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var tmpl = template.Must(template.ParseFiles(
	"templates/login.html",
	"templates/dashboard.html",
	"templates/register.html",
))

var db *sql.DB

type PageData struct {
	Username string
}

func loginPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if checkUser(username, password) {

			cookie := http.Cookie{
				Name:  "session",
				Value: username,
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		w.Write([]byte("AVISO: Usuário ou senha inválido"))
		return
	}

	tmpl.Execute(w, nil)

}

func dashboardPage(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")

	data := PageData{
		Username: cookie.Value,
	}
	tmpl.ExecuteTemplate(w, "dashboard.html", data)

}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func registerPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := createUser(username, password)

		if err != nil {
			w.Write([]byte("Erro ao criar usuário ou usuário existente"))
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "register.html", nil)
}

func createTable() {

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func createUser(username, password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO users(username, password) VALUES(?, ?)",
		username,
		string(hash),
	)

	return err
}

func checkUser(username, password string) bool {

	var hash string

	err := db.QueryRow(
		"SELECT password FROM users WHERE username = ?",
		username,
	).Scan(&hash)

	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		_, err := r.Cookie("session")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func main2() {

	var err error

	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}

	createTable()
	//createUser("Egon", "123")

	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/dashboard", authMiddleware(dashboardPage))
	http.HandleFunc("/logout", authMiddleware(logoutPage))
	http.HandleFunc("/register", authMiddleware(registerPage))

	println("Servidor rodando em http://localhost:8080/login")

	http.ListenAndServe(":8080", nil)

}
