package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Main", func() {

	It("should run successfully with exit code 0", func() {
		createCmd()
		session := runBinary()
		session.Wait()
		Expect(session.ExitCode()).To(Equal(0))
	})

	FIt("should store the stdin and return the value next time I call the binary", func() {
		createCmd()
		stdinPipe := getStdinPipe()
		_, err := stdinPipe.Write([]byte("hey"))
		Expect(err).ToNot(HaveOccurred())

		session := runBinary()

		err = stdinPipe.Close()
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("Storing"))
		session.Wait()
		Expect(session.ExitCode()).To(Equal(0))

		createCmd()
		session = runBinary()
		Eventually(session.Out).Should(gbytes.Say("hey"))
		session.Wait()
		Expect(session.ExitCode()).To(Equal(0))
	})
})
