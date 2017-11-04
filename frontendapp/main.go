package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("index.html").ParseFiles("./assets/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := "please login"
	userDataCookie, _ := r.Cookie("name")
	if userDataCookie != nil {
		msg = fmt.Sprintf("loged as %v", userDataCookie.Value)
	}

	tplMap := map[string]string{
		"title":   "GOCRUD - Welcome",
		"message": msg,
	}

	err = tpl.Execute(w, tplMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	tplMap := map[string]string{
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

		req, err := http.NewRequest(http.MethodPost, "http://localhost:4000", nil)
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
			Name:  "token",
			Value: gatekeeperRet.Token,
		}

		userDataCookie := http.Cookie{
			Name:  "name",
			Value: user,
		}

		http.SetCookie(w, &authCookie)
		http.SetCookie(w, &userDataCookie)

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	err := templateHelper(w, tplMap, "login.html", "./assets/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func templateHelper(w http.ResponseWriter, tplMap map[string]string, name string, file string) (err error) {
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

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
