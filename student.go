package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

//Gets a class object
func handleGetStudentClasses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classCode, _ := vars["cCode"]
	q := `SELECT class_code, class_name FROM classes WHERE class_code = ?`
	class := StudentClass{}
	nilClass := StudentClass{}
	dbErr := db.Get(&class, q, classCode)
	if deepEqualCheck(w, class, nilClass) {
		return
	}
	structPass(class, w, dbErr)
}

//Checks to see if who interfaces are the same
func deepEqualCheck(w http.ResponseWriter, firstStruct interface{}, secondStruct interface{}) bool {
	if reflect.DeepEqual(firstStruct, secondStruct) {
		nrErr := errors.New("No rows returned from SQL query")
		sendError(w, "No rows", nrErr, http.StatusNotFound)
		return true
	}
	return false
}

//Gets a list of questions for a class
func handleGetClassQuestions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classCode, _ := vars["cCode"]
	q := `SELECT question_id, class_code, question_text, question_image
	FROM questions WHERE class_code = ? ORDER BY question_id ASC`
	questList := []GetQuestion{}
	nilQuest := []GetQuestion{}
	dbErr := db.Select(&questList, q, classCode)
	if deepEqualCheck(w, questList, nilQuest) {
		return
	}
	structPass(questList, w, dbErr)

}

//Gets a list of answers for a question
func handleGetQuestionAns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questNum, err := strconv.Atoi(vars["qId"])
	if err != nil {
		sendError(w, "Vars error", err, http.StatusBadRequest)
		return
	}
	q := `SELECT answer_id, question_id, answer_text, answer_image, isCorrect FROM answers
	WHERE question_id = ? ORDER BY answer_id ASC`
	answerList := []Answer{}
	nilAnswer := []Answer{}
	dbErr := db.Select(&answerList, q, questNum)
	if deepEqualCheck(w, answerList, nilAnswer) {
		return
	}
	structPass(answerList, w, dbErr)
}

//Creates a response to a question
func handleCreateResp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questNum, qErr := strconv.Atoi(vars["qId"])
	answerNum, aErr := strconv.Atoi(vars["aId"])
	if aErr != nil || qErr != nil {
		sendError(w, "Vars error", qErr, http.StatusBadRequest)
		return
	}
	respFromDB := Response{}
	q := `INSERT INTO responses (question_id, answer_id)
				VALUES (?, ?)`
	qs := `SELECT response_id, question_id, answer_id FROM
	responses WHERE response_id = ? ORDER BY response_id ASC`
	tx, txErr := db.Beginx()
	if txErr != nil {
		sendError(w, "Transaction error", txErr, http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	txInsResults, txInsErr := tx.Exec(q, questNum, answerNum)
	if txInsErr != nil {
		sendError(w, "Transaction error", txInsErr, http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	lastID, _ := txInsResults.LastInsertId()
	txSelErr := tx.Get(&respFromDB, qs, lastID)
	if txSelErr != nil {
		sendError(w, "Transaction error", txSelErr, http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	txComErr := tx.Commit()
	structPass(respFromDB, w, txComErr)
}

//Updates a response to a question
func handleUpdateResp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	respNum, rErr := strconv.Atoi(vars["rId"])
	answerNum, aErr := strconv.Atoi(vars["aId"])
	if aErr != nil || rErr != nil {
		sendError(w, "Vars error", rErr, http.StatusBadRequest)
		return
	}
	q := `UPDATE responses SET answer_id = ? WHERE response_id = ?`
	_, dbErr := db.Exec(q, answerNum, respNum)
	if dbErr != nil {
		sendError(w, "Database error", dbErr, http.StatusInternalServerError)
		return
	}
}

//structPass takes an interface and outputs it as a JSON response
//will also output error that is passed, if exists
func structPass(curStruct interface{}, w http.ResponseWriter, dbErr error) {
	if dbErr == sql.ErrNoRows {
		sendError(w, "No rows found", dbErr, http.StatusNotFound)
		return
	} else if dbErr != nil {
		sendError(w, "Database error", dbErr, http.StatusInternalServerError)
		return
	}
	jsonResp, jErr := json.Marshal(curStruct)
	if jErr != nil {
		sendError(w, "Json error", jErr, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
