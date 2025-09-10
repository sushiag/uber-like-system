package api

import (
	"encoding/json"
	"log"
	"net/http"
	db "server/database"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) createRider(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	log.Printf("[Rider Registration] Request to create an account!")

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("[Rider Registration] Received request to create an account!")

	// Username
	if err := UsernameField(body.Username); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	// Password
	if err := PasswordField(body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Password Hashed for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = s.DB.CreateRider(r.Context(), db.CreateRiderParams{
		Username: body.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "failed to insert rider to database", http.StatusInternalServerError)
		return
	}
	log.Printf("[Rider Registration] Successfully registered rider %s", body.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "you've been registered successfully!",
	})
}

func (s *Server) createDriver(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	log.Printf("[Driver Registration] Request to create an account!")

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("[Driver Registration] Received request to create an account!")

	// Username
	if err := UsernameField(body.Username); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	// Password
	if err := PasswordField(body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Password Hashed for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = s.DB.CreateDriver(r.Context(), db.CreateDriverParams{
		Username: body.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "failed to insert driver to database", http.StatusInternalServerError)
		return
	}
	log.Printf("[Driver Registration] Successfully registered rider %s", body.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "you've been registered successfully!",
	})
}

func (s *Server) LoginRider(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// If username/password is provided → normal login
	user, err := s.DB.GetRiderByUsername(r.Context(), body.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Login successful
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "You've Login successful!",
		"username": user.Username,
		"password": user.Password,
	})
}

func (s *Server) LoginDriver(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// If username/password is provided → normal login
	user, err := s.DB.GetDriverByUsername(r.Context(), body.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Login successful
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "You've Login successful!",
		"username": user.Username,
		"password": user.Password,
	})
}
