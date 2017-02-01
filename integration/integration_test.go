package integration_test

import (
	"os/user"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
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

	Context("when the .copy-pastarc is not present", func() {
		FIt("should prompt for the credentials", func() {
			createCmd()
			writeContent = []byte("HHHHHHHHHHey\nBye\n")
			stdinPipe := getStdinPipe()
			_, err := stdinPipe.Write(writeContent)
			Expect(err).ToNot(HaveOccurred())
			err = stdinPipe.Close()
			Expect(err).ToNot(HaveOccurred())

			session := runBinary()

			Eventually(session.Out).Should(gbytes.Say("Please input your S3 accesskey"))
			Eventually(session.Out).Should(gbytes.Say("Please input your S3 secret accesskey"))

			session.Wait(5 * time.Second)
			Expect(session.ExitCode()).To(Equal(0))

			usr, err := user.Current()
			Expect(err).ToNot(HaveOccurred())

			Expect(filepath.Join(usr.HomeDir, ".copy-pastarc")).To(BeAnExistingFile())
		})
	})
})
