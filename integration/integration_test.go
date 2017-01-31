package integration_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	var writeContent []byte

	BeforeEach(func() {
		writeContent = []byte("HHHHHHHHHHey\nBye")
	})

	It("should run successfully with exit code 0", func() {
		createCmd()
		session := runBinary()
		session.Wait(5 * time.Second)
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("should store the stdin and return the value next time I call the binary", func() {
		createCmd()
		stdinPipe := getStdinPipe()
		_, err := stdinPipe.Write(writeContent)
		Expect(err).ToNot(HaveOccurred())

		session := runBinary()

		err = stdinPipe.Close()
		Expect(err).ToNot(HaveOccurred())
		session.Wait(5 * time.Second)
		Expect(session.ExitCode()).To(Equal(0))

		createCmd()
		session = runBinary()
		session.Wait(5 * time.Second)

		readString := string(session.Out.Contents())
		Expect(readString).To(Equal(string(writeContent)))
		Expect(session.ExitCode()).To(Equal(0))
	})
})
