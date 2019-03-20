package main

import (
	"customFfmpeg/pkg/media-check-tool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/micro/go-log"
	"net/http"
)


type DefaultRequest struct {
	URL               string  `json:"url"`
	CustomParameters  []string  `json:"custom_parameters"`
}

func main(){
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", checkHealth)
	e.GET("/healthz", checkHealth)
	e.POST("/check/ffprobe", Ffprobe)
	e.POST("/check/ffprobe/custom", FfprobeCustom)
	e.POST("/check/media-info", MediaInfoRun)
	e.Logger.Fatal(e.Start(":1333"))
}

// Health Check for Kubernetes Reasons
func checkHealth(c echo.Context) error {
	return c.String(http.StatusOK, "200")
}

// FFProbe: returns meta data of media
func Ffprobe(c echo.Context) error {
	u := new(DefaultRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	data, error := media_check_tool.FFProbe(u.URL)

	log.Logf("%v", u.URL)

	if error != nil {
		return c.JSON(http.StatusInternalServerError, error)
	}
	return c.JSON(http.StatusOK, data)
}

// returns metadata of media with custom parameters for tool
func FfprobeCustom(c echo.Context) error {
	u := new(DefaultRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	data, error := media_check_tool.FFProbeCustom(u.URL, u.CustomParameters)

	log.Logf("%v, %v", u.URL, u.CustomParameters)

	if error != nil {
		return c.JSON(http.StatusInternalServerError, error)
	}
	return c.JSON(http.StatusOK, data)
}

// returns metadata pulled by mediainfo tool
func MediaInfoRun(c echo.Context) error {
	u := new(DefaultRequest)

	if err := c.Bind(u); err != nil {
		return err
	}

	data, error :=  media_check_tool.MediaInfo(u.URL)

	if error != nil {
		return c.JSON(http.StatusInternalServerError, error)
	}
	return c.JSON(http.StatusOK, data)
}