package ops

import (
	"context"

	"cloud.google.com/go/storage"
)

func CreateBucket() (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket("data.openland.com"), nil
}
