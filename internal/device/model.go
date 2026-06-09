package device

import "time"

type Device struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	LocationLabel   string    `json:"locationLabel"`
	Status          string    `json:"status"`
	IPAddress       *string   `json:"ipAddress"`
	FirmwareVersion *string   `json:"firmwareVersion"`
	Notes           *string   `json:"notes"`
	InstalledAt     string    `json:"installedAt"`
	LastSeen        time.Time `json:"lastSeen"`
	AlertCount      int       `json:"alertCount"`
}

type Alert struct {
	ID         string     `json:"id"`
	DeviceID   string     `json:"deviceId"`
	DeviceName string     `json:"deviceName"`
	Message    string     `json:"message"`
	Severity   string     `json:"severity"`
	CreatedAt  time.Time  `json:"createdAt"`
	ResolvedAt *time.Time `json:"resolvedAt"`
}
