package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func handlerAddContact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		err = r.ParseForm()
		if err != nil {
			http.Error(w, "form error", http.StatusInternalServerError)
			return
		}

		clientID, err := strconv.Atoi(r.FormValue("client_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		contact := contactModel{
			ClientID: clientID,
			Name:     name,
			Phone:    phone,
			Email:    email,
		}
		b, err := json.Marshal(contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := httpClientHelper(
			token.Value,
			http.MethodPost,
			"http://localhost:2015/api/gocrud/public/client_contacts",
			bytes.NewReader(b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Insert Contact:", string(b))
		clientURL := fmt.Sprintf("/clients/edit?id=%d", clientID)
		http.Redirect(w, r, clientURL, http.StatusSeeOther)
		return
	}

	clientID := r.URL.Query().Get("client_id")

	err := templateHelper(
		w,
		map[string]interface{}{
			"title":    "GOCRUD - Add Contact",
			"clientID": clientID,
		},
		"contactsadd.html",
		"./assets/contactsadd.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
