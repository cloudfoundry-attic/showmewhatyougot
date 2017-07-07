package statedetector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStateDetector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StateDetector Suite")
}
