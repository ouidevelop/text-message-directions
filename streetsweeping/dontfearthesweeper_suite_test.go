package streetsweeping_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ouidevelop/dontfearthesweeper/streetsweeping"

	"testing"
	"time"
)

var MockEnv = Env{
	MsgSvc: &MockMessageService{},
}

func TestDontfearthesweeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dontfearthesweeper Suite")
}

type MockMessageService struct {
	from string
	to   string
	body string
}

func (t *MockMessageService) Send(from, to, body string) error {
	t.from = from
	t.to = to
	t.body = body
	return nil
}

func (t *MockMessageService) RequestCode(phoneNumber string) (bool, error) {
	return true, nil
}

func (t *MockMessageService) VerifyCode(phoneNumber, code string) (bool, error) {
	return true, nil
}

func MockNow(t time.Time) func() {
	oldNow := Now
	Now = func() time.Time {
		return t
	}
	doneFunc := func() {
		Now = oldNow
	}
	return doneFunc
}

var _ = BeforeSuite(func() {
	location, _ := time.LoadLocation("America/New_York")
	Now = func() time.Time {
		return time.Date(2017, 4, 6, 0, 0, 0, 0, location)
	}
})
