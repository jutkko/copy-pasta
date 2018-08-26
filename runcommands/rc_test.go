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
		BeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
			Expect(err).ToNot(HaveOccurred())

			os.Setenv("HOME", tmpDir)
		})

		Context("when the .copy-pastarc file does not exist", func() {
			It("should return an error saying there is no copy-pastarc", func() {
				_, err := runcommands.Load()
				Expect(err.Error()).To(ContainSubstring("Unable to load the targets, please check if ~/.copy-pastarc exists"))
			})
		})

		Context("when the .copy-pastarc file exists", func() {
			BeforeEach(func() {
				copyPastaRc = filepath.Join(userHomeDir(), ".copy-pastarc")
				copyPastaRcContents := `currenttarget:
  name: mycurrenttarget
  backend: s3
  accesskey: current-key
  secretaccesskey: current-secret-key
  bucketname: current-bucket-name
  endpoint: current-endpoint
  location: current-location
targets:
  mycurrenttarget:
    name: mycurrenttarget
    backend: s3
    accesskey: current-key
    secretaccesskey: current-secret-key
    bucketname: current-bucket-name
    endpoint: current-endpoint
    location: current-location
  another-target:
    name: another-target
    backend: gist
    gisttoken: 0xdeadbeef`
				ioutil.WriteFile(copyPastaRc, []byte(copyPastaRcContents), 0600)
			})

			It("should load the target to the Rc struct", func() {
				config, err := runcommands.Load()
				Expect(err).ToNot(HaveOccurred())

				currentTarget := config.CurrentTarget
				checkTarget(currentTarget, "mycurrenttarget", "s3", "current-key", "current-secret-key", "current-bucket-name", "current-endpoint", "current-location", "", "")

				targets := config.Targets
				checkTarget(targets["mycurrenttarget"], "mycurrenttarget", "s3", "current-key", "current-secret-key", "current-bucket-name", "current-endpoint", "current-location", "", "")
				checkTarget(targets["another-target"], "another-target", "gist", "", "", "", "", "", "0xdeadbeef", "")
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

	Describe("Update", func() {
		var tmpDir, copyPastaRc string
		BeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir("", "copy-pasta-test")
			Expect(err).ToNot(HaveOccurred())

			os.Setenv("HOME", tmpDir)

			copyPastaRc = filepath.Join(userHomeDir(), ".copy-pastarc")
		})

		Context("when there is a target file already", func() {
			BeforeEach(func() {
				copyPastaRcContents := `currenttarget:
  name: some-target
  backend: s3
  accesskey: some-key
  secretaccesskey: some-secret-key
  bucketname: some-bucket-name
  endpoint: some-endpoint
  location: some-location
targets:
  some-target:
    name: some-target
    backend: s3
    accesskey: some-key
    secretaccesskey: some-secret-key
    bucketname: some-bucket-name
    endpoint: some-endpoint
    location: some-location`
				ioutil.WriteFile(copyPastaRc, []byte(copyPastaRcContents), 0600)
			})

			It("updates the current .copy-pastarc and sets the current target to target", func() {
				err := runcommands.Update("another-target", "gist", "", "", "", "", "", "hiyo", "my-gist-id")
				Expect(err).ToNot(HaveOccurred())

				config, err := runcommands.Load()
				Expect(err).ToNot(HaveOccurred())

				currentTarget := config.CurrentTarget
				checkTarget(currentTarget, "another-target", "gist", "", "", "", "", "", "hiyo", "my-gist-id")

				targets := config.Targets
				Expect(len(targets)).To(Equal(2))
				checkTarget(targets["some-target"], "some-target", "s3", "some-key", "some-secret-key", "some-bucket-name", "some-endpoint", "some-location", "", "")
				checkTarget(targets["another-target"], "another-target", "gist", "", "", "", "", "", "hiyo", "my-gist-id")
			})
		})

		Context("when there is a target file already but corrupted", func() {
			BeforeEach(func() {
				copyPastaRcContents := `currenttarget:
  some-target:
		accesskey: some-key
	  secretaccesskey: some-secret-key
targets:
  some-target:
    name: some-target
    accesskey: some-key
    secretaccesskey: some-secret-key`
				ioutil.WriteFile(copyPastaRc, []byte(copyPastaRcContents), 0600)
			})

			It("creates a new .copy-pastarc", func() {
				err := runcommands.Update("another-target", "s3", "another-key", "another-secret-key", "another-bucket-name", "another-endpoint", "another-location", "", "")
				Expect(err).ToNot(HaveOccurred())

				config, err := runcommands.Load()
				Expect(err).ToNot(HaveOccurred())

				currentTarget := config.CurrentTarget
				checkTarget(currentTarget, "another-target", "s3", "another-key", "another-secret-key", "another-bucket-name", "another-endpoint", "another-location", "", "")

				targets := config.Targets
				Expect(len(targets)).To(Equal(1))
				checkTarget(targets["another-target"], "another-target", "s3", "another-key", "another-secret-key", "another-bucket-name", "another-endpoint", "another-location", "", "")

				Expect(filepath.Join(userHomeDir(), ".copy-pastarc")).To(BeAnExistingFile())
			})
		})

		Context("when there is no target to start with", func() {
			It("should create a new .copy-pasta file with the passed in credentials", func() {
				err := runcommands.Update("some-target", "s3", "some-key", "some-secret-key", "some-bucket-name", "some-endpoint", "some-location", "", "")
				Expect(err).ToNot(HaveOccurred())

				config, err := runcommands.Load()
				Expect(err).ToNot(HaveOccurred())

				currentTarget := config.CurrentTarget
				checkTarget(currentTarget, "some-target", "s3", "some-key", "some-secret-key", "some-bucket-name", "some-endpoint", "some-location", "", "")

				targets := config.Targets
				Expect(len(targets)).To(Equal(1))
				checkTarget(targets["some-target"], "some-target", "s3", "some-key", "some-secret-key", "some-bucket-name", "some-endpoint", "some-location", "", "")
				Expect(filepath.Join(userHomeDir(), ".copy-pastarc")).To(BeAnExistingFile())
			})
		})
	})
})

func checkTarget(t *runcommands.Target, name, backend, accessKey, secretAccessKey, bucketName, endpoint, location, gistToken, gistID string) {
	Expect(t.Name).To(Equal(name))
	Expect(t.Backend).To(Equal(backend))
	Expect(t.AccessKey).To(Equal(accessKey))
	Expect(t.SecretAccessKey).To(Equal(secretAccessKey))
	Expect(t.BucketName).To(Equal(bucketName))
	Expect(t.Endpoint).To(Equal(endpoint))
	Expect(t.Location).To(Equal(location))
	Expect(t.GistToken).To(Equal(gistToken))
	Expect(t.GistID).To(Equal(gistID))
}

func userHomeDir() string {
	return os.Getenv("HOME")
}
