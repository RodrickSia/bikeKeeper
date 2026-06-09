package visitor

import "time"

type VisitorPass struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	VisitorName  string     `json:"visitorName"`
	VisitorPhone *string    `json:"visitorPhone"`
	VehiclePlate string     `json:"vehiclePlate"`
	ValidDate    string     `json:"validDate"`
	Status       string     `json:"status"`
	QRCodeData   string     `json:"qrCodeData"`
	CreatedAt    time.Time  `json:"createdAt"`
	UsedAt       *time.Time `json:"usedAt"`
}
