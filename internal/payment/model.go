package payment

import "time"

const (
	TypeDeposit    = "deposit"
	TypeParkingFee = "parking_fee"

	MethodCardBalance = "card_balance"
	MethodCash        = "cash"
)

type Transaction struct {
	ID            string    `json:"id"`
	CardUID       string    `json:"cardUid"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	PaymentMethod *string   `json:"paymentMethod,omitempty"`
	SessionID     *int64    `json:"sessionId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}
