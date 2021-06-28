package directions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ouidevelop/dontfearthesweeper/text_message_directions/command/directions"
)

var _ = Describe("Directions", func() {
	//It("should get an error if you don't specify a mode", func(){
	//	d, err := directions.Get("berkeley to oakland")
	//	Expect(err).To(HaveOccurred())
	//	Expect(err.Error()).To(Equal("directions malformed. Please submit in the form '[mode] from [origin] to [destination]'. Mode can be either 'walk', 'drive', 'bike', or 'transit'. For example: 'drive from berkeley to sfo'"))
	//	Expect(d).To(Equal(""))
	//})
	//It("should get an error if you specify an incorrect mode", func(){
	//	d, err := directions.Get("badmode berkeley to oakland")
	//	Expect(err).To(HaveOccurred())
	//	Expect(err.Error()).To(Equal("directions malformed. Please submit in the form '[mode] from [origin] to [destination]'. Mode can be either 'walk', 'drive', 'bike', or 'transit'. For example: 'drive from berkeley to sfo'"))
	//	Expect(d).To(Equal(""))
	//})
	It("should be able to get transit directions", func(){
		d, err := directions.Get("transit from uc berkeley to oakland airport")
		Expect(err).NotTo(HaveOccurred())
		Expect(d).NotTo(Equal(""))
	})

	It("should be able to get walking directions", func(){
		d, err := directions.Get("walk from uc berkeley to moe's books")
		Expect(err).NotTo(HaveOccurred())
		Expect(d).NotTo(Equal(""))
	})

	It("should be able to get driving directions", func(){
		d, err := directions.Get("drive from uc berkeley to moe's books")
		Expect(err).NotTo(HaveOccurred())
		Expect(d).NotTo(Equal(""))
	})

	It("should be able to get driving directions", func(){
		d, err := directions.Get("drive from Clarksville to Nashville")
		Expect(err).NotTo(HaveOccurred())
		Expect(d).NotTo(Equal(""))
	})

	It("should be able to get bike directions", func(){
		d, err := directions.Get("bike from uc berkeley to moe's books")
		Expect(err).NotTo(HaveOccurred())
		Expect(d).NotTo(Equal(""))
	})

})
