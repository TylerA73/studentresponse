package main

import "database/sql"

type HttpError struct {
	ErrorStr string `json:"error"`
}

/**
 * Database Models start here.
**/
type User struct {
	ID               int            `db:"user_id" json:"userid"`
	Username         string         `db:"username" json:"username"`
	FirstName        string         `db:"first_name" json:"firstname"`
	LastName         string         `db:"last_name" json:"lastname"`
	Password         string         `db:"password" json:"password"`
	IsTOTPSetup      bool           `db:"isTOTPSetup" json:"istotpsetup"`
	TOTPSecret       sql.NullString `db:"totp_secret" `
	TOTPRecoveryCode sql.NullString `db:"totp_recovery_code"`
	IsAdmin          bool           `db:"isAdmin" json:"isadmin"`
}

/**
 * struct to send out to admin, cannot contain sensitve data
**/
type UserOut struct {
	ID          int    `db:"user_id" json:"userid"`
	Username    string `db:"username" json:"username"`
	FirstName   string `db:"first_name" json:"firstname"`
	LastName    string `db:"last_name" json:"lastname"`
	IsTOTPSetup bool   `db:"isTOTPSetup" json:"istotpsetup"`
	IsAdmin     bool   `db:"isAdmin" json:"isadmin"`
}

//Question struct used to retrieve list of questions for class
type Question struct {
	QuestionID    int            `db:"question_id" json:"questionId"`
	ClassCode     string         `db:"class_code" json:"classcode"`
	QuestionText  string         `db:"question_text" json:"questionText"`
	QuestionImage sql.NullString `db:"question_image" json:"questionImage"`
	ResponseCount int            `json:"count"`
}

//GetQuestion struct used to hold question data
type GetQuestion struct {
	QuestionID    int            `db:"question_id" json:"questionId"`
	QuestionText  string         `db:"question_text" json:"questionText"`
	ClassCode     string         `db:"class_code" json:"classCode"`
	QuestionImage sql.NullString `db:"question_image" json:"questionImage"`
}

type QuestionAnswer struct {
	QuestionInfo Question `json:"question"`
	Answers      []Answer `json:"answers"`
}

type Answer struct {
	AnswerId      int            `db:"answer_id" json:"answerId"`
	QuestionID    int            `db:"question_id" json:"questionId"`
	AnswerText    string         `db:"answer_text" json:"answerText"`
	AnswerImage   sql.NullString `db:"answer_image" json:"answerImage"`
	IsCorrect     bool           `db:"isCorrect" json:"iscorrect"`
	ResponseCount int            `json:"count"`
}

type Class struct {
	UserId    int    `db:"user_id" json:"userid"`
	ClassName string `db:"class_name" json:"classname"`
	ClassCode string `db:"class_code" json:"classcode"`
}

//StudentClass gets class info for student
type StudentClass struct {
	ClassCode string `db:"class_code" json:"classcode"`
	ClassName string `db:"class_name" json:"classname"`
}

//Response struct used to hold response data
type Response struct {
	ResponseID int `db:"response_id" json:"responseId"`
	QuestionID int `db:"question_id" json:"questionId"`
	AnswerID   int `db:"answer_id" json:"answerId"`
}
