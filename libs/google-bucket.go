package libs

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func UploadFileToBucket(w io.Writer, bucket, object, localFilepath string) error {
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

func RemoveFileFromBucket(w io.Writer, bucket, object string) error {
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
