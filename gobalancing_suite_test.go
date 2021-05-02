package gobalancing_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGobalancing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobalancing Suite")
}
