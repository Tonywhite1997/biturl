package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GeoIP struct {
	Status  string `json:"status"`
	Country string `json:"country"`
	City    string `json:"city"`
}

var geoCache = make(map[string]GeoIP)

func GetGeoInfo(ip string) (country, city string, err error) {

	// checking memeory cache if ip is cached
	if geo, ok := geoCache[ip]; ok {
		fmt.Println("ip in cache")
		return geo.Country, geo.City, nil
	}

	// gettiing the ip info from ip-api.com api
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var geo GeoIP
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return "", "", err
	}

	if geo.Status != "success" {
		return "", "", fmt.Errorf("failed to get geo info for %s", ip)
	}

	// catching the ip details in memory when found on ip-api.com
	geoCache[ip] = GeoIP{
		Country: geo.Country,
		City:    geo.City,
		Status:  geo.Status,
	}

	fmt.Println(geoCache)

	return geo.Country, geo.City, nil
}
