// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.p

package api

// DB represents an instance of the service's database.
type DB struct {
	Cars     map[string]Car
	Statuses map[string]CarStatus
}

// Car represents a car managed by the service.
type Car struct {
	ID        string `json:"id"`
	Model     string `json:"model"`
	VehicleID string `json:"vehicle_id"`
	OwnerID   string `json:"owner_id"`
}

// CarStatus reprents the status of a car managed by the service.
type CarStatus struct {
	ID       string      `json:"id"`
	Position CarPosition `json:"position"`
	Mileage  int         `json:"mileage"`
	Speed    int         `json:"speed"`
	Fuel     float64     `json:"fuel"`
}

// CarPosition defines the coordinates of the car.
type CarPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func mockDB() DB {

	cars := []Car{
		{
			ID:        "663dc85d-2455-466c-b2e5-76691b0ce14e",
			Model:     "Honda",
			VehicleID: "127482",
			OwnerID:   "742",
		},
		{
			ID:        "879a273c-a8dc-41a6-9c30-2bb92288e93b",
			Model:     "Ford",
			VehicleID: "3784312",
			OwnerID:   "6928",
		},
		{
			ID:        "6c018cfa-e9c2-4169-a61b-dd3bf3bc19a7",
			Model:     "Toyota",
			VehicleID: "19019",
			OwnerID:   "742",
		},
		{
			ID:        "fca3ab25-a151-4c76-b238-9aa6ee92c374",
			Model:     "Honda",
			VehicleID: "22781",
			OwnerID:   "30390",
		},
	}

	statuses := []CarStatus{
		{
			ID: "663dc85d-2455-466c-b2e5-76691b0ce14e",
			Position: CarPosition{
				Latitude:  -39.91045,
				Longitude: -161.70716,
			},
			Mileage: 742,
			Speed:   90,
			Fuel:    6.42,
		},
		{
			ID: "879a273c-a8dc-41a6-9c30-2bb92288e93b",
			Position: CarPosition{
				Latitude:  -8.86414,
				Longitude: -142.59820,
			},
			Mileage: 9347,
			Speed:   45,
			Fuel:    3.1,
		},
		{
			ID: "6c018cfa-e9c2-4169-a61b-dd3bf3bc19a7",
			Position: CarPosition{
				Latitude:  12.77061,
				Longitude: 9.05115,
			},
			Mileage: 17384,
			Speed:   62,
			Fuel:    8.9,
		},
		{
			ID: "fca3ab25-a151-4c76-b238-9aa6ee92c374",
			Position: CarPosition{
				Latitude:  68.86632,
				Longitude: -92.85048,
			},
			Mileage: 97698,
			Speed:   50,
			Fuel:    3.22,
		},
	}

	var db DB

	db.Cars = make(map[string]Car, len(cars))

	for _, car := range cars {
		db.Cars[car.ID] = car
	}

	db.Statuses = make(map[string]CarStatus, len(statuses))

	for _, car := range statuses {
		db.Statuses[car.ID] = car
	}

	return db
}
