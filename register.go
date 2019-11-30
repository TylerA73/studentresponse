package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
* Handles Post call to /api/v1/register
* creates new teacher users cannot be used to create admins
**/
func handleRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newUser User

	err := decoder.Decode(&newUser)
	if err != nil {
		sendError(w, "Error parsing Registration JSON", err, http.StatusBadRequest)
		return
	}
	fmt.Printf("%v\n", newUser)
	encoded, err := encodePwd(newUser.Password)
	if err != nil {
		sendError(w, "Error encoding password", err, http.StatusInternalServerError)
		return
	}

	schema := `INSERT INTO users (username, first_name, last_name, password)
				VALUES (?, ?,?,?);`

	// execute a query on the server
	_, err = db.Exec(schema, newUser.Username, newUser.FirstName, newUser.LastName, encoded)
	if err != nil {
		if isDuplicateRowError(err) {
			sendError(w, "DB: Error User already exists"+err.Error(), err, http.StatusConflict)
		} else {
			sendError(w, "DB: Error Could not create user"+err.Error(), err, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
