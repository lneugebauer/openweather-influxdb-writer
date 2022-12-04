# openweather-influxdb-writer

A simple OpenWeather InfluxDB writer developed in Go.

## Installation and usage

```shell
$ mkdir /opt/openweather-influxdb-writer
$ cd /opt/openweather-influxdb-writer
$ wget https://github.com/lneugebauer/openweather-influxdb-writer/releases/download/{tag}/openweather-influxdb-writer_{tag}_{os}_{arch}.tar.gz && https://raw.githubusercontent.com/lneugebauer/openweather-influxdb-writer/{tag}/.env.example
$ tar -xzf openweather-influxdb-writer_{tag}_{os}_{arch}.tar.gz
$ mv .env.example .env
$ vim .env
$ chmod +x openweather-influxdb-writer
$ crontab -e
# m h  dom mon dow   command
*/15 * * * * cd /opt/openweather-influxdb-writer && ./openweather-influxdb-writer
```

## MQTT InfluxDB bridge

I've also developed a program to [write weather data from MQTT devices to InfluxDB](https://github.com/lneugebauer/mqtt-influxdb-bridge). Take a look at it in case you are interested.