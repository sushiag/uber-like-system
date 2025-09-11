package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	db "uber-like-system/server/database"

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

	params := db.GetNearbyDriversParams{
		LlToEarth:   fmt.Sprintf("ll_to_earth(%f,%f)", req.PickupLat, req.PickupLong),
		LlToEarth_2: fmt.Sprintf("ll_to_earth(%f,%f)", req.PickupLat, req.PickupLong),
		Lat:         sql.NullFloat64{Float64: req.PickupLat, Valid: true},
	}

	drivers, err := s.DB.GetNearbyDrivers(r.Context(), params)
	if err != nil {
		log.Println("failed to find nearby drivers:", err)
		http.Error(w, "failed to find drivers", http.StatusInternalServerError)
		return
	}
	if len(drivers) == 0 {
		http.Error(w, "no drivers available nearby", http.StatusNotFound)
		return
	}

	driver := drivers[0]

	ride, err := s.DB.CreateRide(r.Context(), db.CreateRideParams{
		RiderID:     req.RiderID,
		DriverID:    sql.NullInt64{Valid: false},
		PickupLat:   req.PickupLat,
		PickupLong:  req.PickupLong,
		DropoffLat:  req.DropoffLat,
		DropoffLong: req.DropoffLong,
	})

	if err != nil {
		log.Println("failed to create ride:", err)
		http.Error(w, "failed to create ride", http.StatusInternalServerError)
		return
	}

	err = s.DB.AssignDriverToRide(r.Context(), db.AssignDriverToRideParams{
		ID:       ride.ID,
		DriverID: sql.NullInt64{Int64: driver.DriverID, Valid: true},
	})
	if err != nil {
		log.Println("failed to assign driver:", err)
		http.Error(w, "failed to assign driver", http.StatusInternalServerError)
		return
	}

	s.Wsm.SendToUser(uint64(driver.DriverID), []byte(fmt.Sprintf(
		`{"event":"ride_assigned","ride_id":%d,"driver_id":%d}`,
		ride.ID, driver.DriverID,
	)))

	resp := map[string]interface{}{
		"ride_id": ride.ID,
		"status":  "driver_assigned",
		"driver": map[string]interface{}{
			"driver_id": driver.DriverID,
			"username":  driver.Username,
		},
	}
	json.NewEncoder(w).Encode(resp)
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

	driverIDs, err := s.Redis.GetNearbyDrivers(r.Context(), lat, long, 5000)
	if err != nil {
		http.Error(w, "failed to fetch nearby drivers", http.StatusInternalServerError)
		log.Println("Redis error:", err)
		return
	}

	if len(driverIDs) == 0 {
		http.Error(w, "no nearby drivers", http.StatusNotFound)
		return
	}

	drivers := []map[string]interface{}{}
	for _, driverIDStr := range driverIDs {
		driverID, err := strconv.ParseInt(driverIDStr, 10, 64) // convert string -> int64
		if err != nil {
			log.Println("invalid driver ID from Redis:", driverIDStr)

			continue
		}

		driver, err := s.DB.GetDriverByID(r.Context(), strconv.FormatInt(driverID, 10))
		if err != nil {
			log.Println("driver not found in DB:", driverID)
			continue
		}

		drivers = append(drivers, map[string]interface{}{
			"driver_id": driver.ID,
			"username":  driver.Username,
		})
	}

	log.Println("Drivers returned:", drivers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"drivers": drivers,
	})

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

func (s *Server) AcceptRideHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DriverID int64 `json:"driver_id"`
		RideID   int64 `json:"ride_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	params := db.AssignDriverToRideParams{
		ID:       req.RideID,
		DriverID: sql.NullInt64{Int64: req.DriverID, Valid: true},
	}

	err := s.DB.AssignDriverToRide(r.Context(), params)
	if err != nil {
		log.Println("failed to assign driver:", err)
		http.Error(w, "failed to assign driver", http.StatusInternalServerError)
		return
	}

	s.Wsm.SendToUser(uint64(req.DriverID), []byte(fmt.Sprintf(
		`{"event":"ride_accepted","ride_id":%d,"driver_id":%d}`,
		req.RideID, req.DriverID,
	)))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "ride accepted"})
}

func (s *Server) updateDriverLocation(w http.ResponseWriter, r *http.Request) {
	driverIDStr := chi.URLParam(r, "id")
	var driverID uint64
	fmt.Sscanf(driverIDStr, "%d", &driverID)

	var req struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := s.Redis.SetDriverLocation(r.Context(), driverID, req.Lat, req.Long)
	if err != nil {
		http.Error(w, "failed to update location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Driver location updated successfully",
	})
}
