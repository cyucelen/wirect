package model

// Router represents database schema of a router
type Router struct {
	MAC        string `gorm:"primary_key"`
	SSID       string
	Sniffer    Sniffer `gorm:"foreignkey:SnifferMAC"`
	SnifferMAC string
}

// RouterExternal holds information about a router in the area of sniffer
type RouterExternal struct {
	MAC  string `json:"MAC"`
	SSID string `json:"SSID"`
}
