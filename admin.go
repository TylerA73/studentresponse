package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

/**
*   Used to parse the URI for the query
**/
func getQuery(w http.ResponseWriter, r *http.Request) (uid, uname, fname, lname string, err error) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		sendError(w, "Error parsing URI", err, http.StatusBadRequest)
		return "", "", "", "", err
	}
	m, _ := url.ParseQuery(u.RawQuery)

	if len(m["id"]) > 0 {
		_, err = strconv.Atoi(m["id"][0])
		if err != nil {
			sendError(w, "Error convert string to int", err, http.StatusInternalServerError)
			return "", "", "", "", err
		}
		uid = m["id"][0]
	} else {
		uid = ".*"
	}
	if len(m["un"]) > 0 {
		uname = m["un"][0]
	} else {
		uname = ".*"
	}
	if len(m["fn"]) > 0 {
		fname = m["fn"][0]
	} else {
		fname = ".*"
	}
	if len(m["ln"]) > 0 {
		lname = m["ln"][0]
	} else {
		lname = ".*"
	}

	return uid, uname, fname, lname, nil
}

/**
 *  handles admin request for users
**/
func handleListUsers(w http.ResponseWriter, r *http.Request) {
	q := `Select user_id, username, first_name, last_name, isTOTPSetup, isAdmin
            From users
            where username rLike ? AND
                    first_name rLike ? AND
                    last_name rLike ? AND
                    user_id rLike ?`
	uid, uname, fname, lname, err := getQuery(w, r)
	if err != nil {
		return
	}

	var usrs []UserOut
	err = db.Select(&usrs, q, uname, fname, lname, uid)
	if err != nil {
		sendError(w, "DB: Error Could not retrive users", err, http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(usrs)

}

/**
 * Handles user deletion either by the user them selves or the admin
 **/
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uId"])
	if err != nil {
		sendError(w, "Error convert string to int", err, http.StatusInternalServerError)
		return
	}

	q := `Delete From users where user_id = ?`

	_, err = db.Exec(q, uid)
	if err != nil {
		sendError(w, "Error could not delete user", err, http.StatusInternalServerError)
		return
	}

}

func handleAdminPwdReset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.Atoi(vars["uId"])
	if err != nil {
		sendError(w, "Error convert string to int", err, http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var u User
	err = decoder.Decode(&u)
	if err != nil {
		sendError(w, "Error parsing password reset JSON", err, http.StatusBadRequest)
		return
	}

	resetPwd(w, u.Password, uid)
}
