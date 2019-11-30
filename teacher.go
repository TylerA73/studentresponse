package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	qrcode "github.com/skip2/go-qrcode"
)

func handleGetClasses(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	q := `Select class_name, class_code
			From srs.classes
			where user_id = ?`
	var c []Class
	err := db.Select(&c, q, uid)
	if err != nil {
		sendError(w, "DB: Error Could not retrive list of class", err, http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(c)
}

func handleMakeNewClass(w http.ResponseWriter, r *http.Request) {
	var c Class

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		sendError(w, "Error parsing Registration JSON", err, http.StatusBadRequest)
		return
	}

	c.UserId = getUID(r)
	code, err := getClassCode(4)
	if err != nil {
		sendError(w, "DB: Error occured getting class code", err, http.StatusInternalServerError)
		return
	}
	for code == "exists" {
		code, err = getClassCode(4)
		if err != nil {
			sendError(w, "DB: Error occured getting class code", err, http.StatusInternalServerError)
			return
		}
	}
	c.ClassCode = code
	q := `INSERT INTO
			classes (user_id, class_name, class_code)
			values (?, ?, ?)`
	_, err = db.Exec(q, c.UserId, c.ClassName, c.ClassCode)
	if err != nil {
		sendError(w, "DB: Error Could not create class", err, http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(c.ClassCode)
}

func handleDeleteClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["qId"])
	if err != nil {
		sendError(w, "Error convert string to int", err, http.StatusInternalServerError)
		return
	}
	q := `Delete From questions where question_id = ?`

	_, err = db.Exec(q, qid)
	if err != nil {
		sendError(w, "Error could not delete question", err, http.StatusInternalServerError)
		return
	}
}

func handleGetTeacherQuestions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	var qs []Question
	q := `SELECT question_id, class_code, question_text
			FROM questions
			where class_code = ?`

	err := db.Select(&qs, q, code)
	if err != nil {
		sendError(w, "Error could not select list of questions for class", err, http.StatusInternalServerError)
		return
	}
	for i, ques := range qs {
		count, errStr, err := getQStat(ques.QuestionID)
		if err != nil {
			sendError(w, errStr, err, http.StatusInternalServerError)
			return
		}
		qs[i].ResponseCount = count
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(qs)
}

func getQStat(qid int) (int, string, error) {
	q := `SELECT count(*)
		FROM responses
		where question_id = ?`
	var c int
	err := db.Get(&c, q, qid)
	if err != nil {
		return 0, "Error could not count of resposnses for question for class", err
	}
	return c, "", nil
}

func handleCreateTeacherQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	var newQs QuestionAnswer

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newQs)
	if err != nil {
		sendError(w, "Error parsing Registration JSON", err, http.StatusBadRequest)
		return
	}
	newQs.QuestionInfo.ClassCode = code

	err = insertQuestions(newQs)
	if err != nil {
		sendError(w, "Error could not create question", err, http.StatusInternalServerError)
		return
	}
}

func handleQRCodeGenerateForClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	var png []byte
	png, err := qrcode.Encode(r.Host+"/joinedclass.html?code="+code, qrcode.Medium, 300)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func insertQuestions(nqs QuestionAnswer) error {
	qs := nqs.QuestionInfo
	q := `insert into 
			questions (class_code, question_text)
			values (?, ?)`
	_, err := db.Exec(q, qs.ClassCode, qs.QuestionText)
	if err != nil {
		return nil
	}

	q = `SELECT LAST_INSERT_ID()`
	var qid int
	err = db.Get(&qid, q)
	if err != nil {
		return err
	}

	err = insertAnswers(qid, nqs.Answers)
	return nil
}

func insertAnswers(qid int, a []Answer) error {
	for _, ans := range a {
		q := `insert  into
		answers (question_id, answer_text, isCorrect)
				values (?, ?, ?)`
		_, err := db.Exec(q, qid, ans.AnswerText, ans.IsCorrect)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleDeleteTeacherQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid := vars["qId"]
	q := `Delete From questions where question_id = ?`

	_, err := db.Exec(q, qid)
	if err != nil {
		sendError(w, "Error could not delete class", err, http.StatusInternalServerError)
		return
	}
}

func handleUpdateTeacherQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["qId"])
	if err != nil {
		sendError(w, "Error convert string to int", err, http.StatusInternalServerError)
		return
	}

	var qs QuestionAnswer

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&qs)
	if err != nil {
		sendError(w, "Error parsing Registration JSON", err, http.StatusBadRequest)
		return
	}

	q := `update questions 
			set  question_text = ?
			where question_id = ?`
	_, err = db.Exec(q, qs.QuestionInfo.QuestionText, qid)
	if err != nil {
		sendError(w, "Error could not update question", err, http.StatusInternalServerError)
		return
	}

	err = updateAnswers(qid, qs.Answers)
	if err != nil {
		sendError(w, "Error could not update answers", err, http.StatusInternalServerError)
		return
	}
}

func updateAnswers(qid int, a []Answer) error {
	q := `SELECT answer_id FROM answers
			where question_id = ?`
	var aID []int
	err := db.Select(&aID, q, qid)
	if err != nil {
		return err
	}

	var presentAID []int

	for _, ans := range a {
		q = `update answers
				set question_id = ?,
					answer_text = ?,
					isCorrect =?
				where answer_id = ?`
		_, err = db.Exec(q, qid, ans.AnswerText, ans.IsCorrect, ans.AnswerId)
		if err != nil {
			return err
		}
		for _, a := range aID {
			if a == ans.AnswerId {
				presentAID = append(presentAID, a)
			}
		}
	}

	for _, a := range aID {
		if !present(a, presentAID) {
			q = `DELETE FROM answers
					WHERE answer_id = ?`
			_, err = db.Exec(q, a)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func present(a int, presentAID []int) bool {
	for _, i := range presentAID {
		if a == i {
			return true
		}
	}
	return false
}

func handleGetQuestionStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qid, err := strconv.Atoi(vars["qId"])
	if err != nil {
		sendError(w, "Error convert string to int", err, http.StatusInternalServerError)
		return
	}
	qcount, errStr, err := getQStat(qid)
	if err != nil {
		sendError(w, errStr, err, http.StatusInternalServerError)
	}

	var qa QuestionAnswer

	q := `SELECT *
			FROM questions
			WHERE question_id = ?`
	err = db.QueryRowx(q, qid).StructScan(&qa.QuestionInfo)
	if err != nil {
		sendError(w, "Error fetching question details, getting question", err, http.StatusInternalServerError)
		return
	}
	qa.QuestionInfo.ResponseCount = qcount

	q = `SELECT *
			FROM answers
			where question_id = ?`
	err = db.Select(&qa.Answers, q, qid)

	if err != nil {
		sendError(w, "Error fetching question details, getting list of answers", err, http.StatusInternalServerError)
		return
	}

	for i, ans := range qa.Answers {
		q := `SELECT count(*)
				FROM responses
				where answer_id = ?`
		var c int
		err := db.Get(&c, q, ans.AnswerId)
		if err != nil {
			sendError(w, "Error fetching question details, getting counts per answer", err, http.StatusInternalServerError)
			return
		}
		qa.Answers[i].ResponseCount = c
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(qa)
}
