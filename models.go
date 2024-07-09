package main

// Car представляет собой структуру автомобиля
type Car struct {
	ID          string `json:"id"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Mileage     int    `json:"mileage"`
	OwnersCount int    `json:"owners_count"`
}
