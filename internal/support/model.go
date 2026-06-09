package support

import "time"

type Ticket struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	UserName    string     `json:"userName"`
	Category    string     `json:"category"`
	Subject     string     `json:"subject"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Responses   []Response `json:"responses"`
}

type Response struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticketId"`
	SenderID   string    `json:"senderId"`
	SenderName string    `json:"senderName"`
	Message    string    `json:"message"`
	IsAdmin    bool      `json:"isAdmin"`
	CreatedAt  time.Time `json:"createdAt"`
}
