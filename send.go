package main

import (
	"github.com/sfreiberg/gotwilio"
	"log"
	"strconv"
	"strings"
)

func send(cl *gotwilio.Twilio, from string, to string, body string) {

	var messages []string

	if len(body) < 1500 {
		messages = append(messages, body)
	} else {
		messages = splitLongBody(body)
	}

	for _, message := range messages {
		res, exception, err := cl.SendSMS(from, to, message, "", "")
		log.Printf("res: %+v, exeption: %+v, err: %+v", res, exception, err)
		if exception != nil {
			errorMessage := "oops! we have an error with code " +
				strconv.Itoa(exception.Code) +
				". If you'd like help, please share this code with ouidevelop@gmail.com"
			res, exception, err = cl.SendSMS(from, to, errorMessage, "", "")
			log.Printf("res: %+v, exeption: %+v, err: %+v", res, exception, err)
		}
	}
}

func splitLongBody(str string) []string {
	var messages []string
	segments := strings.Split(str, "\n\n")

	charCount := 0
	message := ""
	for _, segment := range segments {
		if charCount < 1500 {
			message += segment + "\n\n"
			charCount += len(segment + "\n\n")
		} else {
			messages = append(messages, message)
			charCount = 0
			message = ""
		}
	}

	return messages
}
