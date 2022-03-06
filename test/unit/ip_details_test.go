package unit

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"simple-fh/database"
	"simple-fh/internal/ipDetails"
	"simple-fh/internal/models"
	"testing"
)

func TestGetIpDetailsByIp(t *testing.T) {
	mock := database.InitializePostgresMock()
	rows := mock.NewRows([]string{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"}).
		AddRow("127.0.0.1", "US", "United States", "New York", "40.7128", "-74.0059", "2342534324")
	mock.ExpectQuery("SELECT ip_address,country_code,country,city,latitude,longitude,mystery_value FROM ip_details").WillReturnRows(rows)
	mock.ExpectQuery("SELECT ip_address,country_code,country,city,latitude,longitude,mystery_value FROM ip_details").WillReturnError(sql.ErrNoRows)

	ipd, err := ipDetails.GetIpDetailsForIP("127.0.0.1")
	assert.Nil(t, err)

	assert.Equal(t, ipd, models.IpDetails{
		Ip:           "127.0.0.1",
		CountryCode:  "US",
		Country:      "United States",
		City:         "New York",
		Latitude:     40.7128,
		Longitude:    -74.0059,
		MysteryValue: 2342534324,
	})

	_, err = ipDetails.GetIpDetailsForIP("1.2.3.4")
	assert.Equal(t, err, ipDetails.IpNotFound)
}

func TestImportIpDetails(t *testing.T) {
	mock := database.InitializePostgresMock()
	mock.ExpectExec("INSERT INTO ip_details").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO ip_details").WillReturnResult(sqlmock.NewResult(0, 0))

	err := ipDetails.InsertNewIPDetails(models.IpDetails{
		Ip:           "127.0.0.1",
		CountryCode:  "US",
		Country:      "United States",
		City:         "New York",
		Latitude:     40.7128,
		Longitude:    -74.0059,
		MysteryValue: 2342534324,
	})

	if err != nil {
		t.Errorf("could not import ip details: %s", err)
	}

	err = ipDetails.InsertNewIPDetails(models.IpDetails{
		Ip:           "127.0.0.1",
		CountryCode:  "US",
		Country:      "United States",
		City:         "New York",
		Latitude:     40.7128,
		Longitude:    -74.0059,
		MysteryValue: 2342534324,
	})

	assert.Equal(t, err, ipDetails.IpAlreadyExists)
}
