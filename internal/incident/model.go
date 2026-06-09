package incident

import "time"

type Incident struct {
	ID           string     `json:"id"`
	ReportedBy   string     `json:"reportedBy"`
	ReporterName string     `json:"reporterName"`
	VehiclePlate *string    `json:"vehiclePlate"`
	Type         string     `json:"type"`
	Description  string     `json:"description"`
	Location     *string    `json:"location"`
	Status       string     `json:"status"`
	ResolvedAt   *time.Time `json:"resolvedAt"`
	ResolvedNote *string    `json:"resolvedNote"`
	CreatedAt    time.Time  `json:"createdAt"`
}
