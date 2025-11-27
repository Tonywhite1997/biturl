package helper

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

func GetGeoInfo(ipStr string) (country, city string, err error) {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	ip := net.ParseIP(ipStr)
	record, err := db.City(ip)
	if err != nil {
		return "", "", err
	}

	country = record.Country.Names["en"]
	if len(record.Subdivisions) > 0 {
		city = record.City.Names["en"]
	} else {
		city = ""
	}

	return country, city, nil
}
