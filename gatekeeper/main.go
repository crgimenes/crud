package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/prest/adapters"
	"github.com/prest/adapters/postgres"
	prestConf "github.com/prest/config"
)

type (
	response struct {
		Token string `json:"token"`
	}

	user struct {
		ID         int `json:"id"`
		CustomerID int `json:"customer_id"`
	}
)

const jwtKey = "your jwt key here"

var dbAdapter adapters.Adapter

func token(userID int, customerID int) string {
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	claims := jwt.StandardClaims{
		ExpiresAt: expireToken,
		Id:        strconv.Itoa(userID),
		IssuedAt:  time.Now().Unix(),
		Issuer:    strconv.Itoa(customerID),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := tok.SignedString([]byte(jwtKey))
	return signedToken
}

func getUserData(username, password string) (ret user, err error) {
	// default username=gocrud@example.com password=1234
	// password hash = BB06E371F6B25E45A7F2082115346F5EB78E4501BDFECD59F8DBC4643A7355BA
	passwordSum := sha256.Sum256([]byte(password + username))
	passwordHash := fmt.Sprintf("%X", passwordSum)
	sc := dbAdapter.Query("SELECT id,customer_id FROM users WHERE username=$1 AND password=$2", username, passwordHash)
	err = sc.Err()
	if err != nil {
		return
	}
	var n int
	ret = user{}
	n, err = sc.Scan(&ret)
	if err != nil {
		return
	}
	if n != 1 {
		err = fmt.Errorf("invalid username or password")
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "HTTP Basic Authentication required", http.StatusForbidden)
		return
	}

	u, err := getUserData(username, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "user not found", http.StatusForbidden)
		return
	}

	jwtToken := token(u.ID, u.CustomerID)
	resp := response{
		Token: jwtToken,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "autentication error", http.StatusInternalServerError)
		return
	}
}

func main() {

	prestConf.Load()
	postgres.Load()
	dbAdapter = prestConf.PrestConf.Adapter

	fmt.Println("starting gatekeeper at localhost:4000")

	err := http.ListenAndServe(":4000", http.HandlerFunc(handler))
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
