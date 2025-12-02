package geo

import (
	"biturl/internal/helper"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
)

type GeoRedisCache struct {
	RDB   *redis.Client
	GEODB *geoip2.Reader
}

func InitGeoDB(path string, redisDB *redis.Client) *GeoRedisCache {
	var err error
	db, err := geoip2.Open(path)

	if err != nil {
		log.Fatalf("failed to load geo db: %v", err)
	}
	return &GeoRedisCache{RDB: redisDB, GEODB: db}
}

func (g *GeoRedisCache) LookupIP(ipString string, ctx context.Context) (country, city string, err error) {
	// fmt.Println(g.GEODB.City(net.ParseIP("8.8.8.8")))
	if g == nil || g.GEODB == nil {
		return "", "", fmt.Errorf("GEO db not initialized\n")
	}

	ip := net.ParseIP(ipString)
	record, err := g.GEODB.City(ip)
	if err != nil {
		return "", "", err
	}

	country = record.Country.Names["en"]
	city = record.City.Names["en"]

	if city != "" {
		return country, city, nil
	}

	cached := g.RDB.Get(ctx, "geo:"+ipString)
	if cached.Val() != "" {
		parts := strings.Split(cached.Val(), "|")
		return parts[0], parts[1], nil
	}

	timeDay := helper.GenerateDate(1)
	ttl := time.Until(*timeDay)
	if ttl < 0 {
		ttl = 0
	}

	country, city, err = GetGeoInfoFromIPAPI(ipString)
	err = g.RDB.Set(ctx, "geo:"+ipString, fmt.Sprintf("%v|%v", country, city), ttl).Err()
	if err != nil {
		fmt.Println(err)
	}

	return country, city, nil
}
