package models

type IpDetails struct {
	Id           int
	Ip           string  `validate:"required,ipv4" json:"ip"`
	CountryCode  string  `validate:"required,iso3166_1_alpha2" json:"country_code"`
	Country      string  `validate:"required,alpha" json:"country"`
	City         string  `validate:"required,alpha" json:"city"`
	Latitude     float64 `validate:"required,latitude" json:"latitude"`
	Longitude    float64 `validate:"required,longitude" json:"longitude"`
	MysteryValue int     `validate:"required,numeric" json:"mystery_value"`
}
