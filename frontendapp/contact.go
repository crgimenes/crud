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

func handlerEditContact(w http.ResponseWriter, r *http.Request) {
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
	contactID := r.FormValue("id")

	listURL := fmt.Sprintf("http://localhost:2015/api/gocrud/public/client_contacts?id=$eq.%s", contactID)

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		var clientID int
		clientID, err = strconv.Atoi(r.FormValue("client_id"))
		if err != nil {
			http.Error(w, "form error", http.StatusInternalServerError)
			return
		}

		contact := contactModel{
			ClientID: clientID,
			Name:     name,
			Phone:    phone,
			Email:    email,
		}
		var b []byte
		b, err = json.Marshal(contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp *http.Response
		resp, err = httpClientHelper(
			token.Value,
			http.MethodPatch,
			listURL,
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
		fmt.Println("Update Contact:", string(b))
		clientURL := fmt.Sprintf("/clients/edit?id=%v", clientID)
		http.Redirect(w, r, clientURL, http.StatusSeeOther)
		return
	}
	resp, err := httpClientHelper(
		token.Value,
		http.MethodGet,
		listURL,
		nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contact := []contactModel{}
	err = json.NewDecoder(resp.Body).Decode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("%#v\n", contact[0])

	err = templateHelper(
		w,
		map[string]interface{}{
			"title":   "GOCRUD - Edit Contact",
			"contact": contact[0],
		},
		"contactsedit.html",
		"./assets/contactsedit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
