package ridesharingservice

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Location struct
type Location struct {
	Latitude  float64
	Longitude float64
}

// Enums
type DriverStatus int
type RideStatus int

const (
	Available DriverStatus = iota
	Busy
)

const (
	Requested RideStatus = iota
	Accepted
	InProgress
	Completed
	Cancelled
)

// Domain Models
type Passenger struct {
	ID       int
	Name     string
	Contact  string
	Location *Location
}

type Driver struct {
	ID           int
	Name         string
	Contact      string
	LicensePlate string
	Location     *Location
	Status       DriverStatus
}

type Ride struct {
	ID          int
	Passenger   *Passenger
	Driver      *Driver
	Source      *Location
	Destination *Location
	Status      RideStatus
	Fare        float64
}

// PassengerService
type PassengerService struct {
	passengers map[int]*Passenger
	mu         sync.Mutex
}

func NewPassengerService() *PassengerService {
	return &PassengerService{passengers: make(map[int]*Passenger)}
}

func (ps *PassengerService) AddPassenger(passenger *Passenger) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.passengers[passenger.ID] = passenger
}

// DriverService
type DriverService struct {
	drivers map[int]*Driver
	mu      sync.Mutex
}

func NewDriverService() *DriverService {
	return &DriverService{drivers: make(map[int]*Driver)}
}

func (ds *DriverService) AddDriver(driver *Driver) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.drivers[driver.ID] = driver
}

func (ds *DriverService) GetAvailableDrivers() []*Driver {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	var availableDrivers []*Driver
	for _, driver := range ds.drivers {
		if driver.Status == Available {
			availableDrivers = append(availableDrivers, driver)
		}
	}
	return availableDrivers
}

// RideService
type RideService struct {
	rides map[int]*Ride
	mu    sync.Mutex
}

func NewRideService() *RideService {
	return &RideService{rides: make(map[int]*Ride)}
}

func (rs *RideService) AddRide(ride *Ride) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.rides[ride.ID] = ride
}

func (rs *RideService) UpdateRideStatus(rideID int, status RideStatus) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if ride, exists := rs.rides[rideID]; exists {
		ride.Status = status
	}
}

// NotificationService
type NotificationService struct{}

func (ns *NotificationService) NotifyPassenger(passenger *Passenger, message string) {
	fmt.Printf("Notifying passenger %s: %s\n", passenger.Name, message)
}

func (ns *NotificationService) NotifyDriver(driver *Driver, message string) {
	fmt.Printf("Notifying driver %s: %s\n", driver.Name, message)
}

// FareCalculator
type FareCalculator struct{}

func (fc *FareCalculator) CalculateFare(distance, duration float64) float64 {
	baseFare := 2.0
	perKmFare := 1.5
	perMinuteFare := 0.25
	return baseFare + (distance * perKmFare) + (duration * perMinuteFare)
}

// RideCoordinator
type RideCoordinator struct {
	passengerService    *PassengerService
	driverService       *DriverService
	rideService         *RideService
	notificationService *NotificationService
	fareCalculator      *FareCalculator
	requestedRides      chan *Ride
}

func NewRideCoordinator(ps *PassengerService, ds *DriverService, rs *RideService, ns *NotificationService, fc *FareCalculator) *RideCoordinator {
	return &RideCoordinator{
		passengerService:    ps,
		driverService:       ds,
		rideService:         rs,
		notificationService: ns,
		fareCalculator:      fc,
		requestedRides:      make(chan *Ride, 10),
	}
}

func (rc *RideCoordinator) RequestRide(passenger *Passenger, source, destination *Location) {
	ride := &Ride{
		ID:          int(time.Now().UnixNano()),
		Passenger:   passenger,
		Source:      source,
		Destination: destination,
		Status:      Requested,
	}
	rc.rideService.AddRide(ride)
	rc.requestedRides <- ride
	// Notify available drivers
	for _, driver := range rc.driverService.GetAvailableDrivers() {
		fmt.Printf("Notifying driver %s about ride request: %d\n", driver.Name, ride.ID)
	}
}

func (rc *RideCoordinator) AcceptRide(driver *Driver, ride *Ride) {
	driver.Status = Busy
	ride.Driver = driver
	ride.Status = Accepted
	rc.notificationService.NotifyPassenger(ride.Passenger, fmt.Sprintf("Your ride has been accepted by driver: %s", driver.Name))
}

func (rc *RideCoordinator) StartRide(driver *Driver, ride *Ride) {
	driver.Status = Busy
	ride.Driver = driver
	ride.Status = InProgress
	rc.notificationService.NotifyPassenger(ride.Passenger, fmt.Sprintf("Your ride has been started by driver: %s", driver.Name))
}

func (rc *RideCoordinator) CancelRide(driver *Driver, ride *Ride) {
	driver.Status = Busy
	ride.Driver = driver
	ride.Status = Cancelled
	rc.notificationService.NotifyPassenger(ride.Passenger, fmt.Sprintf("Your ride has been cancelled. Driver: %s", driver.Name))
}

func (rc *RideCoordinator) CompleteRide(ride *Ride) {
	distance := rand.Float64() * 10 // Mocked distance
	duration := distance / 30 * 60  // Mocked duration
	fare := rc.fareCalculator.CalculateFare(distance, duration)
	ride.Status = Completed
	ride.Fare = fare
	ride.Driver.Status = Available

	rc.notificationService.NotifyPassenger(ride.Passenger, fmt.Sprintf("Your ride is completed. Fare: $%.2f", fare))
	rc.notificationService.NotifyDriver(ride.Driver, fmt.Sprintf("Ride completed. Fare: $%.2f", fare))
}

// Main function
func Run() {
	passengerService := NewPassengerService()
	driverService := NewDriverService()
	rideService := NewRideService()
	notificationService := &NotificationService{}
	fareCalculator := &FareCalculator{}

	rideCoordinator := NewRideCoordinator(passengerService, driverService, rideService, notificationService, fareCalculator)

	// Add passengers and drivers
	passenger := &Passenger{ID: 1, Name: "John", Contact: "12345", Location: &Location{Latitude: 37.77, Longitude: -122.42}}
	driver := &Driver{ID: 1, Name: "Alice", Contact: "67890", LicensePlate: "XYZ123", Location: &Location{Latitude: 37.77, Longitude: -122.42}, Status: Available}

	passengerService.AddPassenger(passenger)
	driverService.AddDriver(driver)

	// Request and complete a ride
	rideCoordinator.RequestRide(passenger, passenger.Location, &Location{Latitude: 37.78, Longitude: -122.43})

	ride := <-rideCoordinator.requestedRides

	rideCoordinator.AcceptRide(driver, ride)

	rideCoordinator.StartRide(driver, ride)

	rideCoordinator.CompleteRide(ride)
}
