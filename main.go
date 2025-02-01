package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

type Weather struct {
	Description     string `json:"description"`
	OpenWeatherIcon string `json:"icon"`
}

type Main struct {
	Temperature float64 `json:"temp"`
	Pressure    int     `json:"pressure"`
	Humidity    float64 `json:"humidity"`
}

type OpenWeatherData struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
}

type Config struct {
	InfluxDBHost         string `env:"INFLUXDB_HOST" envDefault:"localhost"`
	InfluxDBPort         string `env:"INFLUXDB_PORT" envDefault:"8086"`
	InfluxDBDatabase     string `env:"INFLUXDB_DATABASE" envDefault:"telemetry"`
	InfluxDBOrganization string `env:"INFLUXDB_ORGANIZATION" envDefault:"iot"`
	InfluxDBToken        string `env:"INFLUXDB_TOKEN" envDefault:""`
	InfluxDBMeasurement  string `env:"INFLUXDB_MEASUREMENT" envDefault:"iot_data"`
	Lat                  string `env:"LAT"`
	Lon                  string `env:"LON"`
	OpenWeatherApiKey    string `env:"OPENWEATHER_API_KEY"`
	OpenWeatherUnits     string `env:"OPENWEATHER_UNITS" envDefault:"metric"`
	OpenWeatherLang      string `env:"OPENWEATHER_LANG" envDefault:"en"`
}

var cfg = Config{}

func getWeatherData() (*OpenWeatherData, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=%s&lang=%s", cfg.Lat, cfg.Lon, cfg.OpenWeatherApiKey, cfg.OpenWeatherUnits, cfg.OpenWeatherLang)

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	data := &OpenWeatherData{}
	jsonErr := json.Unmarshal([]byte(body), &data)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return data, nil
}

func writeToInfluxDb(data *OpenWeatherData) {
	client := influxdb2.NewClient(
		cfg.InfluxDBHost+":"+cfg.InfluxDBPort,
		cfg.InfluxDBToken,
	)
	writeAPI := client.WriteAPIBlocking(cfg.InfluxDBOrganization, cfg.InfluxDBDatabase)

	node := fmt.Sprintf("%s_%d", strings.ToLower(data.Name), data.ID)

	var tempUnit string
	switch cfg.OpenWeatherUnits {
	case "metric":
		tempUnit = "C"
	case "imperial":
		tempUnit = "F"
	default:
		tempUnit = "K"
	}

	tags := map[string]string{"node": node, "tempUnit": tempUnit}
	fields := map[string]interface{}{
		"temperature": data.Main.Temperature,
		"humidity":    data.Main.Humidity,
		"description": data.Weather[0].Description,
		"owIcon":      data.Weather[0].OpenWeatherIcon,
	}
	point := influxdb2.NewPoint(cfg.InfluxDBMeasurement, tags, fields, time.Now())
	err := writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
}

func init() {
	_ = godotenv.Load()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}

func main() {
	data, err := getWeatherData()
	if err != nil {
		log.Fatal(err)
	}
	writeToInfluxDb(data)
}
