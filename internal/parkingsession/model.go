package parkingsession

import (
	"time"

	"github.com/shopspring/decimal"
)

type ParkingSession struct {
	ID                 int64           `json:"id"`
	CardUID            string          `json:"cardUid"`
	PlateIn            *string         `json:"plateIn"`
	ImgPlateInPath     *string         `json:"imgPlateInPath"`
	ImgPersonInPath    *string         `json:"imgPersonInPath"`
	CheckInTime        time.Time       `json:"checkInTime"`
	PlateOut           *string         `json:"plateOut"`
	ImgPlateOutPath    *string         `json:"imgPlateOutPath"`
	ImgPersonOutPath   *string         `json:"imgPersonOutPath"`
	CheckOutTime       *time.Time      `json:"checkOutTime"`
	Cost               decimal.Decimal `json:"cost"`
	IsWarning          bool            `json:"isWarning"`
	Status             string          `json:"status"`
}
