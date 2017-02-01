package runcommands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRuncommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runcommands Suite")
}
