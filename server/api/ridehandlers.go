package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	db "server/database"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RequestRideHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RiderID     int64   `json:"rider_id"`
		PickupLat   float64 `json:"pickup_lat"`
		PickupLong  float64 `json:"pickup_long"`
		DropoffLat  float64 `json:"dropoff_lat"`
		DropoffLong float64 `json:"dropoff_long"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	rideID, err := s.DB.CreateRide(r.Context(), db.CreateRideParams{
		RiderID:     req.RiderID,
		PickupLat:   req.PickupLat,
		PickupLong:  req.PickupLong,
		DropoffLat:  req.DropoffLat,
		DropoffLong: req.DropoffLong,
	})
	if err != nil {
		http.Error(w, "failed to create ride", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ride_id": rideID,
		"status":  "requested",
	})
}

func (s *Server) GetNearbyDriversHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	longStr := r.URL.Query().Get("long")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}
	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}
	params := db.GetNearbyDriversParams{
		Lat:         sql.NullFloat64{Float64: lat, Valid: true},
		LlToEarth:   fmt.Sprintf("ll_to_earth(%f,%f)", lat, long),
		LlToEarth_2: fmt.Sprintf("ll_to_earth(%f,%f)", lat, long),
	}

	drivers, err := s.DB.GetNearbyDrivers(r.Context(), params)
	if err != nil {
		http.Error(w, "failed to fetch nearby drivers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(drivers)
}

func (s *Server) GetRideStatusHandler(w http.ResponseWriter, r *http.Request) {
	rideIDStr := chi.URLParam(r, "id")
	rideID, err := strconv.ParseInt(rideIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid ride ID", http.StatusBadRequest)
		return
	}

	ride, err := s.DB.GetRideByID(r.Context(), rideID)
	if err != nil {
		http.Error(w, "ride not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ride_id":   ride.ID,
		"status":    ride.Status,
		"driver_id": ride.DriverID,
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

	err := s.Redis.SetDriverLocation(r.Context(), uint64(body.DriverID), body.Latitude, body.Longitude)
	if err != nil {
		http.Error(w, "failed to update location", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "location updated successfully",
	})
}

func (s *Server) AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	completedRides, err := s.DB.GetAnalytics(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch analytics", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"average_wait_time_minutes": completedRides, // TODO: both
		"completed_rides":           completedRides,
	})
}
