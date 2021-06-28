package streetsweeping_test

import (
	"log"

	. "github.com/ouidevelop/dontfearthesweeper/streetsweeping"

	"time"

	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Describe("application", func() {
		It("should alert people who's alert is past due", func() {
			err := DB.Ping()
			Expect(err).NotTo(HaveOccurred())

			jsonAlert := []byte(`{"timezone":"America/New_York","times":[{"weekday":0,"nthWeek":1}],"phoneNumber":"1234567890","token":""}`)

			clearDB()
			defer clearDB()

			req := httptest.NewRequest("POST", "/verification/verify", bytes.NewReader(jsonAlert))
			res := httptest.NewRecorder()
			MockEnv.VerificationVerifyHandler(res, req)
			Expect(res.Code).To(Equal(http.StatusOK))

			var nextCall int64
			err = DB.QueryRow("select NEXT_CALL from alerts").Scan(&nextCall)
			Expect(err).NotTo(HaveOccurred())
			Expect(nextCall).To(Equal(int64(1494111600))) //2017-05-06 19:00:00 -0400 EDT

			// 1 second after the next call of the alert
			done := MockNow(time.Unix(1494111601, 0))
			defer done()

			FindReadyAlerts(MockEnv.MsgSvc)
			expected := &MockMessageService{
				from: "5102414070",
				to:   "1234567890",
				body: "Don't forget about street sweeping tomorrow! (to stop getting these reminders, please email mjkurrels@gmail.com)",
			}
			Expect(MockEnv.MsgSvc).To(Equal(expected))
		})
	})

	Describe("CalculateNextCall", func() {
		It("should calculate the next date to send an alert", func() {
			location, err := time.LoadLocation("America/Los_Angeles")
			Now = func() time.Time {
				return time.Date(2017, 7, 31, 19, 1, 1, 0, location)
			}
			Expect(err).NotTo(HaveOccurred())

			weekday := 2
			nthWeek := 1
			timezone := "America/Los_Angeles"

			log.Println("hi:", Now())
			nextAlertTime, err := CalculateNextCall(nthWeek, weekday, timezone)
			Expect(err).NotTo(HaveOccurred())
			Expect(nextAlertTime).To(Equal(int64(1504576800))) //2017-09-04 19:00:00 -0700

			// change the weekday to friday to make sure that the function determines that the next alert is this month
			weekday = 5

			nextAlertTime, err = CalculateNextCall(nthWeek, weekday, timezone)
			log.Println("time: ", time.Unix(nextAlertTime, 0).In(location))
			Expect(err).NotTo(HaveOccurred())
			Expect(nextAlertTime).To(Equal(int64(1501812000))) //2017-08-03 19:00:00 -0700 PDT
		})
	})
})

func clearDB() {
	_, err := DB.Exec("Truncate table alerts")
	Expect(err).NotTo(HaveOccurred())
}
