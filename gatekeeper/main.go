package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Response struct {
	Token string `json:"token"`
}

const jwtKey = "your jwt key here"

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

func handler(w http.ResponseWriter, r *http.Request) {

	username, password, ok := r.BasicAuth()

	// dummy authentication
	if !ok || password != "1234" || username != "gocrud@example.com" {
		http.Error(w, "user not found", http.StatusForbidden)
		return
	}
	userID := 1
	customerID := 1

	jwtToken := token(userID, customerID)
	resp := Response{
		Token: jwtToken,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "autentication error", http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("starting gatekeeper at localhost:4000")

	err := http.ListenAndServe(":4000", http.HandlerFunc(handler))
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
