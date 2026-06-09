package cardrequest

import "time"

type Status = string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
	StatusBlocked  Status = "blocked"
)

type CardRequest struct {
	ID             string     `json:"id"`
	MemberID       string     `json:"memberId"`
	VehiclePlate   string     `json:"vehiclePlate"`
	VehicleBrand   string     `json:"vehicleBrand"`
	VehicleModel   string     `json:"vehicleModel"`
	VehicleColor   string     `json:"vehicleColor"`
	IDCardNumber   string     `json:"idCardNumber"`
	Note           *string    `json:"note"`
	Status         Status     `json:"status"`
	CardUID        *string    `json:"cardUid"`
	RejectedReason *string    `json:"rejectedReason"`
	SubmittedAt    time.Time  `json:"submittedAt"`
	ReviewedAt     *time.Time `json:"reviewedAt"`
	ReviewedBy     *string    `json:"reviewedBy"`
}
