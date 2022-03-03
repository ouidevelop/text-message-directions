package text_info

import (
	"log"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

const newUserMessage = `Due to a large increase in usage, and therefore cost, I'll unfortunately have to start charging for this service.

If you want to keep using this service, please contact me (michael) at ouidevelop@gmail.com or (805)423-4224, and we can discuss payment options.
If you can't pay but still want to use the service send me a message and we may be able to work something out.
If you're already donating monthly, thanks! Send me a message and I'll make sure you can continue to use the service.

You have 5 free messages (after this one) before you will need to be subscribed to continue. (feel free to retry your request)`

const freeMessagesDone = `You've reached the end of your free messages.

If you want to keep using this service, please contact me (michael) at ouidevelop@gmail.com or (805)423-4224, and we can discuss payment options.
If you can't pay but still want to use the service send me a message and we may be able to work something out.`

const freeMessagesLimit = 5

type query struct {
	Query string `json:"query"`
}

type message struct {
	Body        string
	From        string
	To          string
	FromCountry string
	AccountSid  string
}

func DirectionsHandler(w http.ResponseWriter, r *http.Request) {
	q := query{}
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		_, _ = fmt.Fprint(w, "problem decoding json", err)
		return
	}
	commandResponse := CommandSwitch(q.Query)

	var messages []string

	messages = SplitLongBody(commandResponse)

	if len(messages) > 0 {
		fmt.Println("length of total message, length of first message, different segments", len(commandResponse), len(messages[0]), len(messages))
	} else {
		fmt.Println("no messages", len(commandResponse))
	}

	_, _ = fmt.Fprint(w, messages)
}

func ReceiveTextsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("message recieved!")
	err := r.ParseForm()
	if err != nil {
		log.Println("problem parsing form", err)
	}

	message := message{}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&message, r.PostForm)
	if err != nil {
		log.Println("problem decoding form", err)
	}

	response := formResponse(message)

	fmt.Println("SID:", message.AccountSid)
	Send(message.To, message.From, message.AccountSid, response)
}

func formResponse(message message) string {
	// for people using their own twilio number, don't bother them with any of this other stuff
	if message.To != "+13123131234" && message.To != "+12244847717" {
		fmt.Println("outside number", message.To)
		return CommandSwitch(message.Body)
	}

	u, err := getUser(message.From)
	if u == nil {
		addNewUser(message.From)
		return newUserMessage
	}
	if err != nil {
		log.Println("error checking to see if we need to send initial Message: ", err)
		return "there's been a problem on our end. sorry"
	}
	if u.FreeMessagesUsed >= freeMessagesLimit && !u.Subscribed {
		return freeMessagesDone
	}

	if !u.Private {
		log.Printf("incoming message: %+v", message)
	}

	if !u.Subscribed {
		u.FreeMessagesUsed = u.FreeMessagesUsed + 1
		upsertUser(u)
	}

	return CommandSwitch(message.Body)
}

func CommandSwitch(input string) string {

	input = strings.ToLower(input)

	var err error
	var response string

	if strings.HasPrefix(input, "info for ") {
		response, err = GetPlace(input)
	} else if strings.HasPrefix(input, "[") {
		input2 := strings.ReplaceAll(input, "[", "")
		input3 := strings.ReplaceAll(input2, "]", "")
		return "input problem, try removing the brackets like this: '" + input3 + "'"
	} else if strings.HasPrefix(input, "find ") && strings.Contains(input, " near ") {
		response, err = GetPlace(input)
	} else if strings.HasPrefix(input, "weather for ") {
		parts := strings.Fields(input)
		if (len(parts) != 3) && (len(parts) != 4) {
			return "command should be in form 'weather for [zip]' (works in USA and Canada)"
		}
		response, err = GetWeather(strings.TrimPrefix(input, "weather for "))
	} else if strings.HasPrefix(input, "weather ") {
		response, err = GetWeather(strings.TrimPrefix(input, "weather "))
	} else {
		response, err = GetDirections(input)
	}

	if err != nil {
		return err.Error()
	}
	return response
}
