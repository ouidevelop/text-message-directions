package main

import (
	"log"

	"fmt"
	"net/http"
	"os"

	"encoding/json"

	"math/rand"
	"time"

	"github.com/gorilla/schema"
	"github.com/sfreiberg/gotwilio"
	"googlemaps.github.io/maps"
)

var (
	twilioID        string
	twilioAuthToken string
	mapsAPIKey      string
	port            = os.Getenv("PORT")
	mapsClient      *maps.Client
)

type query struct {
	Query string `json:"query"`
}

type message struct {
	Body        string
	From        string
	FromCountry string
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

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

	var err error
	mapsClient, err = maps.NewClient(maps.WithAPIKey(mapsAPIKey))
	if err != nil {
		log.Fatal("problem setting up google maps client: ", err)
	}

	http.HandleFunc("/sms", receiveTextsHandler)
	http.HandleFunc("/directions", directionsHandler) //for demo purposes
	log.Println("Magic happening on port " + port)
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
	log.Println("message recieved!")
	err := r.ParseForm()
	if err != nil {
		log.Println("problem parsing form")
	}

	m := message{}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&m, r.PostForm)
	if err != nil {
		log.Println("problem decoding form", err)
	}

	log.Printf("incoming message: %+v", m)

	directions, err := getDirections(m.Body)

	cl := gotwilio.NewTwilioClient(twilioID, twilioAuthToken)
	send(cl, "+13123131234", m.From, directions)
}

// todo: set up Status Callback URL on twilio to make sure message was sent
// todo: set up premium sms provider?
//http://www.truesenses.com/website/pages/smspricingpremiumrate/
//http://gateway.txtnation.com/solutions/sms/premium/
// todo: add tests
// todo: clean up code
