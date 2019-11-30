package main

import (
	"encoding/json"
	"net/http"
)

func handlePwdReset(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r) // pass it the cookie and returns the id of the currently logged in user

	decoder := json.NewDecoder(r.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		sendError(w, "Error parsing Registration JSON", err, http.StatusBadRequest)
		return
	}

	resetPwd(w, u.Password, uid)
}

func handleAccountDelete(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r) // pass it the cookie and returns the id of the currently logged in user
	q := `Delete From users where user_id = ?`

	_, err := db.Exec(q, uid)
	if err != nil {
		sendError(w, "Error could not delete user", err, http.StatusInternalServerError)
		return
	}
}
