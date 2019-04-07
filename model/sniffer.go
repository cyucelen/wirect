package model

// Sniffer holds information about a Sniffer
type Sniffer struct {
	MAC      string `gorm:"primary_key" json:"MAC"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
