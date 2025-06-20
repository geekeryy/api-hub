package google

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"googlemaps.github.io/maps"
)

type GeoLocation struct {
	Lng float64
	Lat float64
}

func Geocode(address, apiKey string) *GeoLocation {
	if address == "" {
		return nil
	}
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		logx.Errorf("google maps new error: %v", err)
		return nil
	}
	geocode, err := c.Geocode(context.Background(), &maps.GeocodingRequest{Address: address})
	if err != nil {
		logx.Errorf("google Geocode error: %v", err)
		return nil
	}
	logx.Errorf("google Geocode, address:%v, Geocode: %+v", geocode)
	if len(geocode) != 1 {
		logx.Errorf("google Geocode Expected length of response is 1, was %+v", geocode)
		return nil
	}
	return &GeoLocation{
		Lat: geocode[0].Geometry.Location.Lat,
		Lng: geocode[0].Geometry.Location.Lng,
	}
}
