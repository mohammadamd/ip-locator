package endToEnd

import (
	"bytes"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"simple-fh/database"
	"simple-fh/internal/restHandler"
	"strings"
	"testing"
)

func TestGetIpDetails(t *testing.T) {
	e := echo.New()

	mock := database.InitializePostgresMock()
	rows := mock.NewRows([]string{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"}).
		AddRow("127.0.0.1", "US", "United States", "New York", "40.7128", "-74.0059", "2342534324")
	mock.ExpectQuery("SELECT ip_address,country_code,country,city,latitude,longitude,mystery_value FROM ip_details").WillReturnRows(rows)
	mock.ExpectQuery("SELECT ip_address,country_code,country,city,latitude,longitude,mystery_value FROM ip_details").WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"ip\":\"127.0.0.1\"}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, restHandler.GetIpDetails()(c))

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"ip\":\"127.0.0.1\"}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)

	assert.Error(t, restHandler.GetIpDetails()(c))

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"sp\":\"127.0.0.1\"}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)

	assert.Error(t, restHandler.GetIpDetails()(c))
}

func TestInsertIpDetailsByCSV(t *testing.T) {
	e := echo.New()

	mock := database.InitializePostgresMock()
	mock.ExpectExec("INSERT INTO ip_details").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO ip_details").WillReturnResult(sqlmock.NewResult(0, 0))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	csvFile, _ := writer.CreateFormFile("ip_details", "ip_details.csv")
	_, _ = csvFile.Write([]byte("ip_address,country_code,country,city,latitude,longitude,mystery_value\n" +
		"200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n" +
		"200.106..15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n" +
		"200.106.141.15,D,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n" +
		"200.106.141.15,SI,0,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n" +
		"200.106.141.15,SI,Nepal,0,-84.87503094689836,7.206435933364332,7823011346\n" +
		"200.106.141.15,SI,Nepal,DuBuquemouth,a,7.206435933364332,7823011346\n" +
		"200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,a,7823011346\n" +
		"200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,a\n" +
		"160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115"))
		_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, restHandler.InsertIpDetailsByCSV()(c))
	res, _ := ioutil.ReadAll(rec.Body)
	assert.Equal(t, "{\"failed\":8,\"success\":1}\n", string(res))
}
