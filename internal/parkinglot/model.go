package parkinglot

import "time"

type ParkingLot struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Address          string    `json:"address"`
	Type             string    `json:"type"`
	Status           string    `json:"status"`
	TotalCapacity    int       `json:"totalCapacity"`
	CurrentOccupancy int       `json:"currentOccupancy"`
	OpenTime         string    `json:"openTime"`
	CloseTime        string    `json:"closeTime"`
	ContactPhone     *string   `json:"contactPhone"`
	ManagerName      *string   `json:"managerName"`
	Description      *string   `json:"description"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
