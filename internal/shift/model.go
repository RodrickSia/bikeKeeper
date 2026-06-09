package shift

import "time"

type Shift struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	StartTime   string    `json:"startTime"`
	EndTime     string    `json:"endTime"`
	Date        string    `json:"date"`
	Status      string    `json:"status"`
	Notes       *string   `json:"notes"`
	StaffIDs    []string  `json:"staffIds"`
	StaffNames  []string  `json:"staffNames"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Assignment struct {
	ShiftID    string    `json:"shiftId"`
	UserID     string    `json:"userId"`
	AssignedAt time.Time `json:"assignedAt"`
	AssignedBy *string   `json:"assignedBy"`
}
