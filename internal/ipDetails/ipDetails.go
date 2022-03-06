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

// GetIpDetails will return the details of an IP address
// If the IP address is not found, IpNotFound will be returned
func GetIpDetails(ip string) (models.IpDetails, error) {
	details, err := getIpDetailsByIp(ip)
	if err == sql.ErrNoRows {
		return details, IpNotFound
	}

	return details, err
}

// ImportIpDetails will create a new IP address
// If the IP address exists will update its details
// If the IP address exists and details was same will return IpAlreadyExists
func ImportIpDetails(details models.IpDetails) error {
	err := upsertIpDetails(details)
	if err == noRowsAffected {
		return IpAlreadyExists
	}

	return err
}
