package ipDetails

import (
	"errors"
	"simple-fh/database"
	"simple-fh/internal/models"
)

var noRowsAffected = errors.New("no rows affected")

const (
	getIpDetailsByIpQuery = "SELECT ip_address,country_code,country,city,latitude,longitude,mystery_value FROM ip_details WHERE ip_address = $1"
	upsertIpDetailsQuery  = "INSERT INTO ip_details (ip_address,country_code,country,city,latitude,longitude,mystery_value,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,now(),now()) ON CONFLICT (ip_address) " +
		"DO UPDATE SET country_code = EXCLUDED.country_code, country = EXCLUDED.country, city = EXCLUDED.city, latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude, mystery_value = EXCLUDED.mystery_value, updated_at = now()"
)

// getIpDetailsByIpFromDatabase returns the ip details for the given ip address
// If the ip address is not found, it returns sql.ErrNoRows
func getIpDetailsByIpFromDatabase(ip string) (models.IpDetails, error) {
	var ipDetails models.IpDetails
	r := database.GetDatabaseConnection().QueryRow(getIpDetailsByIpQuery, ip)
	if r.Err() != nil {
		return ipDetails, r.Err()
	}

	err := r.Scan(&ipDetails.Ip, &ipDetails.CountryCode, &ipDetails.Country, &ipDetails.City, &ipDetails.Latitude, &ipDetails.Longitude, &ipDetails.MysteryValue)
	if err != nil {
		return ipDetails, err
	}

	return ipDetails, nil
}

// upsertIpDetailsIntoDatabase inserts or updates the ip details for the given ip address
// If the ip address exists it will try to update it, if it does not exist it will insert it
// If there were no rows affected, it returns noRowsAffected
func upsertIpDetailsIntoDatabase(ipDetails models.IpDetails) error {
	res, err := database.GetDatabaseConnection().Exec(upsertIpDetailsQuery, ipDetails.Ip, ipDetails.CountryCode, ipDetails.Country, ipDetails.City, ipDetails.Latitude, ipDetails.Longitude, ipDetails.MysteryValue)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return noRowsAffected
	}

	return nil
}
