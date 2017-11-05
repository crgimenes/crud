package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("index.html").ParseFiles("./assets/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logged := false
	msg := "please login"
	userDataCookie, _ := r.Cookie("name")
	if userDataCookie != nil {
		logged = true
		msg = fmt.Sprintf("logged as %v", userDataCookie.Value)
	}

	tplMap := map[string]interface{}{
		"title":   "GOCRUD - Welcome",
		"message": msg,
		"logged":  logged,
	}

	err = tpl.Execute(w, tplMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	tplMap := map[string]interface{}{
		"title":   "GOCRUD - Login",
		"message": "",
	}

	if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "form error", http.StatusInternalServerError)
			return
		}

		user := r.FormValue("user")
		password := r.FormValue("password")

		req, err := http.NewRequest(http.MethodPost, "http://localhost:2015/gatekeeper", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.SetBasicAuth(user, password)
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode == http.StatusForbidden {
			tplMap["message"] = "User not found or incrorrect password"
			err = templateHelper(w, tplMap, "login.html", "./assets/login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if resp.StatusCode != http.StatusOK {
			tplMap["message"] = "Some error message that can be seen by the user"
			err = templateHelper(w, tplMap, "login.html", "./assets/login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		gatekeeperRet := struct {
			Token string `json:"token"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&gatekeeperRet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("token", gatekeeperRet.Token)

		authCookie := http.Cookie{
			Name:    "token",
			Value:   gatekeeperRet.Token,
			Expires: time.Now().Add(time.Duration(1) * time.Hour),
		}

		userDataCookie := http.Cookie{
			Name:    "name",
			Value:   user,
			Expires: time.Now().Add(time.Duration(1) * time.Hour),
		}

		http.SetCookie(w, &authCookie)
		http.SetCookie(w, &userDataCookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	err := templateHelper(w, tplMap, "login.html", "./assets/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	authCookie := http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	}

	userDataCookie := http.Cookie{
		Name:    "name",
		Value:   "",
		Expires: time.Now(),
	}

	http.SetCookie(w, &authCookie)
	http.SetCookie(w, &userDataCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func templateHelper(w http.ResponseWriter, tplMap map[string]interface{}, name string, file string) (err error) {
	tpl, err := template.New(name).ParseFiles(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, tplMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func main() {
	fmt.Println("starting frontend app at localhost:5000")

	http.HandleFunc("/", handler)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/clients", handlerListClient)
	http.HandleFunc("/clients/add", handlerAddClient)
	http.HandleFunc("/clients/edit", handlerEditClient)
	http.HandleFunc("/clients/delete", handlerDeleteClient)
	http.HandleFunc("/contacts/add", handlerAddContact)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
