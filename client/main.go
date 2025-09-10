package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Configurable parameters
const (
	BaseURL         = "http://localhost:8080"
	NumRiders       = 50
	NumDrivers      = 20
	SimulationDelay = 100 * time.Millisecond
)

// RideRequest represents the request body for requesting a ride
type RideRequest struct {
	RiderID     int     `json:"rider_id"`
	PickupLat   float64 `json:"pickup_lat"`
	PickupLong  float64 `json:"pickup_long"`
	DropoffLat  float64 `json:"dropoff_lat"`
	DropoffLong float64 `json:"dropoff_long"`
}

// DriverLocation represents driver location updates
type DriverLocation struct {
	DriverID  int     `json:"driver_id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("Starting load test simulation...")

	// Simulate drivers moving
	for i := 1; i <= NumDrivers; i++ {
		wg.Add(1)
		go func(driverID int) {
			defer wg.Done()
			for j := 0; j < 10; j++ { // 10 location updates
				loc := DriverLocation{
					DriverID:  driverID,
					Latitude:  37.7749 + float64(j)*0.001,
					Longitude: -122.4194 + float64(j)*0.001,
				}
				postJSON(fmt.Sprintf("%s/drivers/%d/location", BaseURL, driverID), loc)
				time.Sleep(SimulationDelay)
			}
		}(i)
	}

	// riders requesting rides
	for i := 1; i <= NumRiders; i++ {
		wg.Add(1)
		go func(riderID int) {
			defer wg.Done()
			rideReq := RideRequest{
				RiderID:     riderID,
				PickupLat:   37.7749,
				PickupLong:  -122.4194,
				DropoffLat:  37.7849,
				DropoffLong: -122.4094,
			}
			postJSON(fmt.Sprintf("%s/rides/request", BaseURL), rideReq)
		}(i)
	}

	wg.Wait()
	fmt.Println("Load test simulation finished!")
}

func postJSON(url string, data interface{}) {
	body, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[%s] Status: %d\n", url, resp.StatusCode)
}
