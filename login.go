package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lhecker/argon2"
	"github.com/pquerna/otp/totp"
)

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	login(w, r, true)
}

/**
 *  API Route: /api/v1/login
 *  Description: Handle login attempts.
**/
func handleLogin(w http.ResponseWriter, r *http.Request) {
	login(w, r, false)
}

func login(w http.ResponseWriter, r *http.Request, admin bool) {
	decoder := json.NewDecoder(r.Body)
	var attemptedUser User

	err := decoder.Decode(&attemptedUser)
	if err != nil {
		logger.Printf("%s: %v\n", "Error parsing Login JSON", err.Error())
		// Tried to parse malformed JSON.
		w.WriteHeader(http.StatusBadRequest)

		returnedError := HttpError{ErrorStr: "Malformed JSON Detected"}
		encoder := json.NewEncoder(w)
		encoder.Encode(returnedError)

		return
	}
	// fmt.Printf("%v\n", attemptedUser)

	var knownUser User
	err = db.Get(&knownUser, "SELECT * FROM `users` WHERE UPPER(username) LIKE UPPER(?) LIMIT 1;", attemptedUser.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// Invalid Username.
			w.WriteHeader(http.StatusUnauthorized)
			returnedError := HttpError{ErrorStr: "Invalid Username or Password"}
			encoder := json.NewEncoder(w)
			encoder.Encode(returnedError)
		} else {
			dbLogger.Printf("%s: %s\n", "Error selecting user from Users", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Check password validity using Argon2.
	var passValid bool
	passValid, err = argon2.VerifyEncoded([]byte(attemptedUser.Password), []byte(knownUser.Password))
	if err != nil {
		// Error with CGO bindings for Argon2
		logger.Printf("%s: %s\n", "Error comparing passwords using argon2 bindings", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !passValid {
		// Invalid Password.
		w.WriteHeader(http.StatusUnauthorized)
		returnedError := HttpError{ErrorStr: "Invalid Username or Password"}
		encoder := json.NewEncoder(w)
		encoder.Encode(returnedError)
		return
	}

	if admin && !knownUser.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		returnedError := HttpError{ErrorStr: "Invalid Username or Password"}
		encoder := json.NewEncoder(w)
		encoder.Encode(returnedError)
		return
	}

	// At this point the user is authenticated. Let's tell them.
	sessionExpire := time.Now().Add(2 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": knownUser.ID,
		"exp": sessionExpire.Unix(),
	})

	var tokenstr string
	tokenstr, err = token.SignedString([]byte(os.Getenv("SESSION_SIGNKEY")))
	if err != nil {
		// Error signing JWT.
		logger.Printf("%s: %s\n", "Error signing JWT", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cookie = http.Cookie{
		Name:    "srs",
		Value:   tokenstr,
		Expires: sessionExpire,
	}
	http.SetCookie(w, &cookie)

	// Add Token to the value set at Key:UserID.
	// Set / Push expiry of Key until given time.
	sessionStore.SAdd(tokenstr, strconv.Itoa(knownUser.ID))
	sessionStore.ExpireAt(tokenstr, sessionExpire)

	if knownUser.IsTOTPSetup {
		w.Write([]byte("{\"2fa\":\"challenge\"}"))
	} else {
		w.Write([]byte("{\"2fa\":\"setup\"}"))
	}
}

/**
 * Login Validator Middleware
**/
func authMiddleware(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := r.Cookie("srs")
		if err != nil {
			// Return unauthorized, as cookie wasn't found.
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{\"error\":\"Unauthorized\"}"))
			return
		}

		token, err := jwt.Parse(authToken.Value, func(token *jwt.Token) (interface{}, error) {
			// TODO: Verify token method.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: [%v]", token.Header["alg"])
			}

			return []byte(os.Getenv("SESSION_SIGNKEY")), nil
		})

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{\"error\":\"Unauthorized\"}"))
			return
		}

		var sessionValid int64
		sessionValid, err = sessionStore.Exists(authToken.Value).Result()
		if err != nil {
			sessionLog.Println("Error checking session store: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var twofactorValid int64
		twofactorValid, err = sessionStore.Exists(authToken.Value + ".2fa").Result()
		if err != nil {
			sessionLog.Println("Error checking session store: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Key does not exist in Redis. Session has been invalidated server side.
		if sessionValid == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{\"error\":\"Unauthorized\"}"))
			return
		}

		if !strings.Contains(r.RequestURI, "/api/v1/2fa") && twofactorValid == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{\"error\":\"Unauthorized\"}"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{\"error\":\"Unauthorized\"}"))
		} else {
			ctx := context.WithValue(r.Context(), "UserID", claims["sub"])
			f.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func handleTotpQR(w http.ResponseWriter, r *http.Request) {
	userid := getUID(r)

	// 1. Check user doesn't already have TOTP setup.
	var user User
	err := db.Get(&user, "SELECT * FROM `users` WHERE user_id = ? LIMIT 1;", userid)
	if err != nil {
		// User doesn't exist. 500.
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println("Error obtaining user record: ", err.Error())
		return
	}

	if user.IsTOTPSetup {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Two Factor Authentication is already enabled.\"}"))
		return
	}

	// 2. Generate TOTP Secret Key.
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "Student Response System",
		AccountName: user.Username,
	})

	// 3. Generate Recovery Key randomly.
	var recoveryKey string
	recoveryKey, err = GenerateRandomString(24)

	// 4. Save secret key and recovery token to DB for this user.
	_, err = db.Exec("UPDATE `users` SET totp_secret = ?, totp_recovery_code = ? WHERE user_id = ? LIMIT 1", key.Secret(), recoveryKey, userid)
	if err != nil {
		// Error writing to DB.
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println("Error saving TOTP secret: ", err.Error())
		return
	}

	// 4. Convert TOTP Key to QR Code as PNG.
	// Thanks to https://godoc.org/github.com/pquerna/otp for the example usage.
	var outbuf bytes.Buffer
	img, err := key.Image(300, 300)
	png.Encode(&outbuf, img)

	// Write the PNG out.
	w.Header().Set("Content-Type", "image/png")
	w.Write(outbuf.Bytes())
}

func handleTotpChallenge(w http.ResponseWriter, r *http.Request) {
	userid := getUID(r)

	// 1. Retrieve user data.
	var user User
	err := db.Get(&user, "SELECT * FROM `users` WHERE user_id = ? LIMIT 1", userid)
	if err != nil {
		// User doesn't exist.
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println("Error obtaining user record: ", err.Error())
	}

	// 2. Check if totp is setup.
	if !user.TOTPSecret.Valid {
		w.WriteHeader(http.StatusFailedDependency)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\":\"TOTP Keys not assigned for User.\"}"))
		return
	}

	var authAttempt struct {
		Passcode string `json:"code"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&authAttempt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\":\"Malformed JSON received.\"}"))
		return
	}

	totpRes := totp.Validate(authAttempt.Passcode, user.TOTPSecret.String)
	totpRecover := (authAttempt.Passcode == user.TOTPRecoveryCode.String)
	if !totpRes && !totpRecover {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authToken, err := r.Cookie("srs")
	if err != nil {
		// Return unauthorized, as cookie wasn't found.
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"error\":\"Unauthorized\"}"))
		return
	}
	sessionStore.Set(authToken.Value+".2fa", "Valid", 2*time.Hour)

	// Authorized at this point, let's make sure the DB knows TOTP is confirmed setup.
	if !user.IsTOTPSetup {
		_, err = db.Exec("UPDATE `users` SET isTOTPSetup = 1 WHERE user_id = ? LIMIT 1", user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Println("Error setting isTOTPSetup value: ", err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"recovery_code\":\"" + user.TOTPRecoveryCode.String + "\"}"))
	}
	if totpRecover {
		_, err = db.Exec("UPDATE `users` SET isTOTPSetup = 0, totp_recovery_code = NULL, totp_secret = NULL WHERE user_id = ? LIMIT 1", user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Println("Error resetting 2FA state: ", err.Error())
		}
	}
}

func handleTotpPassOption(w http.ResponseWriter, r *http.Request) {
	userid := getUID(r)
	authToken, err := r.Cookie("srs")
	if err != nil {
		// Return unauthorized, as cookie wasn't found.
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"error\":\"Unauthorized\"}"))
		return
	}

	// 1. Retrieve user data.
	var user User
	err = db.Get(&user, "SELECT * FROM `users` WHERE user_id = ? LIMIT 1", userid)
	if err != nil {
		// User doesn't exist.
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println("Error obtaining user record: ", err.Error())
	}

	// 2. Check if totp is setup.
	if !user.IsTOTPSetup {
		// Since totp isn't setup, user has the option to pass on 2fa enrollment.
		// Grant a pass for this session.
		sessionStore.Set(authToken.Value+".2fa", "Valid", 2*time.Hour)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	authToken, err := r.Cookie("srs")
	if err != nil {
		// Return unauthorized, as cookie wasn't found.
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"error\":\"Unauthorized\"}"))
	}

	if authToken.Value == "" {
		return
	}

	// Expire active session keys.
	sessionStore.Del(authToken.Value, authToken.Value+".2fa")
}
