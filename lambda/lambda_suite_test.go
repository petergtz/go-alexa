package lambda_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLambda(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lambda Suite")
}
