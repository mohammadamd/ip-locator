package cmd

import (
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"simple-fh/config"
	"simple-fh/internal/restHandler"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "This command will summon the ip-locator!",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	e := echo.New()
	e.HideBanner = true

	config.Initialize(configPath)
	p := prometheus.NewPrometheus("ip_locator", nil)
	p.Use(e)
	jaegertracing.New(e, nil)

	e.POST("/api/v1/ip/details", restHandler.GetIpDetails())
	e.POST("/api/v1/ip/upload-csv", restHandler.InsertIpDetailsByCSV())
	e.Logger.Fatal(e.Start(":1323"))
}

func init() {
	rootCMD.AddCommand(serveCmd)
}
