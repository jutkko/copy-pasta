package integration_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Main", func() {

	PIt("should run successfully with exit code 0", func() {
		session := runBinary()
		fmt.Printf("Stdout: %+#v", session.Out.Contents())
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("should store the stdin and return the value next time I call the binary", func() {
		createCmd()
		stdinPipe := getStdinPipe()
		session := runBinary()

		_, err := stdinPipe.Write([]byte("hey"))
		Expect(err).ToNot(HaveOccurred())

		err = stdinPipe.Close()
		Expect(err).ToNot(HaveOccurred())

		createCmd()
		session = runBinary()
		Eventually(session.Out).Should(gbytes.Say("hey"))
	})
})
