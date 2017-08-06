package store_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"github.com/jutkko/copy-pasta/store"
	"github.com/jutkko/copy-pasta/store/storefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("S3", func() {
	Describe("Write", func() {
		var (
			fakeClient                                         *storefakes.FakeMinioClient
			exampleContent                                     io.Reader
			actualBucketName, actualObjectName, actualLoaction string
			testStore                                          *store.S3Store
		)

		BeforeEach(func() {
			exampleContent = bytes.NewReader([]byte("He is a banana\nand an apple"))
			fakeClient = new(storefakes.FakeMinioClient)
			actualBucketName = "this-bucket"
			actualObjectName = "this-object"
			actualLoaction = "that-location"
			testStore = &store.S3Store{MinioClient: fakeClient}
		})

		Context("when the bucket exists command returns an error", func() {
			BeforeEach(func() {
				fakeClient.BucketExistsReturns(true, errors.New("No action should be taken"))
			})

			It("should return the error", func() {
				err := testStore.Write(actualBucketName, actualObjectName, actualLoaction, exampleContent)
				Expect(err).To(MatchError("No action should be taken"))
			})
		})

		Context("when the bucket doesn't exist", func() {
			BeforeEach(func() {
				fakeClient.BucketExistsReturns(false, nil)
			})

			It("should create it and put the object there", func() {
				err := testStore.Write(actualBucketName, actualObjectName, actualLoaction, exampleContent)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
				bucketName := fakeClient.BucketExistsArgsForCall(0)
				Expect(bucketName).To(Equal(actualBucketName))

				Expect(fakeClient.MakeBucketCallCount()).To(Equal(1))
				bucketName, location := fakeClient.MakeBucketArgsForCall(0)
				Expect(bucketName).To(Equal(actualBucketName))
				Expect(location).To(Equal(actualLoaction))

				Expect(fakeClient.PutObjectCallCount()).To(Equal(1))
				bucketName, objectName, reader, contentType := fakeClient.PutObjectArgsForCall(0)
				Expect(bucketName).To(Equal(actualBucketName))
				Expect(objectName).To(Equal(actualObjectName))
				Expect(reader).To(Equal(exampleContent))
				Expect(contentType).To(Equal("text/html"))
			})

			Context("when the make bucket fails", func() {
				BeforeEach(func() {
					fakeClient.MakeBucketReturns(errors.New("Arrr"))
				})

				It("should return a corresponding error", func() {
					err := testStore.Write(actualBucketName, actualObjectName, actualLoaction, exampleContent)
					Expect(err).To(MatchError("Arrr"))
					Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
					bucketName := fakeClient.BucketExistsArgsForCall(0)
					Expect(bucketName).To(Equal(actualBucketName))

					Expect(fakeClient.MakeBucketCallCount()).To(Equal(1))
					bucketName, location := fakeClient.MakeBucketArgsForCall(0)
					Expect(bucketName).To(Equal(actualBucketName))
					Expect(location).To(Equal(actualLoaction))

					Expect(fakeClient.PutObjectCallCount()).To(Equal(0))
				})
			})
		})

		Context("when the bucket exists", func() {
			BeforeEach(func() {
				fakeClient.BucketExistsReturns(true, nil)
			})

			It("should create an object in the bucket", func() {
				err := testStore.Write(actualBucketName, actualObjectName, actualLoaction, exampleContent)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
				bucketName := fakeClient.BucketExistsArgsForCall(0)
				Expect(bucketName).To(Equal(actualBucketName))

				Expect(fakeClient.MakeBucketCallCount()).To(Equal(0))

				Expect(fakeClient.PutObjectCallCount()).To(Equal(1))
				bucketName, objectName, reader, contentType := fakeClient.PutObjectArgsForCall(0)
				Expect(bucketName).To(Equal(actualBucketName))
				Expect(objectName).To(Equal(actualObjectName))
				Expect(reader).To(Equal(exampleContent))
				Expect(contentType).To(Equal("text/html"))
			})

			Context("when the create object returns an error", func() {
				BeforeEach(func() {
					fakeClient.PutObjectReturns(0, errors.New("Hey don't put!"))
				})

				It("should return the error", func() {
					err := testStore.Write(actualBucketName, actualObjectName, actualLoaction, exampleContent)
					Expect(err).To(MatchError("Hey don't put!"))

					Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
					bucketName := fakeClient.BucketExistsArgsForCall(0)
					Expect(bucketName).To(Equal(actualBucketName))

					Expect(fakeClient.MakeBucketCallCount()).To(Equal(0))

					Expect(fakeClient.PutObjectCallCount()).To(Equal(1))
				})
			})
		})
	})

	Describe("Read", func() {
		var (
			fakeClient                         *storefakes.FakeMinioClient
			actualBucketName, actualObjectName string
			actualContent                      []byte
			testStore                          *store.S3Store
		)

		BeforeEach(func() {
			fakeClient = new(storefakes.FakeMinioClient)
			actualBucketName = "read-bucket"
			actualObjectName = "read-thing"
			actualContent = []byte("Arrgggh!\nOooops")
			testStore = &store.S3Store{MinioClient: fakeClient}
		})

		It("should return the string", func() {
			fakeClient.FGetObjectStub = func(bucketName, objectName, filePath string) error {
				if bucketName == "read-bucket" && objectName == "read-thing" {
					err := ioutil.WriteFile(filePath, actualContent, 0600)
					Expect(err).ToNot(HaveOccurred())
				}
				return nil
			}

			content, err := testStore.Read(actualBucketName, actualObjectName)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(string(actualContent)))
		})

		It("should not leave temp files around", func() {
			files, err := ioutil.ReadDir("/tmp")
			Expect(err).ToNot(HaveOccurred())
			for _, file := range files {
				Expect(file.Name()).ShouldNot(ContainSubstring("tempS3ObjectFile"))
			}
		})

		Context("when the get fails", func() {
			It("should return the corresponding error", func() {
				fakeClient.FGetObjectReturns(errors.New("Yo-failed"))

				_, err := testStore.Read(actualBucketName, actualObjectName)
				Expect(err).To(MatchError("Yo-failed"))
			})

		})
	})
})
