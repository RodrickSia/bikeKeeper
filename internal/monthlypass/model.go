package monthlypass

type MonthlyPass struct {
	ID           string  `json:"id"`
	UserID       string  `json:"userId"`
	VehicleID    string  `json:"vehicleId"`
	VehiclePlate string  `json:"vehiclePlate"`
	VehicleBrand string  `json:"vehicleBrand"`
	Month        string  `json:"month"`
	StartDate    string  `json:"startDate"`
	EndDate      string  `json:"endDate"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	IsAutoRenew  bool    `json:"isAutoRenew"`
	PurchasedAt  string  `json:"purchasedAt"`
}

type CreateParams struct {
	UserID       string
	VehicleID    string
	VehiclePlate string
	VehicleBrand string
	Month        string
	StartDate    string
	EndDate      string
	Price        float64
}
