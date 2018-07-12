package main

import (
	"log"

	"fmt"
	"net/http"
	"os"

	"encoding/json"

	"math/rand"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/schema"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"googlemaps.github.io/maps"
	"database/sql"
)

var (
	twilioID        string
	twilioAuthToken string
	mapsAPIKey      string
	port            = os.Getenv("PORT")
	mapsClient      *maps.Client
	DB   *sql.DB
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

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		log.Fatal("MYSQL_PASSWORD environment variable not set")
	}
	DB = startDB(mysqlPassword)

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

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Fatal("MAP_API_KEY environment variable not set")
	}

	var err error
	mapsClient, err = maps.NewClient(maps.WithAPIKey(mapsAPIKey))
	if err != nil {
		log.Fatal("problem setting up google maps client: ", err)
	}

	http.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/sms", receiveTextsHandler)
	http.HandleFunc("/directions", directionsHandler) //for demo purposes
	http.HandleFunc("/payment", paymentHandler)       //for demo purposes
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

func paymentHandler(w http.ResponseWriter, r *http.Request) {

	token := r.FormValue("stripeToken")
	email := r.FormValue("stripeEmail")

	params := &stripe.ChargeParams{
		Amount:       stripe.Int64(200),
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Description:  stripe.String("20 Text message Directions"),
		ReceiptEmail: stripe.String(email),
	}
	params.SetSource(token)
	ch, err := charge.New(params)
	spew.Dump(ch, err)
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

	if err != nil {
		log.Println(err)
		directions = err.Error()
	}

	cl := gotwilio.NewTwilioClient(twilioID, twilioAuthToken)
	send(cl, "+13123131234", m.From, directions)
}

// todo: set up Status Callback URL on twilio to make sure message was sent
// todo: add tests
// todo: clean up code
// todo: track free messages and include in the message how many free messages they have left
// todo: provide message in the case that payment didn't work
// todo: provide success message in the case that payment does work.
