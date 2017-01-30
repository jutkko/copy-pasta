package store_test

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/jutkko/copy-pasta/store"
	"github.com/jutkko/copy-pasta/store/storefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("S3", func() {
	Describe("Write", func() {
		var fakeClient *storefakes.FakeMinioClient
		var exampleContent []string
		var bucketName, location string

		BeforeEach(func() {
			exampleContent = []string{"He is a banana", "and an apple"}
			fakeClient = new(storefakes.FakeMinioClient)
			bucketName = "this-bucket"
			location = "that-location"
		})

		Context("when the bucket exists command returns an error", func() {
			BeforeEach(func() {
				fakeClient.BucketExistsReturns(true, errors.New("No action should be taken"))
			})

			It("should return the error", func() {
				err := store.S3Write(fakeClient, bucketName, location, exampleContent)
				Expect(err).To(MatchError("No action should be taken"))
			})
		})

		Context("when the bucket doesn't exist", func() {
			BeforeEach(func() {
				fakeClient.BucketExistsReturns(false, nil)
			})

			It("should create it and put the object there", func() {
				err := store.S3Write(fakeClient, bucketName, location, exampleContent)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
				bucketName := fakeClient.BucketExistsArgsForCall(0)
				Expect(bucketName).To(Equal(bucketName))

				Expect(fakeClient.MakeBucketCallCount()).To(Equal(1))
				bucketName, location := fakeClient.MakeBucketArgsForCall(0)
				Expect(bucketName).To(Equal(bucketName))
				Expect(location).To(Equal(location))

				Expect(fakeClient.PutObjectCallCount()).To(Equal(1))
				bucketName, _, reader, contentType := fakeClient.PutObjectArgsForCall(0)
				Expect(bucketName).To(Equal(bucketName))
				Expect(reader).To(Equal(bytes.NewReader([]byte("He is a banana"))))
				Expect(contentType).To(Equal("text/html"))
			})

			Context("when the make bucket fails", func() {
				BeforeEach(func() {
					fakeClient.MakeBucketReturns(errors.New("Arrr"))
				})

				It("should return a corresponding error", func() {
					err := store.S3Write(fakeClient, bucketName, location, exampleContent)
					Expect(err).To(MatchError("Arrr"))
					Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
					bucketName := fakeClient.BucketExistsArgsForCall(0)
					Expect(bucketName).To(Equal(bucketName))

					Expect(fakeClient.MakeBucketCallCount()).To(Equal(1))
					bucketName, location := fakeClient.MakeBucketArgsForCall(0)
					Expect(bucketName).To(Equal(bucketName))
					Expect(location).To(Equal(location))

					Expect(fakeClient.PutObjectCallCount()).To(Equal(0))
				})
			})
		})

		Context("when the bucket exists", func() {
			BeforeEach(func() {
				fakeClient.BucketExistsReturns(true, nil)
			})

			It("should create an object in the bucket", func() {
				err := store.S3Write(fakeClient, bucketName, location, exampleContent)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
				bucketName := fakeClient.BucketExistsArgsForCall(0)
				Expect(bucketName).To(Equal(bucketName))

				Expect(fakeClient.MakeBucketCallCount()).To(Equal(0))

				Expect(fakeClient.PutObjectCallCount()).To(Equal(1))
				bucketName, _, reader, contentType := fakeClient.PutObjectArgsForCall(0)
				Expect(bucketName).To(Equal(bucketName))
				Expect(reader).To(Equal(bytes.NewReader([]byte("He is a banana"))))
				Expect(contentType).To(Equal("text/html"))
			})

			Context("when the create object returns an error", func() {
				BeforeEach(func() {
					fakeClient.PutObjectReturns(0, errors.New("Hey don't put!"))
				})

				It("should return the error", func() {
					err := store.S3Write(fakeClient, bucketName, location, exampleContent)
					Expect(err).To(MatchError("Hey don't put!"))

					Expect(fakeClient.BucketExistsCallCount()).To(Equal(1))
					bucketName := fakeClient.BucketExistsArgsForCall(0)
					Expect(bucketName).To(Equal(bucketName))

					Expect(fakeClient.MakeBucketCallCount()).To(Equal(0))

					Expect(fakeClient.PutObjectCallCount()).To(Equal(1))
				})
			})
		})
	})

	Describe("Read", func() {
		var fakeClient *storefakes.FakeMinioClient
		var bucketName, objectName string

		BeforeEach(func() {
			fakeClient = new(storefakes.FakeMinioClient)
			bucketName = "read-bucket"
			objectName = "read-thing"
		})

		It("should return the string", func() {
			fakeClient.FGetObjectStub = func(bucketName, objectName, filePath string) error {
				if bucketName == "read-bucket" && objectName == "read-thing" {
					err := ioutil.WriteFile(filePath, []byte("Arrgggh!"), 0600)
					Expect(err).ToNot(HaveOccurred())
				}
				return nil
			}

			content, err := store.S3Read(fakeClient, bucketName, objectName)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal([]string{"Arrgggh!"}))
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

				_, err := store.S3Read(fakeClient, bucketName, objectName)
				Expect(err).To(MatchError("Yo-failed"))
			})

		})
	})
})
