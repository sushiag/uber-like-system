package api

import (
	"encoding/json"
	"net/http"
)

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
