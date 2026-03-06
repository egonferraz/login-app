package handlers

import (
	"html/template"
	"net/http"

	"login-app/models"
)

var Tmpl *template.Template

type PageData struct {
	Nickname string
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		username := r.FormValue("username")
		password := r.FormValue("password")

		if models.CheckUser(username, password) {

			cookie := http.Cookie{
				Name:  "session",
				Value: username,
				Path:  "/",
			}

			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		Tmpl.ExecuteTemplate(w, "login.html", map[string]string{
			"Error": "Usuário ou senha inválidos",
		})
		return
	}

	Tmpl.ExecuteTemplate(w, "login.html", nil)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session")

	data := PageData{
		Nickname: cookie.Value,
	}
	Tmpl.ExecuteTemplate(w, "dashboard.html", data)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		nickname := r.FormValue("nickname")
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := models.CreateUser(nickname, username, password)

		if err != nil {
			w.Write([]byte("Erro ao criar usuário ou usuário existente"))
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	Tmpl.ExecuteTemplate(w, "register.html", nil)
}

func Users(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()

	if err != nil {
		http.Error(w, "Erro ao buscar usuários", 500)
		return
	}

	data := struct {
		Users []models.User
	}{
		Users: users,
	}

	Tmpl.ExecuteTemplate(w, "user.html", data)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	err := models.DeleteUser(id)

	if err != nil {
		http.Error(w, "Erro ao deletar usuário", 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func EditUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		id := r.FormValue("id")
		nickname := r.FormValue("nickname")
		username := r.FormValue("username")

		models.UpdateUser(id, nickname, username)

		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id")

	user, err := models.GetUserByID(id)

	if err != nil {
		http.Error(w, "Usuário não encontrado", 404)
		return
	}

	Tmpl.ExecuteTemplate(w, "edit-user.html", user)
}
