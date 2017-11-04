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
