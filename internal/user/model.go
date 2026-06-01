package user

import "time"

const (
	RoleStudent = "student"
	RoleStaff   = "staff"
	RoleFaculty = "faculty"
	RoleAdmin   = "admin"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	MemberID     *string   `json:"memberId,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}
