package main

import (
	"html/template"
	"net/http"
	"os"

	"login-app/database"
	"login-app/handlers"
	"login-app/middleware"
)

func main() {

	err := database.InitDB()
	if err != nil {
		panic(err)
	}

	handlers.Tmpl = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", middleware.Auth(handlers.Register))
	http.HandleFunc("/dashboard", middleware.Auth(handlers.Dashboard))
	http.HandleFunc("/logout", middleware.Auth(handlers.Logout))
	http.HandleFunc("/users", middleware.Auth(handlers.Users))
	http.HandleFunc("/delete-user", middleware.Auth(handlers.DeleteUser))
	http.HandleFunc("/edit-user", middleware.Auth(handlers.EditUser))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	println("Servidor rodando na porta:", port)

	http.ListenAndServe(":"+port, nil)
}
