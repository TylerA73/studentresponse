package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("'%v' to [%v]", r.Method, r.URL.Path)
		f.ServeHTTP(w, r)
	})

}

func createRouter() *mux.Router {
	// Site routes.
	// TODO: If we're doing a SPA, we should only be returning Static docs.
	// For now, I'll leave this be, and add an API subrouter.
	r := mux.NewRouter()
	r.Use(logging)

	apiRoutes(r)

	//Serve files
	r.Handle("/socket.io/", startSocket())
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("htdocs/"))))
	return r
}

func apiRoutes(r *mux.Router) {
	// API Routes.(No Authorization needed)
	r.HandleFunc("/api/v1/login", handleLogin).Methods("POST")
	r.HandleFunc("/api/v1/logout", handleLogout).Methods("POST")
	r.HandleFunc("/api/v1/register", handleRegister).Methods("POST")
	r.HandleFunc("/api/v1/admins/login", handleAdminLogin).Methods("POST")

	studentApi(r)

	// API Routes (Authorization required)
	s := r.PathPrefix("/api/v1").Subrouter()
	s.Use(authMiddleware)
	adminApi(s)
	userApi(s)
	teacherApi(s)
	twoFactorApi(s)
}

func twoFactorApi(r *mux.Router) {
	r.HandleFunc("/2fa/qr", handleTotpQR).Methods("GET")
	r.HandleFunc("/2fa/challenge", handleTotpChallenge).Methods("POST")
	r.HandleFunc("/2fa/pass", handleTotpPassOption).Methods("POST")
}

func userApi(r *mux.Router) {
	u := r.PathPrefix("/users").Subrouter()
	u.HandleFunc("/passwords", handlePwdReset).Methods("PUT")
	u.HandleFunc("", handleAccountDelete).Methods("DELETE")
}

func teacherApi(r *mux.Router) {
	t := r.PathPrefix("/teachers").Subrouter()
	t.HandleFunc("/classes", handleGetClasses).Methods("GET")
	t.HandleFunc("/classes", handleMakeNewClass).Methods("POST")
	t.HandleFunc("/classes/{code}", handleDeleteClass).Methods("Delete")

	t.HandleFunc("/classes/{code}", handleGetTeacherQuestions).Methods("GET")
	t.HandleFunc("/classes/{code}", handleCreateTeacherQuestion).Methods("POST")
	t.HandleFunc("/classes/{code}/qrjoin", handleQRCodeGenerateForClass).Methods("GET")
	t.HandleFunc("/questions/{qId:[0-9]+}", handleUpdateTeacherQuestion).Methods("PUT")
	t.HandleFunc("/questions/{qId:[0-9]+}", handleDeleteTeacherQuestion).Methods("Delete")
	t.HandleFunc("/questions/{qId:[0-9]+}", handleGetQuestionStats).Methods("GET")
}

func adminApi(r *mux.Router) {
	a := r.PathPrefix("/admins").Subrouter()

	a.HandleFunc("/users", handleListUsers).Methods("GET")
	a.HandleFunc("/passwords/{uId:[0-9]+}", handleAdminPwdReset).Methods("PUT")
	a.HandleFunc("/users/{uId:[0-9]+}", handleDeleteUser).Methods("DELETE")
}

func studentApi(r *mux.Router) {
	r.HandleFunc("/api/v1/classes/{cCode:[0-9A-Za-z]+}", handleGetStudentClasses).Methods("GET")
	r.HandleFunc("/api/v1/classes/{cCode:[0-9A-Za-z]+}/questions", handleGetClassQuestions).Methods("GET")
	r.HandleFunc("/api/v1/questions/{qId:[0-9]+}/answers", handleGetQuestionAns).Methods("GET")
	r.HandleFunc("/api/v1/questions/{qId:[0-9]+}/answers/{aId:[0-9]+}", handleCreateResp).Methods("POST")
	r.HandleFunc("/api/v1/responses/{rId:[0-9]+}/answers/{aId:[0-9]+}", handleUpdateResp).Methods("PUT")
}
