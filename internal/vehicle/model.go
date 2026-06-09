package vehicle

import "time"

type Vehicle struct {
	ID           string    `json:"id"`
	LicensePlate string    `json:"licensePlate"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Color        string    `json:"color"`
	OwnerID      string    `json:"ownerId"`
	IsActive     bool      `json:"isActive"`
	RegisteredAt time.Time `json:"registeredAt"`
}
