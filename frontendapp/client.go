package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*

list
add
edit
delete

*/

type clientModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func handlerListClient(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:2015/api/gocrud/public/client", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %v", token.Value))

	client := http.Client{}
	resp, err := client.Do(req)
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
