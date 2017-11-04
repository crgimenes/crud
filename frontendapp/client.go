package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

/*

list
add
edit
delete

*/

type clientModel struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omiempty"`
}

type contactModel struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omiempty"`
	Phone string `json:"fone,omiempty"`
	Email string `json:"email,omiempty"`
}

func handlerListClient(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	resp, err := httpClientHelper(token.Value, http.MethodGet, "http://localhost:2015/api/gocrud/public/client", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret := []clientModel{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templateHelper(
		w,
		map[string]interface{}{"clients": ret},
		"clientslist.html",
		"./assets/clientslist.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handlerAddClient(w http.ResponseWriter, r *http.Request) {
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
		name := r.FormValue("name")
		client := clientModel{Name: name}
		b, err := json.Marshal(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := httpClientHelper(
			token.Value,
			http.MethodPost,
			"http://localhost:2015/api/gocrud/public/client",
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
		fmt.Println("Insert Client:", string(b))
		http.Redirect(w, r, "/clients", http.StatusTemporaryRedirect)
		return
	}
	err := templateHelper(
		w,
		map[string]interface{}{
			"title": "GOCRUD - Add Client",
		},
		"clientsadd.html",
		"./assets/clientsadd.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlerEditClient(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	clientID := r.FormValue("id")

	listURL := fmt.Sprintf("http://localhost:2015/api/gocrud/public/client?id=$eq.%s", clientID)
	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			http.Error(w, "form error", http.StatusInternalServerError)
			return
		}
		name := r.FormValue("name")
		client := clientModel{Name: name}
		var b []byte
		b, err = json.Marshal(client)
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
		fmt.Println("Update Client:", string(b))
		http.Redirect(w, r, "/clients", http.StatusTemporaryRedirect)
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

	client := []clientModel{}
	err = json.NewDecoder(resp.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contactsList := fmt.Sprintf("http://localhost:2015/api/gocrud/public/client_contacts?client_id=%s", clientID)
	resp, err = httpClientHelper(
		token.Value,
		http.MethodGet,
		contactsList,
		nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contacts := []contactModel{}
	err = json.NewDecoder(resp.Body).Decode(&contacts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templateHelper(
		w,
		map[string]interface{}{
			"title":    "GOCRUD - Edit Client",
			"client":   client[0],
			"contacts": contacts,
		},
		"clientsedit.html",
		"./assets/clientsedit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlerDeleteClient(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	clientID := r.URL.Query().Get("id")
	deleteURL := fmt.Sprintf("http://localhost:2015/api/gocrud/public/client_contacts?client_id=%s", clientID)
	_, err = httpClientHelper(
		token.Value,
		http.MethodDelete,
		deleteURL,
		nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	deleteURL = fmt.Sprintf("http://localhost:2015/api/gocrud/public/client?id=%s", clientID)
	_, err = httpClientHelper(
		token.Value,
		http.MethodDelete,
		deleteURL,
		nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/clients", http.StatusTemporaryRedirect)
}

func httpClientHelper(token, method, url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %v", token))

	client := http.Client{}
	resp, err = client.Do(req)
	return
}
