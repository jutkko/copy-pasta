package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Main", func() {
	var writeContent []byte

	Describe("no flags", func() {
		Context("when the .copy-pastarc is present", func() {
			var tmpDir, copyPastaRc string
			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
				Expect(err).ToNot(HaveOccurred())

				os.Setenv("S3ENDPOINT", "play.minio.io:9000")
				os.Setenv("S3LOCATION", "us-east-1")
				os.Setenv("HOME", tmpDir)

				// this example uses the test minio endpoint
				copyPastaRc = filepath.Join(userHomeDir(), ".copy-pastarc")
				copyPastaRcContents := `currenttarget:
  name: some-target
  accesskey: Q3AM3UQ867SPQQA43P2F
  secretaccesskey: zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
  bucketname: bucket-name
targets:
  some-target:
    name: some-target
    accesskey: Q3AM3UQ867SPQQA43P2F
    secretaccesskey: zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
    bucketname: bucket-name`
				ioutil.WriteFile(copyPastaRc, []byte(copyPastaRcContents), 0600)
				writeContent = []byte("HHHHHHHHHHey\nBye")
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tmpDir)).To(Succeed())
			})

			It("should store the stdin and return the value next time I call the binary", func() {
				createCmd()
				stdinPipe := getStdinPipe()
				_, err := stdinPipe.Write(writeContent)
				Expect(err).ToNot(HaveOccurred())

				session := runBinary()

				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())
				session.Wait(10 * time.Second)
				Expect(session.ExitCode()).To(Equal(0))

				createCmd()
				session = runBinary()
				session.Wait(10 * time.Second)

				readString := string(session.Out.Contents())
				Expect(readString).To(Equal(string(writeContent)))
				Expect(session.ExitCode()).To(Equal(0))
			})
		})

		Context("when the .copy-pastarc is not present", func() {
			It("prompts to log in", func() {
				createCmd()
				writeContent = []byte("HHHHHHHHHHey\nBye\n")
				stdinPipe := getStdinPipe()
				_, err := stdinPipe.Write(writeContent)
				Expect(err).ToNot(HaveOccurred())
				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())

				session := runBinary()

				Eventually(session.Out).Should(gbytes.Say("Please log in"))
				session.Wait(10 * time.Second)
				Expect(session.ExitCode()).ToNot(Equal(0))
			})
		})
	})

	Describe("when flags are passed", func() {
		Context("login", func() {
			var tmpDir string
			var err error
			BeforeEach(func() {
				tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
				Expect(err).ToNot(HaveOccurred())

				os.Setenv("HOME", tmpDir)
				os.Setenv("S3ENDPOINT", "play.minio.io:9000")
				os.Setenv("S3LOCATION", "us-east-1")
				writeContent = []byte("This is copy-pasta\nBye")
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tmpDir)).To(Succeed())
			})

			// this example uses the test minio endpoint
			It("should prompt for credentials and next time it should work", func() {
				args = []string{"login", "--target", "myTarget"}
				createCmd()
				writeContent = []byte("Q3AM3UQ867SPQQA43P2F\nzuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG\n")
				stdinPipe := getStdinPipe()
				_, err := stdinPipe.Write(writeContent)
				Expect(err).ToNot(HaveOccurred())
				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())

				session := runBinary()

				Eventually(session.Out).Should(gbytes.Say("Please enter key"))
				Eventually(session.Out).Should(gbytes.Say("Please enter secret key"))
				Eventually(session.Out).Should(gbytes.Say("Log in information saved"))

				Expect(filepath.Join(userHomeDir(), ".copy-pastarc")).To(BeAnExistingFile())

				args = []string{}
				createCmd()
				stdinPipe = getStdinPipe()
				_, err = stdinPipe.Write([]byte("something"))
				Expect(err).ToNot(HaveOccurred())

				session = runBinary()

				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())
				session.Wait(10 * time.Second)
				Expect(session.ExitCode()).To(Equal(0))

				createCmd()
				session = runBinary()
				session.Wait(10 * time.Second)

				readString := string(session.Out.Contents())
				Expect(readString).To(Equal("something"))
				Expect(session.ExitCode()).To(Equal(0))
			})
		})

		Context("target", func() {
			var tmpDir string
			var err error
			BeforeEach(func() {
				tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
				Expect(err).ToNot(HaveOccurred())

				os.Setenv("HOME", tmpDir)
				os.Setenv("S3ENDPOINT", "play.minio.io:9000")
				os.Setenv("S3LOCATION", "us-east-1")
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tmpDir)).To(Succeed())
			})

			// this example uses the test minio endpoint
			It("should prompt for credentials and next time it should work", func() {
				args = []string{"login", "--target", "myTargetOne"}
				createCmd()
				loginWriteContent := []byte("Q3AM3UQ867SPQQA43P2F\nzuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG\n")
				stdinPipe := getStdinPipe()
				_, err := stdinPipe.Write(loginWriteContent)
				Expect(err).ToNot(HaveOccurred())
				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())

				session := runBinary()

				Eventually(session.Out).Should(gbytes.Say("Please enter key"))
				Eventually(session.Out).Should(gbytes.Say("Please enter secret key"))
				Eventually(session.Out).Should(gbytes.Say("Log in information saved"))

				Eventually(filepath.Join(userHomeDir(), ".copy-pastarc")).Should(BeAnExistingFile())

				// copy something into it
				args = []string{}
				createCmd()
				stdinPipe = getStdinPipe()
				writeContent = []byte("Hi from targetOne")
				_, err = stdinPipe.Write(writeContent)
				Expect(err).ToNot(HaveOccurred())

				session = runBinary()

				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())
				session.Wait(10 * time.Second)
				Expect(session.ExitCode()).To(Equal(0))

				// login as another target
				args = []string{"login", "--target", "myTargetTwo"}
				createCmd()
				loginWriteContent = []byte("Q3AM3UQ867SPQQA43P2F\nzuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG\n")
				stdinPipe = getStdinPipe()
				_, err = stdinPipe.Write(loginWriteContent)
				Expect(err).ToNot(HaveOccurred())
				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())

				session = runBinary()

				Eventually(session.Out).Should(gbytes.Say("Please enter key"))
				Eventually(session.Out).Should(gbytes.Say("Please enter secret key"))
				Eventually(session.Out).Should(gbytes.Say("Log in information saved"))

				Expect(filepath.Join(userHomeDir(), ".copy-pastarc")).To(BeAnExistingFile())

				// implicitly targeted as myTargetTwo
				// copy something into it
				args = []string{}
				createCmd()
				stdinPipe = getStdinPipe()
				writeContent = []byte("Hi from targetTwo")
				_, err = stdinPipe.Write(writeContent)
				Expect(err).ToNot(HaveOccurred())

				session = runBinary()

				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())
				session.Wait(10 * time.Second)
				Expect(session.ExitCode()).To(Equal(0))

				// set target
				args = []string{"target", "myTargetOne"}
				createCmd()

				session = runBinary()
				session.Wait(10 * time.Second)

				readString := string(session.Out.Contents())
				Expect(readString).To(BeEmpty())
				Expect(session.ExitCode()).To(Equal(0))

				// get hi from targetOne  out
				args = []string{}
				createCmd()

				session = runBinary()
				session.Wait(10 * time.Second)

				readString = string(session.Out.Contents())
				Expect(readString).To(Equal("Hi from targetOne"))
				Expect(session.ExitCode()).To(Equal(0))
			})

			Context("when there is no target", func() {
				It("should display the current target", func() {
					args = []string{"login", "--target", "myTargetOne"}
					createCmd()
					loginWriteContent := []byte("Q3AM3UQ867SPQQA43P2F\nzuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG\n")
					stdinPipe := getStdinPipe()
					_, err := stdinPipe.Write(loginWriteContent)
					Expect(err).ToNot(HaveOccurred())
					err = stdinPipe.Close()
					Expect(err).ToNot(HaveOccurred())

					session := runBinary()
					session.Wait(10 * time.Second)

					args = []string{"target"}
					createCmd()
					session = runBinary()
					Eventually(session.Out).Should(gbytes.Say("myTargetOne"))

					session.Wait(10 * time.Second)
					Expect(session.ExitCode()).To(Equal(0))
				})
			})

		})

		Context("targets", func() {
			var tmpDir string
			var err error
			BeforeEach(func() {
				tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
				Expect(err).ToNot(HaveOccurred())

				os.Setenv("HOME", tmpDir)
				os.Setenv("S3ENDPOINT", "play.minio.io:9000")
				os.Setenv("S3LOCATION", "us-east-1")
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tmpDir)).To(Succeed())
			})

			It("should list the targets", func() {
				args = []string{"login", "--target", "myTargetOne"}
				createCmd()
				loginWriteContent := []byte("Q3AM3UQ867SPQQA43P2F\nzuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG\n")
				stdinPipe := getStdinPipe()
				_, err := stdinPipe.Write(loginWriteContent)
				Expect(err).ToNot(HaveOccurred())
				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())

				session := runBinary()
				session.Wait(10 * time.Second)

				Expect(session.ExitCode()).To(Equal(0))
				Eventually(filepath.Join(userHomeDir(), ".copy-pastarc")).Should(BeAnExistingFile())
				f, err := ioutil.ReadFile(filepath.Join(userHomeDir(), ".copy-pastarc"))
				Expect(err).ToNot(HaveOccurred())
				fmt.Printf("file is: %s", string(f))

				args = []string{"login", "--target", "myTargetTwo"}
				createCmd()
				stdinPipe = getStdinPipe()
				_, err = stdinPipe.Write(loginWriteContent)
				Expect(err).ToNot(HaveOccurred())
				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())

				session = runBinary()
				session.Wait(10 * time.Second)

				Expect(session.ExitCode()).To(Equal(0))
				Eventually(filepath.Join(userHomeDir(), ".copy-pastarc")).Should(BeAnExistingFile())

				args = []string{}
				createCmd()
				stdinPipe = getStdinPipe()
				_, err = stdinPipe.Write([]byte("Hi from targetTwo"))
				Expect(err).ToNot(HaveOccurred())

				session = runBinary()
				err = stdinPipe.Close()
				Expect(err).ToNot(HaveOccurred())
				session.Wait(10 * time.Second)

				Expect(session.ExitCode()).To(Equal(0))

				args = []string{"targets"}
				createCmd()

				session = runBinary()
				session.Wait(10 * time.Second)

				// lists the current target
				Expect(string(session.Out.Contents())).To(ContainSubstring("copy-pasta current target:"))

				Expect(string(session.Out.Contents())).To(ContainSubstring("myTargetOne"))
				Expect(string(session.Out.Contents())).To(ContainSubstring("myTargetTwo"))
				Expect(string(session.Out.Contents())).ToNot(ContainSubstring("Hi from targetTwo"))
			})
		})

		Context("help", func() {
			It("should show some help for the app in stderr", func() {
				args = []string{"help"}
				createCmd()
				session := runBinary()
				Eventually(session.Err).Should(gbytes.Say("Usage: copy-pasta"))

				session.Wait(10 * time.Second)
				Expect(session.ExitCode()).NotTo(Equal(0))
			})
		})

		Context("something invalid", func() {
			It("should inform that the command is not valid", func() {
				args = []string{"ligon", "--target", "myTarget"}
				createCmd()
				session := runBinary()
				Eventually(session.Err).Should(gbytes.Say("Usage: copy-pasta"))

				session.Wait(10 * time.Second)
				Expect(session.ExitCode()).ToNot(Equal(0))
			})
		})
	})
})

func userHomeDir() string {
	return os.Getenv("HOME")
}
