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
	fmt.Println("input2", "$$"+input+"%%")

	if strings.HasPrefix(input, "info for ") {
		return places.Get(input)
	} else if strings.HasPrefix(input, "[") {
		input2 := strings.ReplaceAll(input, "[", "")
		input3 := strings.ReplaceAll(input2, "]", "")
		return "", errors.New("input problem, try removing the brackets like this: '" + input3 + "'")
	} else if strings.HasPrefix(input, "find ") && strings.Contains(input, " near "){
		return places.Get(input)
	} else if strings.HasPrefix(input, "weather for ") {
		parts := strings.Fields(input)
		if (len(parts) != 3) && (len(parts) != 4) {
			return "", errors.New("command should be in form 'weather for [zip]' (works in USA and Canada)")
		}
		return weather.GetWeather(strings.TrimPrefix(input, "weather for "))
	} else {
		givenMode := strings.Split(input, " from ")[0]
		fmt.Println("givenMode", givenMode)

		return directions.Get(input)
	}
}