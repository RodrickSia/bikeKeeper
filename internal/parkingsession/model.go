package parkingsession

import (
	"time"
)

type ParkingSession struct {
	ID                 int64      `json:"id"`
	CardUID            string     `json:"cardUid"`
	PlateIn            *string    `json:"plateIn"`
	ImgPlateInPath     *string    `json:"imgPlateInPath"`
	ImgPersonInPath    *string    `json:"imgPersonInPath"`
	CheckInTime        time.Time  `json:"checkInTime"`
	PlateOut           *string    `json:"plateOut"`
	ImgPlateOutPath    *string    `json:"imgPlateOutPath"`
	ImgPersonOutPath   *string    `json:"imgPersonOutPath"`
	CheckOutTime       *time.Time `json:"checkOutTime"`
	Status             string     `json:"status"`
}
