package places

import (
	"context"
	"fmt"
	"log"
	"os"

	"googlemaps.github.io/maps"
)

var (
	mapsClient *maps.Client
	mapsAPIKey string
)

func init() {

	mapsAPIKey = os.Getenv("MAP_API_KEY")
	if mapsAPIKey == "" {
		log.Fatal("MAP_API_KEY environment variable not set")
	}

	var err error
	mapsClient, err = maps.NewClient(maps.WithAPIKey(mapsAPIKey))
	if err != nil {
		log.Fatal("problem setting up google maps client: ", err)
	}
}

func Get(input string) (string, error) {
	resp, err := mapsClient.TextSearch(context.Background(), &maps.TextSearchRequest{
		Query: input,
	})

	info := ""

	fmt.Println("from places: ", resp.Results)

	for index, result := range resp.Results {
		placeDetails, err := mapsClient.PlaceDetails(context.Background(), &maps.PlaceDetailsRequest{
			PlaceID: result.PlaceID,
		})
		if err != nil {
			return "", err
		}

		info += result.Name + "\n" + placeDetails.FormattedAddress + "\n" + placeDetails.FormattedPhoneNumber + "\n\n"

		if index > 10 {
			break
		}
	}

	if info == "" {
		info = "couldn't find matching results"
	}

	return info, err
}
