package main

import (
	"context"
	"errors"
	"googlemaps.github.io/maps"
	"log"
	"regexp"
	"strings"
)

func getDirections(m string) (string, error) {
	log.Println("getting directions for: ", m)
	m = strings.ToLower(m)

	var directions string
	var from string
	var to string

	fromAndTo := strings.Split(m, " to ")

	if len(fromAndTo) > 1 {
		from = fromAndTo[0]
		to = fromAndTo[1]
	} else {
		return directions, errors.New("directions malformed. Please submit in the form [origin] to [destination]. For example: 'berkeley to sfo'")
	}

	r := &maps.DirectionsRequest{
		Origin:      from,
		Destination: to,
	}

	resp, _, err := mapsClient.Directions(context.Background(), r)
	if err != nil {
		errorCode := randSeq(5)
		log.Println("problem getting directions: ", errorCode, err)
		err = errors.New("oops! something went wrong, if you'd like to help us figure it out, please share the code: " + errorCode + " with ouidevelop@gmail.com")
		return directions, err
	}

	for _, route := range resp {
		legs := route.Legs
		for _, leg := range legs {
			directions += "estimated time: " + leg.Duration.String() + "\n\n"
			steps := leg.Steps
			for _, step := range steps {
				directions += removeBold(step.HTMLInstructions) + " (" + step.Distance.HumanReadable + ")" + "\n\n"
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

	return newStr
}
