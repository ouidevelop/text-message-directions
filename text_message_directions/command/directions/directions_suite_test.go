package directions_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDirections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Directions Suite")
}
