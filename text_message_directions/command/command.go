package command

import (
	"errors"
	"fmt"
	"github.com/ouidevelop/dontfearthesweeper/text_message_directions/command/directions"
	"github.com/ouidevelop/dontfearthesweeper/text_message_directions/command/places"
	"github.com/ouidevelop/dontfearthesweeper/text_message_directions/command/weather"
	"strings"
)

func Do (input string) (string, error) {

	input = strings.ToLower(input)

	var response string
	var err error

	if strings.HasPrefix(input, "info for ") {
		response, err = places.Get(input)
	} else if strings.HasPrefix(input, "[") {
		input2 := strings.ReplaceAll(input, "[", "")
		input3 := strings.ReplaceAll(input2, "]", "")
		return "", errors.New("input problem, try removing the brackets like this: '" + input3 + "'")
	} else if strings.HasPrefix(input, "find ") && strings.Contains(input, " near "){
		response, err = places.Get(input)
	} else if strings.HasPrefix(input, "weather for ") {
		parts := strings.Fields(input)
		if (len(parts) != 3) && (len(parts) != 4) {
			return "", errors.New("command should be in form 'weather for [zip]' (works in USA and Canada)")
		}
		response, err = weather.GetWeather(strings.TrimPrefix(input, "weather for "))
	} else if strings.HasPrefix(input, "weather ") {
		response, err = weather.GetWeather(strings.TrimPrefix(input, "weather "))
	} else {
		givenMode := strings.Split(input, " from ")[0]
		fmt.Println("givenMode", givenMode)

		response, err = directions.Get(input)
	}

	response = randomlyAddDonationMessage(true, response, err, 7)
	return response, err
}


var count = 9
func randomlyAddDonationMessage(should bool, directions string, err error, every int) string {
	if err != nil {
		return directions
	}
	
	count++
	donationMessage := "\n If you'd like to support this service and are able to, please consider donating. If you search online for this phone number + 'github', I have a page with a donation button. thanks -mike"
	
	if should && count % every == 0 {
		directions = directions + donationMessage
	}
	return  directions
}