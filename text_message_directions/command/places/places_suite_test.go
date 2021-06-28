package places_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlaces(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Places Suite")
}
