package user

import "time"

const (
	RoleStudent = "student"
	RoleStaff   = "staff"
	RoleFaculty = "faculty"
	RoleAdmin   = "admin"

	StatusPending   = "pending_approval"
	StatusActive    = "active"
	StatusRejected  = "rejected"
	StatusSuspended = "suspended"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	MemberID     *string   `json:"memberId,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}
