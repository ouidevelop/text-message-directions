package streetsweeping

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"database/sql"

	"net/http/httputil"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Env contains the interfaces for any external API's used. This way we can mock out those API's in tests.
type Env struct {
	MsgSvc MessageServicer
}

var (
	//All global environment variables should be set at the beginning of the application, then remain unchanged.

	// DB is a database handle representing a pool of zero or more
	// underlying connections. It's safe for concurrent use by multiple
	// goroutines.
	// **(from sql package)**
	DB   *sql.DB
	from string
)

type startVerification struct {
	Via         string `json:"via"`
	PhoneNumber string `json:"phoneNumber"`
}

type alert struct {
	Timezone    string `json:"timezone"`
	Times       []day  `json:"times"`
	PhoneNumber string `json:"phoneNumber"`
	Token       string `json:"token"`
}

type removeAlert struct {
	PhoneNumber string `json:"phoneNumber"`
	Token       string `json:"token"`
}

type day struct {
	Weekday int `json:"weekday"`
	NthWeek int `json:"nthWeek"`
}

func init() {
	from = os.Getenv("TWILIO_PHONE_NUMBER")
	if from == "" {
		log.Fatal("TWILIO_PHONE_NUMBER environment variable not set")
	}

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" && os.Getenv("STREETSWEEP_PRODUCTION") != "" {
		log.Fatal("MYSQL_PASSWORD environment variable not set")
	}
	if mysqlPassword != "" {
		DB = StartDB(mysqlPassword)
	}
}

func (env *Env) StopAlertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in stopAlertHandler")
	decoder := json.NewDecoder(r.Body)
	var t removeAlert
	err := decoder.Decode(&t)
	log.Println("1")
	if err != nil {
		log.Println("2")
		log.Println("error decoding json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "oops! we made a mistake")
		return
	}
	log.Println("3")
	defer r.Body.Close()

	verification, err := env.MsgSvc.VerifyCode(t.PhoneNumber, t.Token)
	if err != nil {
		log.Println("4")
		log.Println("error verifying code: error: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "validation code incorrect")
		return
	}
	log.Println("5")
	if !verification {
		log.Println("6")
		//todo: do this better. figure out all the ways that CheckPhoneVerification could fail and handle all of them well
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "validation code incorrect")
		return
	}

	log.Println("7")
	err = removeAlerts(t)
	if err != nil {
		log.Println("8")
		log.Println("problem deleting alert to database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "oops! we made a mistake")
		return
	}
	log.Println("9")

	w.WriteHeader(http.StatusOK)
}

func (env *Env) VerificationStartHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))

	decoder := json.NewDecoder(r.Body)
	var t startVerification
	err = decoder.Decode(&t)
	if err != nil {
		log.Println("error decoding json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "oops! we made a mistake")
		return
	}
	defer r.Body.Close()

	verified, err := env.MsgSvc.RequestCode(t.PhoneNumber)
	if !verified {
		//todo: do this better. figure out all the ways that start phone verification could fail and handle all of them well
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "problem starting phone verification")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// VerificationVerifyHandler verifies that a user has the correct verification code.
func (env *Env) VerificationVerifyHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println("VerificationVerifyHandler", string(requestDump))

	decoder := json.NewDecoder(r.Body)
	var t alert
	err = decoder.Decode(&t)
	if err != nil {
		log.Println("error decoding json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "oops! we made a mistake")
		return
	}
	defer r.Body.Close()

	verified, err := env.MsgSvc.VerifyCode(t.PhoneNumber, t.Token)
	if !verified {
		//todo: do this better. figure out all the ways that CheckPhoneVerification could fail and handle all of them well
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "validation code incorrect")
		return
	}

	err = save(t)
	if err != nil {
		log.Println("problem saving new alert to database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "oops! we made a mistake")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Now provides a rapper to time.Now and can be used to mock calls to time.Now in tests.
var Now = func() time.Time {
	return time.Now()
}

// CalculateNextCall takes an nth week (first, second, third, forth), a weekday, and a timezone and calculates
// the next time that a person should be alerted for street sweeping.
func CalculateNextCall(nthWeek int, weekday int, timezone string) (int64, error) {

	var NextCallUnixTime int64

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return NextCallUnixTime, err
	}

	now := Now().In(location)
	timeToSendMessageThisMonth := timeAtNthDayOfMonth(now, nthWeek, weekday, 19)
	if now.After(timeToSendMessageThisMonth) { // if right now is after the time to send a message this month
		dateNextMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, 0)
		timeToSendMessageThisMonth = timeAtNthDayOfMonth(dateNextMonth, nthWeek, weekday, 19)
		if now.After(timeToSendMessageThisMonth) {
			dateInTwoMonths := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 2, 0)
			timeToSendMessageThisMonth = timeAtNthDayOfMonth(dateInTwoMonths, nthWeek, weekday, 19)
		}
	}

	NextCallUnixTime = timeToSendMessageThisMonth.Unix()

	return NextCallUnixTime, nil
}

func timeAtNthDayOfMonth(t time.Time, nthDay int, weekday int, hour int) time.Time {
	firstDayOfThisMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	dateOfFirstWeekday := ((weekday+7)-int(firstDayOfThisMonth.Weekday()))%7 + 1
	dateOfNthWeekday := dateOfFirstWeekday + ((nthDay - 1) * 7)
	TimeAtNthDayOfMonth := time.Date(t.Year(), t.Month(), dateOfNthWeekday, hour, 0, 0, 0, t.Location())
	return TimeAtNthDayOfMonth.Add(-24 * time.Hour)
}

func remind(phoneNumber string, sender smsMessager, id int) {
	log.Println("sending message to: ", id)
	message := "Don't forget about street sweeping tomorrow! (to stop getting these reminders, go to dontfearthesweeper.com/remove or email ouidevelop@gmail.com)"
	err := sender.Send(from, phoneNumber, message)
	if err != nil {
		log.Println("problem sending message: ", err)
	}
}

// todo: get rid of /<wbr/>
// test locally.
// stop super long messages.
