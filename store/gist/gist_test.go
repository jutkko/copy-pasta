package gist_test

import (
	"bytes"
	"errors"

	"github.com/google/go-github/github"
	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/jutkko/copy-pasta/store/gist"
	"github.com/jutkko/copy-pasta/store/gist/gistfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("S3", func() {
	Describe("Write", func() {
		var (
			fakeGistClient              *gistfakes.FakeGistClient
			actualTarget                *runcommands.Target
			writeContent, token, gistID string
			testStore                   *gist.GistStore
		)

		BeforeEach(func() {
			fakeGistClient = new(gistfakes.FakeGistClient)
			writeContent = "This is a secret"
			token = "my-token"
			gistID = "gist-IDID"
		})

		Context("when the target has gitID", func() {
			It("should modify the existing gist", func() {
				actualTarget = &runcommands.Target{
					Backend:   "gist",
					GistToken: token,
					GistID:    gistID,
				}

				testStore = gist.NewGistStore(fakeGistClient, actualTarget)

				content := "Arrgggh!\nOooops"
				contentReader := bytes.NewReader([]byte(content))
				filename := "copy-pasta"

				err := testStore.Write(contentReader)
				Expect(err).ToNot(HaveOccurred())
				_, gistIDCall, gistCall := fakeGistClient.EditArgsForCall(0)
				Expect(gistIDCall).To(Equal(gistID))
				Expect(*gistCall.Files[github.GistFilename(filename)].Content).To(Equal("Arrgggh!\nOooops"))
			})

			Context("when editing fails", func() {
				It("should try to create a new gist", func() {
					actualTarget = &runcommands.Target{
						Name:      "test-gist",
						Backend:   "gist",
						GistToken: token,
						GistID:    gistID,
					}

					newGistID := "myGistID"
					fakeGistClient.CreateReturns(&github.Gist{ID: &newGistID}, nil, nil)
					fakeGistClient.EditReturns(nil, nil, errors.New("Cannot edit"))
					testStore = gist.NewGistStore(fakeGistClient, actualTarget)
					content := "Arrgggh!\nOooops"
					contentReader := bytes.NewReader([]byte(content))

					err := testStore.Write(contentReader)
					Expect(err).NotTo(HaveOccurred())
					_, gistCall := fakeGistClient.CreateArgsForCall(0)
					Expect(*gistCall.Files[github.GistFilename("copy-pasta")].Content).To(Equal("Arrgggh!\nOooops"))
				})
			})
		})

		Context("when the target does not have a gitID", func() {
			It("should create a new gist, and update the target", func() {
				actualTarget = &runcommands.Target{
					Backend:   "gist",
					GistToken: token,
				}

				testStore = gist.NewGistStore(fakeGistClient, actualTarget)

				content := "Arrgggh!\nOooops"
				contentReader := bytes.NewReader([]byte(content))
				filename := "copy-pasta"
				newGistID := "myGistID"
				gist := &github.Gist{ID: &newGistID}
				fakeGistClient.CreateReturns(gist, nil, nil)

				err := testStore.Write(contentReader)
				Expect(err).ToNot(HaveOccurred())
				_, gistCall := fakeGistClient.CreateArgsForCall(0)
				Expect(*gistCall.Files[github.GistFilename(filename)].Content).To(Equal("Arrgggh!\nOooops"))
			})

			Context("when creating fails", func() {
				It("should fail with error", func() {
					actualTarget = &runcommands.Target{
						Backend:   "gist",
						GistToken: token,
					}

					testStore = gist.NewGistStore(fakeGistClient, actualTarget)

					fakeGistClient.CreateReturns(nil, nil, errors.New("Cannot create"))
					testStore = gist.NewGistStore(fakeGistClient, actualTarget)
					content := "Arrgggh!\nOooops"
					contentReader := bytes.NewReader([]byte(content))

					err := testStore.Write(contentReader)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("Read", func() {
		var (
			fakeGistClient       *gistfakes.FakeGistClient
			actualTarget         *runcommands.Target
			actualContent, token string
			testStore            *gist.GistStore
		)

		BeforeEach(func() {
			fakeGistClient = new(gistfakes.FakeGistClient)
			actualContent = "This is a secret.\n"
			token = "my-token"
			actualTarget = &runcommands.Target{
				Backend:   "gist",
				GistToken: token,
				GistID:    "gist-IDID",
			}

			testStore = gist.NewGistStore(fakeGistClient, actualTarget)
		})

		It("should return the string", func() {
			rawURL := "https://gist.githubusercontent.com/jutkko/c58b1318786d7b89ff224dcc3ad87ac1/raw/21226c78ff3bd4422544fd05be8a816216f00897/copy-pasta"
			gist := &github.Gist{
				Files: map[github.GistFilename]github.GistFile{
					"copy-pasta": github.GistFile{
						RawURL: &rawURL,
					},
				},
			}
			fakeGistClient.GetReturns(gist, nil, nil)

			content, err := testStore.Read()
			_, gistID := fakeGistClient.GetArgsForCall(0)
			Expect(gistID).To(Equal("gist-IDID"))
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(actualContent))
		})

		It("should return the error when get fails", func() {
			fakeGistClient.GetReturns(nil, nil, errors.New("Failed to fetch gist"))

			content, err := testStore.Read()
			_, gistID := fakeGistClient.GetArgsForCall(0)
			Expect(gistID).To(Equal("gist-IDID"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to fetch gist"))
			Expect(content).To(Equal(""))
		})
	})
})
