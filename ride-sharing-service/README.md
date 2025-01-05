### Ride Sharing Service - Low-Level Design (LLD)

This project is a **Ride Sharing Service** implemented in Golang that simulates a ride-hailing platform. It incorporates modular design principles and provides key features like ride requests, driver-passenger interaction, ride status updates, and fare calculation.

---

### Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Design Components](#design-components)
4. [Key Classes and Responsibilities](#key-classes-and-responsibilities)
5. [Run Instructions](#run-instructions)
6. [Example Usage](#example-usage)

---

### Overview

The **Ride Sharing Service** models core functionality for a ride-hailing platform:
- Passengers request rides.
- Drivers accept ride requests and complete rides.
- Fare calculation based on distance and time.
- Notifications for passengers and drivers.

---

### Features

1. **Passenger Management**: Add and manage passenger details.
2. **Driver Management**: Add, track, and manage driver availability.
3. **Ride Lifecycle**: Handle ride requests, acceptance, status updates, and completion.
4. **Fare Calculation**: Calculate fare dynamically based on distance and duration.
5. **Notifications**: Notify passengers and drivers at key stages of the ride.

---

### Design Components

The service is designed with modular components for better maintainability and scalability:

- **Domain Models**:
  - `Passenger`
  - `Driver`
  - `Ride`
  - `Location`

- **Services**:
  - `PassengerService`: Manages passengers.
  - `DriverService`: Manages drivers and their availability.
  - `RideService`: Handles ride creation and status updates.
  - `NotificationService`: Sends notifications to passengers and drivers.
  - `FareCalculator`: Calculates ride fare.

- **Coordinator**:
  - `RideCoordinator`: Orchestrates the entire ride lifecycle and integrates various services.

---

### Key Classes and Responsibilities

#### **Domain Models**
- **Passenger**: Represents a passenger with an ID, name, contact, and location.
- **Driver**: Represents a driver with an ID, name, contact, license plate, location, and availability status.
- **Ride**: Represents a ride with details like source, destination, passenger, driver, fare, and ride status.
- **Location**: Represents geographical coordinates (latitude and longitude).

#### **PassengerService**
- Adds and manages passengers in the system.

#### **DriverService**
- Adds drivers and fetches a list of available drivers.

#### **RideService**
- Creates and updates rides in the system.

#### **NotificationService**
- Sends notifications to passengers and drivers.

#### **FareCalculator**
- Calculates ride fare based on distance and duration.

#### **RideCoordinator**
- Integrates services and coordinates ride requests, acceptance, starting, completion, and cancellation.

---

### Run Instructions

1. Clone this repository to your local machine.
2. Install [Go](https://golang.org/) if not already installed.
3. Run the `Run` function in the `main.go` file:
   ```bash
   go run cmd/main.go
   ```

---

### Example Usage

The `Run` function demonstrates the flow of a typical ride in the system:

1. **Setup**:
   - Create services for passengers, drivers, rides, notifications, and fare calculation.
   - Initialize the `RideCoordinator`.

2. **Passenger and Driver Creation**:
   ```go
   passenger := &Passenger{ID: 1, Name: "John", Contact: "12345", Location: &Location{Latitude: 37.77, Longitude: -122.42}}
   driver := &Driver{ID: 1, Name: "Alice", Contact: "67890", LicensePlate: "XYZ123", Location: &Location{Latitude: 37.77, Longitude: -122.42}, Status: Available}

   passengerService.AddPassenger(passenger)
   driverService.AddDriver(driver)
   ```

3. **Request and Manage Rides**:
   - A passenger requests a ride:
     ```go
     rideCoordinator.RequestRide(passenger, passenger.Location, &Location{Latitude: 37.78, Longitude: -122.43})
     ```
   - A driver accepts the ride:
     ```go
     ride := <-rideCoordinator.requestedRides
     rideCoordinator.AcceptRide(driver, ride)
     ```
   - The driver starts the ride:
     ```go
     rideCoordinator.StartRide(driver, ride)
     ```
   - The ride is completed:
     ```go
     rideCoordinator.CompleteRide(ride)
     ```

4. **Notifications**:
   Passengers and drivers are notified at key stages (e.g., ride accepted, started, completed).

---

### Output Example

```plaintext
Notifying driver Alice about ride request: 169478293490233
Notifying passenger John: Your ride has been accepted by driver: Alice
Notifying passenger John: Your ride has been started by driver: Alice
Notifying passenger John: Your ride is completed. Fare: $8.50
Notifying driver Alice: Ride completed. Fare: $8.50
```

---

### Future Enhancements

1. **Driver Matching**: Use advanced algorithms to match passengers with drivers.
2. **Route Optimization**: Integrate real-time maps for route and fare calculation.
3. **Payment Gateway**: Add a payment service for ride payments.
4. **Error Handling**: Improve error handling for edge cases.
