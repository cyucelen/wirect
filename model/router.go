package model

// Router represents database schema of a router
type Router struct {
	SSID       string  `gorm:"primary_key"`
	Sniffer    Sniffer `gorm:"foreignkey:SnifferMAC"`
	SnifferMAC string  `gorm:"primary_key"`
	LastSeen   int64
}

// RouterExternal holds information about a router in the area of sniffer
type RouterExternal struct {
	SSID     string `json:"SSID"`
	LastSeen int64  `json:"lastSeen"`
}
