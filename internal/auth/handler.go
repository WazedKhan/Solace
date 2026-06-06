package auth

import (
	"fmt"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	password := "wazed"
	hashed_password, err := HashPassword(password)
	if err != nil {
		log.Println("failed to hash the password!", err)
	}
	fmt.Fprintln(w, hashed_password)
}
