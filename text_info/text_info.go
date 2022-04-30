package text_info

import (
	"errors"
	"log"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

const newUserMessage = `Due to a large increase in usage, and therefore cost, I'll unfortunately have to start charging for this service.

To learn more about how to pay, text "pay". It's 10 cents per message.

You have 5 free messages (after this one) before you will need to be subscribed to continue. (feel free to retry your request)`

const freeMessagesDone = `You've reached the end of your free messages. Unfortunately I can't support unlimited free use of this service at this time, but may be able to in the future. 

This thing costs me on average 10 cents per message. As I'm not trying to make a profit, that is what I am charging.

To learn more about how to pay, text "pay"`

const paymentMessage = `This thing costs me on average 10 cents per message. As I'm not trying to make a profit, that is what I am charging.You would pay for messages ahead of time. So for example 100 messages would be 10$. 

If you can't access the internet, I can take a credit card payment over the phone at (805) 423 4224.

If you can access the internet, you can pay here: https://github.com/ouidevelop/text-message-directions

There is a paypal button at the bottom of the page. There you can pay with either paypal or a credit card. Once you've paid, let me know and I'll add the messages to your number.`

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
	commandResponse, err := CommandSwitch(q.Query)
	if err != nil {
		commandResponse = err.Error()
	}

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
		resp, err := CommandSwitch(message.Body)
		if err != nil {
			resp = err.Error()
		}
		return resp
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

	if !u.Private {
		log.Printf("incoming message: %+v", message)
	}

	messagesLeft := (freeMessagesLimit - u.FreeMessagesUsed) + u.PaidMessages
	input := strings.ToLower(message.Body)

	if !u.Subscribed && messagesLeft <= 0 && !strings.HasPrefix(input, "pay") {
		return freeMessagesDone
	}

	response, err := CommandSwitch(input)
	if err != nil {
		return err.Error()
	}

	response = updateDB(u, response)

	return response
}

func updateDB(u *User, response string) string {
	if !u.Subscribed {
		if u.FreeMessagesUsed < freeMessagesLimit {
			u.FreeMessagesUsed = u.FreeMessagesUsed + 1
		} else if u.PaidMessages > 0 {
			u.PaidMessages = u.PaidMessages - 1
		}
		upsertUser(u)
		messagesLeft := (freeMessagesLimit - u.FreeMessagesUsed) + u.PaidMessages
		response = response + fmt.Sprintf("\n\n %v msgs left", messagesLeft)
	}
	return response
}

func CommandSwitch(input string) (string, error) {

	var err error
	var response string

	if strings.HasPrefix(input, "pay") {
		return "", errors.New(paymentMessage) // dumb hack to make payment info not count towards message count
	} else if strings.HasPrefix(input, "define ") {
		parts := strings.Split(input, " ")
		if len(parts) != 2 {
			return "", errors.New("command should be of the form 'define car'")
		}
		response, err = getDefinition(parts[1])
		if err != nil {
			log.Println("error getting definition: ", err)
			return "", errors.New("error getting definition")
		}
	} else if strings.HasPrefix(input, "info ") {
		if !strings.HasPrefix(input, "info for ") {
			input = strings.Replace(input, "info ", "info for ", 1)
		}
		response, err = GetPlace(input)
	} else if strings.HasPrefix(input, "[") {
		input2 := strings.ReplaceAll(input, "[", "")
		input3 := strings.ReplaceAll(input2, "]", "")
		return "", errors.New("input problem, try removing the brackets like this: '" + input3 + "'")
	} else if strings.HasPrefix(input, "find ") && strings.Contains(input, " near ") {
		response, err = GetPlace(input)
	} else if strings.HasPrefix(input, "weather for ") {
		parts := strings.Fields(input)
		if (len(parts) != 3) && (len(parts) != 4) {
			return "", errors.New("command should be in form 'weather [zip]' (works in USA and Canada)")
		}
		response, err = GetWeather(strings.TrimPrefix(input, "weather for "))
	} else if strings.HasPrefix(input, "weather ") {
		response, err = GetWeather(strings.TrimPrefix(input, "weather "))
	} else {
		response, err = GetDirections(input)
	}

	response = strings.TrimSuffix(response, "\n\n")
	return response, err
}
