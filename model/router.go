package model

// Router holds information about a router in the area of sniffer
type Router struct {
	MAC  string `gorm:"primary_key"`
	SSID string
}
