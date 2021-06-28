package places_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ouidevelop/dontfearthesweeper/text_message_directions/command/places"
)

var _ = Describe("Places", func() {
	It("should be able to get the phone number and address of a business", func() {
		placeString, err := places.Get("info for starbucks slo")
		Expect(err).NotTo(HaveOccurred())
		Expect(placeString).NotTo(Equal(""))
	})
})
