package directions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"googlemaps.github.io/maps"
	//"github.com/davecgh/go-spew/spew"
)

var (
	mapsClient *maps.Client
	mapsAPIKey string
)

var ErrMalformed = errors.New("Input malformed.\n\n"+
	"For directions, please submit in the form '[mode] from [origin] to [destination]'. "+
	"For example: 'drive from berkeley to sfo'. " +
	"Mode can be either 'walk', 'drive', 'bike', or 'transit'.  \n\n"+
	"You can also get phone number and address info for a business with 'info for [place]'. " +
	"For example: 'info for UC Berkeley'.\n\n" +
	"You can get up to 10 nearby places with 'find [type of place] near [specific place]'. " +
	"For example: 'find grocery near uc berkeley'\n\n"+
	"Lastly, you can get weather info with 'weather for [zip code]'. For example: 'weather for 95555'.")


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

func Get(m string) (string, error) {
	log.Println("getting directions for: ", m)
	m = strings.ToLower(m)

	givenMode := strings.TrimSpace(strings.Split(m, " from ")[0])
	fmt.Println("split", strings.Split(m, " from "))
	fmt.Println("givenMode", "&&"+givenMode+"&&")

	apiMode := maps.TravelModeDriving

	switch givenMode {
	case "walk", "walk!":
		fmt.Println("get", 1)
		apiMode = maps.TravelModeWalking
	case "drive", "drive!":
		fmt.Println("get", 2)
		apiMode = maps.TravelModeDriving
	case "bike", "bike!":
		fmt.Println("get", 3)
		apiMode = maps.TravelModeBicycling
	case "transit", "transit!":
		fmt.Println("get", 4)
		apiMode = maps.TravelModeTransit
	default:
		fmt.Println("get", 5)
		return "", ErrMalformed
	}

	var directions string
	var from string
	var to string

	directionPart := strings.Split(m, " from ")[1]
	fromAndTo := strings.Split(directionPart, " to ")

	if len(fromAndTo) > 1 {
		from = fromAndTo[0]
		to = fromAndTo[1]
	} else {
		return directions, ErrMalformed
	}

	r := &maps.DirectionsRequest{
		Origin:      from,
		Destination: to,
		Mode:        apiMode,
	}

	resp, _, err := mapsClient.Directions(context.Background(), r)

	if err != nil {
		errorCode := randSeq(5)
		log.Println("problem getting directions: ", errorCode, err)
		if strings.Contains(err.Error(), "NOT_FOUND") {
			err = errors.New("sorry, google maps couldn't find either the start or end of your directions. this may happen if something is misspelled or the address is wrong.")
		} else {
			err = errors.New("oops! something went wrong, if you'd like to help us figure it out, please share the code: " + errorCode + " with ouidevelop@gmail.com")
		}
		return directions, err
	}

	for _, route := range resp {
		legs := route.Legs
		fmt.Println("num legs: ", len(legs))
		for _, leg := range legs {
			directions += "estimated time: " + leg.Duration.String() + "\n\n"
			steps := leg.Steps
			fmt.Println("num steps: ", len(steps))
			if len(steps) > 100 {
				return "Directions for your request have over 100 steps. " +
					"That usually happens when the provided start or destination are not specific enough and Google guesses two locations that are very far apart. " +
					"\n\nTry adding a city name to the starting point and destination. ", nil
			}
			for _, step := range steps {
				directions += "** " + removeBold(step.HTMLInstructions) + " (" + step.Distance.HumanReadable + ")"
				if step.TransitDetails != nil {
					directions += "\ndeparture time: " + step.TransitDetails.DepartureTime.String() + "\n"
					directions += "line: " + step.TransitDetails.Line.Name
				}

				directions += "\n\n"

				if len(step.Steps) > 0 {
					for _, subStep := range step.Steps {
						if step.TransitDetails != nil {
							directions += "bob2" + step.TransitDetails.DepartureTime.String() + step.TransitDetails.Line.Name + "\n\n"
						}
						directions += removeBold(subStep.HTMLInstructions) + " (" + subStep.Distance.HumanReadable + ")" + "\n\n"
					}
				}
			}
		}
	}

	log.Println("directions, error : ", directions, err)
	return directions, err
}

func removeBold(str string) string {
	inBold := false
	newStr := ""
	for i, r := range str {
		if inBold {
			newStr += strings.ToUpper(string(r))
		} else {
			newStr += string(r)
		}
		if string(r) == ">" && string(str[i-1]) == "b" && string(str[i-2]) == "<" {
			inBold = true
		}
		if string(r) == ">" && string(str[i-1]) == "b" && string(str[i-2]) == "/" && string(str[i-3]) == "<" {
			inBold = false
		}
	}

	re := regexp.MustCompile(`<b>|<\/B>`)
	newStr = re.ReplaceAllString(newStr, "")
	re = regexp.MustCompile(`<div style="font-size:0.9em">`)
	newStr = re.ReplaceAllString(newStr, " (")
	re = regexp.MustCompile(`<\/div>`)
	newStr = re.ReplaceAllString(newStr, ") ")

	newStr = strings.ReplaceAll(newStr, "/<wbr/>", " ")
	newStr = strings.ReplaceAll(newStr, "/<WBR/>", " ")
	newStr = strings.ReplaceAll(newStr, "&nbsp;", " ")
	return newStr
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
