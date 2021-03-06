package sms

import (
	"log"
	"os"
	"strconv"
	"strings"
	"fmt"

	"github.com/sfreiberg/gotwilio"
)

var (
	twilioAuths = make(map[string]string)
)

func init(){
	envVars := os.Environ()
	for _, envVar := range envVars {
		parts := strings.Split(envVar, "=")
		key := parts[0]
		value := parts[1]
		if strings.HasPrefix(key, "TWILIO_AUTHS_") {
			parts := strings.Split(value, ":")
			twilioAuths[parts[0]] = parts[1]
		}
	}
}

func Send(from string, to string, id string, body string) {

	cl := gotwilio.NewTwilioClient(id, twilioAuths[id])

	messages := SplitLongBody(body)
	
	for _, message := range messages {
		fmt.Println("bob", len(message))
		res, exception, err := cl.SendSMS(from, to, message, "", "")
		log.Printf("res: %+v, exeption: %+v, err: %+v", res, exception, err)
		if exception != nil {
			errorMessage := "oops! we have an error with code " +
				strconv.Itoa(int(exception.Code)) +
				". If you'd like help, please share this code with ouidevelop@gmail.com"
			res, exception, err = cl.SendSMS(from, to, errorMessage, "", "")
			log.Printf("res: %+v, exeption: %+v, err: %+v", res, exception, err)
		}
	}
}

func SplitLongBody(str string) []string {
	var messages []string
	segments := strings.Split(str, "\n\n")

	message := ""
	for _, segment := range segments {
		if len(message + segment + "\n\n") < 1550 {
			message = message + segment + "\n\n"
		} else {
			messages = append(messages, message)
			// message = segment + "\n\n"
			message = ""
		}
	}
	
	if message != "" {
		messages = append(messages, message)
	}

	return messages
}