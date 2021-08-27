package putters

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pilillo/igovium/utils"
)

var s3Client *minio.Client
var s3Bucket string

func s3put(tmpPath string, partName string, tmpFile string, config *utils.S3Config) {
	log.Println("Uploading to remote s3 volume")

	if s3Client == nil {
		endpoint := config.Endpoint
		accessKeyID := os.Getenv(config.AccessKeyVarName)
		secretKey := os.Getenv(config.SecretKeyVarName)
		useSSL := config.UseSSL
		s3Bucket = config.Bucket

		var err error
		s3Client, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Println(err.Error())
			return
			//panic(err)
		}
	}

	// Upload the file to s3
	tmpFilePath := path.Join(tmpPath, partName, tmpFile)
	fi, _ := os.Stat(tmpFilePath)
	fileReader, _ := os.Open(tmpFilePath)
	// write asset definition to bucket
	defer fileReader.Close()
	// use partName as prefix for the target filename
	objectName := path.Join(partName, tmpFile)
	_, err := s3Client.PutObject(context.Background(),
		config.Bucket, objectName, fileReader, fi.Size(),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)

	if err != nil {
		log.Println(err.Error())
		return
	}

	// todo : delete file?
}
