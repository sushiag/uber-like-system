package api

import (
	"encoding/json"
	"log"
	"net/http"
	db "server/database"
	redis "server/redis"
	ws "server/ws"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	DB    *db.Queries
	Redis *redis.Client
	Wsm   *ws.WebSocketManager
}

func (s *Server) RegisterRoute(r chi.Router) {
	r.Post("/riders/register", s.createRider)
	r.Post("/drivers/register", s.createDriver)
	r.Post("/loginRider", s.LoginRider)
	r.Post("/loginDriver", s.LoginDriver)
	// r.Post("/rides/request", s.createDriverRequest)
	r.Post("/drivers/{id}/location", s.updateDriverLocation)
	//	r.Get("/rides/{id}/status", s.getRideStatus)
}

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

func (s *Server) updateDriverLocation(w http.ResponseWriter, r *http.Request) {
	type req struct {
		DriverID  int64   `json:"driver_id"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"long"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//	err := s.Redis.HSet(r.Context(),
	//		fmt.Sprintf("%d", body.DriverID),
	//		"lat", body.Latitude,
	//		"long", body.Longitude,
	//	).Err()
	//	if err != nil {
	//		http.Error(w, "failed to update location", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	w.WriteHeader(http.StatusOK)
}

//func (s *Server) getRideStatus(w http.ResponseWriter, r *http.Request) {
//	rideID := chi.URLParam(r, "id")
//
//	var status string
//	err := s.DB.QueryRow("SELECT status FROM ride_requests WHERE id=$1", rideID).Scan(&status)
//	if err != nil {
//		http.Error(w, "ride not found", http.StatusNotFound)
//		return
//	}
//
//	json.NewEncoder(w).Encode(map[string]string{"status": status})
//}
