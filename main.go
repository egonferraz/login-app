package main

import (
	"html/template"
	"net/http"

	"login-app/database"
	"login-app/handlers"
	"login-app/middleware"
)

func main() {

	err := database.InitDB()
	if err != nil {
		panic(err)
	}

	//models.CreateUser("Egon", "egon", "123")

	handlers.Tmpl = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", middleware.Auth(handlers.Register))
	http.HandleFunc("/dashboard", middleware.Auth(handlers.Dashboard))
	http.HandleFunc("/logout", middleware.Auth(handlers.Logout))
	http.HandleFunc("/users", middleware.Auth(handlers.Users))
	http.HandleFunc("/delete-user", middleware.Auth(handlers.DeleteUser))
	http.HandleFunc("/edit-user", middleware.Auth(handlers.EditUser))

	println("Servidor em http://localhost:8080/login")

	http.ListenAndServe(":8080", nil)
}
