package geo

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

func GetGeoInfoFromIPAPI(ip string) (country, city string, err error) {

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

	return geo.Country, geo.City, nil
}
