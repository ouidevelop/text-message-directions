package text_message_directions

import (
	"log"

	"github.com/ouidevelop/dontfearthesweeper/text_message_directions/command"
	"github.com/ouidevelop/dontfearthesweeper/text_message_directions/sms"

	"encoding/json"
	"fmt"
	"net/http"

	"math/rand"
	"time"

	"github.com/gorilla/schema"
)

type query struct {
	Query string `json:"query"`
}

type message struct {
	Body        string
	From        string
	To 			string
	FromCountry string
	AccountSid  string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func DirectionsHandler(w http.ResponseWriter, r *http.Request) {
	q := query{}
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		_, _ = fmt.Fprint(w, "problem decoding json", err)
		return
	}
	d, err := command.Do(q.Query)
	if err != nil {
		_, _ = fmt.Fprint(w, "problem getting directions: ", err)
		return
	}

	var messages []string

	messages = sms.SplitLongBody(d)

	if len(messages) > 0 {
		fmt.Println("bob::::", len(d), len(messages[0]), len(messages))
	} else {
		fmt.Println("no messages", len(d))	
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

	log.Printf("incoming message: %+v", message)

	response, err := command.Do(message.Body)

	if err != nil {
		log.Println(err)
		response = err.Error()
	}

	fmt.Println("SID:", message.AccountSid)
	sms.Send(message.To, message.From, message.AccountSid, response)
}

// todo: set up Status Callback URL on twilio to make sure message was sent
// todo: set up premium sms provider?
//http://www.truesenses.com/website/pages/smspricingpremiumrate/
//http://gateway.txtnation.com/solutions/sms/premium/
// todo: add tests
// todo: clean up code
