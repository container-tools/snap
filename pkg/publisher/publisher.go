package publisher

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	minio "github.com/minio/minio-go"
)

type Publisher struct{}

func (p *Publisher) Publish(dir, dest string) error {
	endpoint := envOr("MINIO_ENDPOINT", "minio-service-snap.192.168.42.192.nip.io")
	accessKeyID := envOr("MINIO_ACCESS_KEY_ID", "minio")
	secretAccessKey := envOr("MINIO_ACCESS_KEY_SECRET", "minio123")
	useSSL := false

	log.Printf("Using Minio endpoint %s", endpoint)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	parts := strings.FieldsFunc(dest, func(c rune) bool { return c == '/' })
	if len(parts) < 1 {
		return errors.New("invalid destination: " + dest)
	}

	bucket := parts[0]
	if err := p.getOrCreateBucket(minioClient, bucket); err != nil {
		return err
	}

	bucketPath := ""
	if len(parts) > 1 {
		bucketPath = strings.Join(parts[1:], "/")
	}

	return filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var relative string
		relative, err = filepath.Rel(dir, filePath)
		if err != nil {
			return err
		}
		relative = path.Join(bucketPath, relative)
		if err := p.publish(minioClient, bucket, filePath, relative); err != nil {
			return err
		}
		return nil
	})
}

func (p *Publisher) publish(client *minio.Client, bucket, sourceFile, destinationFile string) error {
	_, err := client.FPutObject(bucket, destinationFile, sourceFile, minio.PutObjectOptions{})
	return err
}

func (p *Publisher) getOrCreateBucket(client *minio.Client, bucket string) error {
	if exists, err := client.BucketExists(bucket); err != nil {
		return err
	} else if !exists {
		if err := client.MakeBucket(bucket, ""); err != nil {
			return err
		}
	}
	return nil
}

func envOr(env, def string) string {
	res := os.Getenv(env)
	if res != "" {
		return res
	}
	return def
}
