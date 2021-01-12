package parser_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArithmetic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arithmetic Suite")
}
