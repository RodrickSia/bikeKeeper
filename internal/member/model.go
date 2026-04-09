package member

import "time"

type Member struct {
	ID        string    `json:"id"`
	StudentID string    `json:"studentId"`
	FullName  string    `json:"fullName"`
	Phone     *string   `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
}
