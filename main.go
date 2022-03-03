package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/ouidevelop/dontfearthesweeper/streetsweeping"

	"github.com/ouidevelop/dontfearthesweeper/text_info"

	"github.com/dcu/go-authy"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sfreiberg/gotwilio"

	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	twilioID := os.Getenv("TWILIO_ID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioID == "" {
		log.Fatal("TWILIO_ID environment variable not set")
	}
	if twilioAuthToken == "" {
		log.Fatal("TWILIO_AUTH_TOKEN environment variable not set")
	}

	authyAPIKey := os.Getenv("STREETSWEEP_AUTHY_API_KEY")
	if authyAPIKey == "" {
		log.Fatal("STREETSWEEP_AUTHY_API_KEY environment variable not set")
	}

	msgSvc := streetsweeping.TwilioMessageService{
		Twilio: gotwilio.NewTwilioClient(twilioID, twilioAuthToken),
		Authy:  authy.NewAuthyAPI(authyAPIKey),
	}

	env := streetsweeping.Env{
		MsgSvc: &msgSvc,
	}

	isProduction := os.Getenv("STREETSWEEP_PRODUCTION")

	if isProduction == "true" {
		go func() {
			for range time.Tick(10 * time.Second) {
				streetsweeping.FindReadyAlerts(env.MsgSvc)
			}
		}()
		go func() {
			for range time.Tick(20 * time.Minute) {
				resp, err := http.Get("https://dontfearthesweeper.herokuapp.com/")
				if err != nil {
					log.Println("problem pinging website: ", err)
				}
				if resp.StatusCode != http.StatusOK {
					log.Println("non-200 status code from healthcheck: ", resp.Status)
				}
			}
		}()
	}

	// dontfearthesweeper
	http.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir("./public"))))
	http.Handle("/remove/", http.StripPrefix("/remove/", http.FileServer(http.Dir("./public/remove"))))
	http.Handle("/remove", http.StripPrefix("/remove", http.FileServer(http.Dir("./public/remove"))))
	http.HandleFunc("/verification/start", env.VerificationStartHandler)
	http.HandleFunc("/verification/verify", env.VerificationVerifyHandler)
	http.HandleFunc("/alerts/stop", env.StopAlertHandler)

	// text message directions
	http.HandleFunc("/sms", text_info.ReceiveTextsHandler)

	paths := os.Getenv("PATHS")
	if paths == "" {
		panic("PATHS env variable not set")
	}
	pathsSlice := strings.Split(paths, ",")

	for _, path := range pathsSlice {
		fmt.Println("path: ", path)
		http.HandleFunc("/"+path+"/sms", text_info.ReceiveTextsHandler)
	}
	http.HandleFunc("/directions", text_info.DirectionsHandler) //for demo purposes

	log.Println("Magic happening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
