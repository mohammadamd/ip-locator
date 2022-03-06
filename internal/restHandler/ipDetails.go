package restHandler

import (
	"encoding/csv"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"simple-fh/internal/ipDetails"
	"simple-fh/internal/models"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type getIpDetailsRequest struct {
	IP string `json:"ip" validate:"required,ipv4"`
}

type getIpDetailsResponse struct {
	Ip           string  `json:"ip"`
	CountryCode  string  `json:"country_code"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	MysteryValue int     `json:"mystery_value"`
}

// GetIpDetails returns the details of an IP address
func GetIpDetails() func(c echo.Context) error {
	return func(c echo.Context) error {
		req := new(getIpDetailsRequest)
		if err := c.Bind(req); err != nil {
			logrus.Error("could not bind request: ", err)
			return echo.ErrInternalServerError
		}

		if err := validator.New().Struct(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		details, err := ipDetails.GetIpDetails(req.IP)
		if err == ipDetails.IpNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "ip not found")
		}

		if err != nil {
			logrus.Error("could not get ip details: ", err)
			return echo.ErrInternalServerError
		}

		res := getIpDetailsResponse{
			Ip:           details.Ip,
			CountryCode:  details.CountryCode,
			Country:      details.Country,
			City:         details.City,
			Latitude:     details.Latitude,
			Longitude:    details.Longitude,
			MysteryValue: details.MysteryValue,
		}

		return c.JSON(http.StatusOK, res)
	}
}

// InsertIpDetailsByCSV receive a csv file with ip details and upsert them into the database
func InsertIpDetailsByCSV() func(c echo.Context) error {
	return func(c echo.Context) error {
		f, err := c.FormFile("ip_details")
		if err == http.ErrMissingFile {
			return echo.NewHTTPError(http.StatusBadRequest, "missing file")
		}

		if err != nil {
			logrus.Error("could not get file: ", err)
			return echo.ErrInternalServerError
		}

		of, err := f.Open()
		if err != nil {
			logrus.Error("could not open file: ", err)
			return echo.ErrInternalServerError
		}

		cr := csv.NewReader(of)
		header, err := cr.Read()
		if err != nil {
			logrus.Error("could not read header: ", err)
			return echo.ErrInternalServerError
		}

		if len(header) != 7 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid file header")
		}

		mp := createCsvHeaderMap(header)
		var successCount uint64
		var failedCount uint64
		workerChan := make(chan byte, 50)
		var wg sync.WaitGroup

		for {
			record, err := cr.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				atomic.AddUint64(&failedCount, 1)
				logrus.Error("could not read record: ", err)
				continue
			}

			workerChan <- 1
			wg.Add(1)
			go checkAndInsertIpDetails(mp, record, &failedCount, &successCount, workerChan, &wg)
		}

		wg.Wait()
		return c.JSON(http.StatusOK, map[string]interface{}{"failed": failedCount, "success": successCount})
	}
}

// createCsvHeaderMap creates a dictionary from csv header to map the header to the index in the csv records
func createCsvHeaderMap(header []string) map[string]int {
	m := make(map[string]int)
	for i, h := range header {
		m[strings.ToLower(h)] = i
	}

	return m
}

// checkAndInsertIpDetails checks if the ip details are valid and upsert them into the database
// WorkerChan is a channel which used to limit the number of goroutines
// SuccessCount and FailedCount are used to count the number of successful and failed inserts
func checkAndInsertIpDetails(mp map[string]int, record []string, failedCount *uint64, successCount *uint64, workerChan chan byte, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		_ = <-workerChan
	}()

	Ip := strings.TrimSpace(record[mp["ip"]])
	CountryCode := strings.TrimSpace(record[mp["country_code"]])
	Country := strings.TrimSpace(record[mp["country"]])
	City := strings.TrimSpace(record[mp["city"]])
	Latitude, err := strconv.ParseFloat(strings.TrimSpace(record[mp["latitude"]]), 64)
	if err != nil {
		atomic.AddUint64(failedCount, 1)
		return
	}

	Longitude, err := strconv.ParseFloat(strings.TrimSpace(record[mp["longitude"]]), 64)
	if err != nil {
		atomic.AddUint64(failedCount, 1)
		return
	}

	MysteryValue, err := strconv.ParseInt(strings.TrimSpace(record[mp["mystery_value"]]), 10, 64)
	if err != nil {
		atomic.AddUint64(failedCount, 1)
		return
	}

	ipd := models.IpDetails{
		Ip:           Ip,
		CountryCode:  CountryCode,
		Country:      Country,
		City:         City,
		Latitude:     Latitude,
		Longitude:    Longitude,
		MysteryValue: int(MysteryValue),
	}

	err = validator.New().Struct(ipd)
	if err != nil {
		atomic.AddUint64(failedCount, 1)
		return
	}

	err = ipDetails.ImportIpDetails(ipd)
	if err != nil {
		atomic.AddUint64(failedCount, 1)
		logrus.Error("could not insert into database: ", err)
		return
	}

	atomic.AddUint64(successCount, 1)
}
