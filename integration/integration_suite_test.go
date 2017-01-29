package integration_test

import (
	"io"
	"log"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var args []string
var cmd *exec.Cmd
var err error
var execPath string

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	srcPath := "github.com/jutkko/copy-pasta"
	execPath, err = gexec.Build(srcPath)
	if err != nil {
		log.Fatalf("executable %s could not be built: %s", srcPath, err)
	}
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func createCmd() {
	cmd = exec.Command(execPath, args...)
}

func getStdinPipe() io.WriteCloser {
	stdinPipe, err := cmd.StdinPipe()
	Expect(err).ToNot(HaveOccurred())

	return stdinPipe
}

func runBinary() *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())

	return session
}
