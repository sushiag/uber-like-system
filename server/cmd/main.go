package main

// user that requests the rides
type Rider struct {
	ID            uint64   `json:"id"`
	Request       *Request `json:"request,omitempty"`
	RiderLocation Location `json:"rider_location"`
}

// user that accepts rides
type Driver struct {
	ID             uint64       `json:"id"`
	Status         DriverStatus `json:"status"`
	DriverLocation Location     `json:"driver_location"`
}

// this is geo coordinates
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type RiderStatus uint8

const (
	RideRequested RiderStatus = iota
	AssignedRider
	AcceptedRider
	CompletedRider
)

type DriverStatus uint8

const (
	DriverAvailable DriverStatus = iota
	AssignedDriver
	DriverRoute
	DriverCompleted
)

type Request struct {
	ID       uint64      `json:"id"`
	RiderID  uint64      `json:"rider_id"`
	DriverID uint64      `json:"driver_id"`
	Status   RiderStatus `json:"status"`
	PickUp   Location    `json:"pickup"`
	DropOff  Location    `json:"dropoff"`
}

type TrackDriver struct {
	DriverID        uint64     `json:"driver_id"`
	Path            []Location `json:"path"`
	CurrentLocation Location   `json:"current_location"`
}
