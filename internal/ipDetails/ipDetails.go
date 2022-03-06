package ipDetails

import (
	"database/sql"
	"fmt"
	"simple-fh/internal/models"
)

var (
	IpAlreadyExists = fmt.Errorf("IP already exists")
	IpNotFound      = fmt.Errorf("IP not found")
)

// GetIpDetailsForIP will return the details of an IP address
// If the IP address is not found, IpNotFound will be returned
func GetIpDetailsForIP(ip string) (models.IpDetails, error) {
	details, err := getIpDetailsByIpFromDatabase(ip)
	if err == sql.ErrNoRows {
		return details, IpNotFound
	}

	return details, err
}

// InsertNewIPDetails will create a new IP address
// If the IP address exists will update its details
// If the IP address exists and details was same will return IpAlreadyExists
func InsertNewIPDetails(details models.IpDetails) error {
	err := upsertIpDetailsIntoDatabase(details)
	if err == noRowsAffected {
		return IpAlreadyExists
	}

	return err
}
