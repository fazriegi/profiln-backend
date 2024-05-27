package libs

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type IGoogleBucket interface {
	HandleObjectUpload(imageFile *multipart.FileHeader, newObjectPath string) (string, error)
	HandleObjectDeletion(objectUrl ...string) error
	// HandleBatchObjectDeletion(objectUrls ...string) error
}

type GoogleBucket struct {
	fs         IFileSystem
	log        *logrus.Logger
	bucketName string
}

func NewGoogleBucket(fs IFileSystem, log *logrus.Logger) IGoogleBucket {
	return &GoogleBucket{
		bucketName: os.Getenv("BUCKET_NAME"),
		fs:         fs,
		log:        log,
	}
}

func (g *GoogleBucket) HandleObjectDeletion(objectUrls ...string) error {
	if len(objectUrls) < 1 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(objectUrls))

	for _, objectUrl := range objectUrls {
		// Extract the previous object path from the current document URL
		objectPath, err := extractBucketObjectUrl(objectUrl)
		if err != nil {
			return fmt.Errorf("extractBucketObjectUrl: %w", err)
		}

		wg.Add(1)
		go func(objectPath string) {
			defer wg.Done()
			if err := removeBucketObject(g.bucketName, objectPath); err != nil {
				errChan <- fmt.Errorf("removeBucketObject (%s): %v", objectPath, err)
			}
		}(objectPath)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// func (g *GoogleBucket) HandleObjectDeletion(objectUrl string) error {
// 	if objectUrl == "" {
// 		return nil
// 	}

// 	// Extract the previous object path from the current document URL
// 	objectPath, err := extractBucketObjectUrl(objectUrl)
// 	if err != nil {
// 		return fmt.Errorf("extractBucketObjectUrl: %w", err)
// 	}

// 	if err := removeBucketObject(g.bucketName, objectPath); err != nil {
// 		return fmt.Errorf("removeBucketObject (%s): %v", objectPath, err)
// 	}

// 	return nil
// }

func (g *GoogleBucket) HandleObjectUpload(imageFile *multipart.FileHeader, newObjectPath string) (string, error) {
	// Generate a new filename and save the file locally
	newFilename := g.fs.GenerateNewFilename(imageFile.Filename)
	fileDest := fmt.Sprintf("./storage/temp/file/%s", newFilename)

	// Save file to local temporary storage
	if err := g.fs.SaveFile(imageFile, fileDest); err != nil {
		return "", fmt.Errorf("fileSystem.SaveFile: %w", err)
	}

	// Upload the new file to the bucket
	bucketObject := fmt.Sprintf("%s/%s", newObjectPath, newFilename)
	if err := uploadBucketObject(g.bucketName, bucketObject, fileDest); err != nil {
		return "", fmt.Errorf("uploadBucketObject: %w", err)
	}

	// If file exists, delete it from the local temporary storage asynchronously
	if fileDest != "" {
		go func() {
			if err := g.fs.RemoveFile(fileDest); err != nil {
				g.log.Errorf("fileSystem.RemoveFile: %v", err)
			}
		}()
	}

	// Construct the new object URL
	objectUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, bucketObject)
	return objectUrl, nil
}

func uploadBucketObject(bucket, object, localFilepath string) error {
	credFilepath := os.Getenv("GOOGLE_APP_CREDENTIALS_FILEPATH")
	ctx := context.Background()

	// Initialize storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credFilepath))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Open local file
	f, err := os.Open(localFilepath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	obj := client.Bucket(bucket).Object(object)

	// Upload file to the bucket
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	return nil
}

func removeBucketObject(bucket, object string) error {
	credFilepath := os.Getenv("GOOGLE_APP_CREDENTIALS_FILEPATH")
	ctx := context.Background()

	// Initialize storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credFilepath))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Delete object from bucket
	obj := client.Bucket(bucket).Object(object)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%s).Delete: %w", object, err)
	}

	return nil
}

// Get the filepath from object url
func extractBucketObjectUrl(objectUrl string) (string, error) {
	parsedURL, err := url.Parse(objectUrl)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}

	// Split the path and get the relevant part
	parts := strings.SplitN(parsedURL.Path, "/", 3)

	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected URL format")
	}

	return parts[2], nil
}
