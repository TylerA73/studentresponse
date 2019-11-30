package main

import (
	"crypto/rand"
	"encoding/json"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/lhecker/argon2"
)

func getUID(r *http.Request) int {
	ctx := r.Context()
	var userID = int(ctx.Value("UserID").(float64))
	return userID
}

func encodePwd(newPwd string) ([]byte, error) {
	cfg := argon2.DefaultConfig()
	encoded, err := cfg.HashEncoded([]byte(newPwd))
	if err != nil {
		return []byte(""), err
	}
	return encoded, nil
}

func sendError(w http.ResponseWriter, errStrm string, err error, hcode int) {
	logger.Printf("%s: %v\n", errStrm, err.Error())
	w.WriteHeader(hcode)

	returnedError := HttpError{ErrorStr: errStrm}
	encoder := json.NewEncoder(w)
	encoder.Encode(returnedError)
}

func resetPwd(w http.ResponseWriter, pwd string, uid int) {
	encoded, err := encodePwd(pwd)
	if err != nil {
		sendError(w, "Error encoding password", err, http.StatusInternalServerError)
		return
	}

	q := `update users
			set password = ?
			where user_id = ?
			limit 1`
	_, err = db.Exec(q, encoded, uid)
	if err != nil {
		sendError(w, "DB: Error Could not reset user password", err, http.StatusInternalServerError)
		return
	}
}

func getClassCode(n int) (string, error) {
	tmp, err := GenerateRandomString(n)
	if err != nil {
		return "", err
	}
	var c []Class
	q := `Select *
			From classes
			where class_code = ?`
	err = db.Select(&c, q, tmp)
	if err != nil {
		return "", err
	}
	if len(c) > 0 {
		return "exists", nil
	}
	return tmp, nil

}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// Taken from:
// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// Taken from:
// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

/**
 * isDuplicateRowError(error err)
 * Check if an error generated from a MySQL call is a unique row constraint violation.
 * Taken from Nick's code presented in CMPT 315.
**/
func isDuplicateRowError(err error) bool {
	mysqlerr := err.(*mysql.MySQLError)
	return mysqlerr.Number == 1062 // Duplicate Row Error Code. (MySQL Protocol)
}
