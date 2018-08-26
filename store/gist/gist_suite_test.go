package gist_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGist(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gist Suite")
}
