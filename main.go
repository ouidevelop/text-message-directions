package main

import (
	"log"
	"regexp"
	"strings"

	"fmt"
	"net/http"
	"os"

	"context"

	"errors"

	"encoding/json"

	"github.com/gorilla/schema"
	"github.com/sfreiberg/gotwilio"
	"googlemaps.github.io/maps"
)

var (
	twilioID        string
	twilioAuthToken string
	mapsAPIKey      string
	port            = os.Getenv("PORT")
)

type query struct {
	Query string `json:"query"`
}

type message struct {
	NumSegments   int
	FromState     string
	ApiVersion    string
	ToState       string
	SmsMessageSid string
	FromCity      string
	Body          string
	To            string
	MessageSid    string
	ToCountry     string
	FromCountry   string
	From          string
	NumMedia      int
	ToZip         int
	ToCity        string
	FromZip       string
	SmsSid        string
	SmsStatus     string
	AccountSid    string
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

func main() {

	http.HandleFunc("/sms", receiveTextsHandler)
	http.HandleFunc("/directions", directionsHandler)
	log.Println("Magic happening on port " + "8080")

	twilioID = os.Getenv("TWILIO_ID")
	if twilioID == "" {
		log.Fatal("TWILIO_ID environment variable not set")
	}

	twilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioAuthToken == "" {
		log.Fatal("TWILIO_AUTH_TOKEN environment variable not set")
	}

	mapsAPIKey = os.Getenv("MAP_API_KEY")
	if mapsAPIKey == "" {
		log.Fatal("MAP_API_KEY environment variable not set")
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func directionsHandler(w http.ResponseWriter, r *http.Request) {
	q := query{}
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		fmt.Fprint(w, "problem decoding json", err)
		return
	}
	d, err := getDirections(q.Query)
	if err != nil {
		fmt.Fprint(w, "problem getting directions: ", err)
		return
	}
	fmt.Fprint(w, d)
}

func receiveTextsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("message recieved!")
	err := r.ParseForm()
	if err != nil {
		fmt.Println("problem parsing form")
	}

	m := message{}

	decoder := schema.NewDecoder()
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(&m, r.PostForm)
	if err != nil {
		fmt.Println("problem decoding form", err)
	}

	fmt.Printf("incoming message: %+v", m)

	directions, err := getDirections(m.Body)

	cl := gotwilio.NewTwilioClient(twilioID, twilioAuthToken)
	send(cl, "+13123131234", m.From, directions)
}

func getDirections(m string) (string, error) {
	fmt.Println("getting directions for: ", m)
	m = strings.ToLower(m)

	var directions string
	var from string
	var to string

	fromAndTo := strings.Split(m, " to ")

	if len(fromAndTo) > 1 {
		from = fromAndTo[0]
		to = fromAndTo[1]
	} else {
		return directions, errors.New("directions malformed")
	}
	fmt.Println("1")

	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBYsgGoxFUK8cDI2cm177AJb9jbthfpNsY"))
	if err != nil {
		fmt.Println("problem setting up google maps client: ", err)
		return directions, err
	}

	r := &maps.DirectionsRequest{
		Origin:      from,
		Destination: to,
	}

	resp, _, err := c.Directions(context.Background(), r)
	if err != nil {
		fmt.Println("problem getting directions: ", err)
		return directions, err
	}
	fmt.Println("3")

	for _, route := range resp {
		legs := route.Legs
		for _, leg := range legs {
			steps := leg.Steps
			for _, step := range steps {
				directions += removeBold(step.HTMLInstructions) + " (" + step.Distance.HumanReadable + ")" + "\n\n"
			}
		}
	}

	fmt.Println("4")
	fmt.Println("directions, error : ", directions, err)
	return directions, err
}

func send(cl *gotwilio.Twilio, from string, to string, body string) {
	fmt.Println("about to send message!")
	res, exeption, err := cl.SendSMS(from, to, body, "", "") // todo: should not ignore those returns
	fmt.Println("res, exeption, err", res, exeption, err)
}
