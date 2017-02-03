package runcommands_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jutkko/copy-pasta/runcommands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rc", func() {
	Describe("Load", func() {
		var tmpDir, copyPastaRc string
		Context("when the .copy-pastarc file does not exist", func() {
			It("should return an error saying there is no copy-pastarc", func() {
				_, err := runcommands.Load()
				Expect(err.Error()).To(ContainSubstring("Unable to load the targets, please check if ~/.copy-pastarc exists"))
			})
		})

		Context("when the .copy-pastarc file exists", func() {
			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
				Expect(err).ToNot(HaveOccurred())

				os.Setenv("HOME", tmpDir)

				copyPastaRc = filepath.Join(userHomeDir(), ".copy-pastarc")
				copyPastaRcContents := `some-target:
  accesskey: some-key
  secretaccesskey: some-secret-key
another-target:
  accesskey: another-key
  secretaccesskey: another-secret-key`
				ioutil.WriteFile(copyPastaRc, []byte(copyPastaRcContents), 0600)
			})

			It("should load the target to the Rc struct", func() {
				targets, err := runcommands.Load()
				Expect(err).ToNot(HaveOccurred())

				Expect(targets["some-target"].AccessKey).To(Equal("some-key"))
				Expect(targets["some-target"].SecretAccessKey).To(Equal("some-secret-key"))
				Expect(targets["another-target"].AccessKey).To(Equal("another-key"))
				Expect(targets["another-target"].SecretAccessKey).To(Equal("another-secret-key"))
			})
		})

		Context("when the .copy-pastarc is not valid", func() {
			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
				Expect(err).ToNot(HaveOccurred())

				os.Setenv("HOME", tmpDir)

				copyPastaRc = filepath.Join(userHomeDir(), ".copy-pastarc")
				copyPastaRcContents := `some-target:
	accewhaa: some-key
  whaaaa: some-secret-key`
				ioutil.WriteFile(copyPastaRc, []byte(copyPastaRcContents), 0600)
			})

			It("should return an parsing error", func() {
				_, err := runcommands.Load()
				Expect(err.Error()).To(ContainSubstring("Parsing failed"))
			})
		})
	})

	FDescribe("Update", func() {
		var tmpDir, copyPastaRc string

		Context("when there is a target file already", func() {
			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
				Expect(err).ToNot(HaveOccurred())

				os.Setenv("HOME", tmpDir)

				copyPastaRc = filepath.Join(userHomeDir(), ".copy-pastarc")
				copyPastaRcContents := `some-target:
  accesskey: some-key
  secretaccesskey: some-secret-key`
				ioutil.WriteFile(copyPastaRc, []byte(copyPastaRcContents), 0600)
			})

			It("updates the current .copy-pastarc", func() {
				err := runcommands.Update("another-target", "another-accesskey", "another-secret-key")
				Expect(err).ToNot(HaveOccurred())

				targets, err := runcommands.Load()
				Expect(err).ToNot(HaveOccurred())
				Expect(targets["another-target"].AccessKey).To(Equal("another-accesskey"))
				Expect(targets["another-target"].SecretAccessKey).To(Equal("another-secret-key"))
			})
		})

		Context("when there is no target to start with", func() {
			It("should create a new .copy-pasta file with the passed in credentials", func() {
			})
		})
	})
})

func userHomeDir() string {
	return os.Getenv("HOME")
}
