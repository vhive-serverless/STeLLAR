package util

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const (
	awsRegion      = "us-west-1"
	S3Bucket       = "benchmarking-aws"
	localZipName = "benchmarking.zip"
	randomFileName = "random.file"
)

var S3ZipName string

//GenerateZIPLocation will create the serverless function zip deployment for the given provider,
//in the given language and of the given size in bytes.
func GenerateZIPLocation(provider string, language string, sizeBytes int) string {
	if fileExists(localZipName) {
		log.Infof("Local ZIP archive `%s` already exists, removing...", localZipName)
		if err := os.Remove(localZipName); err != nil {
			log.Fatalf("Failed to remove local ZIP archive `%s`", localZipName)
		}
	}

	log.Infof("Building %s handler...", language)
	codePath := fmt.Sprintf("code/producer/%s/%s-handler.go", language, provider)
	if !fileExists(codePath) {
		log.Fatalf("Code path `%s` does not exist, cannot deploy/update code.", codePath)
	}

	switch language {
	case "go1.x":
		RunCommandAndLog(exec.Command("go", "build", "-v", "-o", "producer-handler",
			"code/producer/go1.x/aws-handler.go"))
	//case "python3.8":
	//	RunCommandAndLog(exec.Command("go", "build", "-v", "-race", "-o", "producer-handler"))
	default:
		log.Fatalf("Unrecognized language %s", language)
	}

	generateRandomFile(sizeBytes)
	RunCommandAndLog(exec.Command("zip", localZipName, "producer-handler", randomFileName))

	if strings.Compare(provider, "aws") == 0 && sizeBytes > 50_000_000 {
		log.Infof(`Deploying to AWS and package size (~%dMB) > 50 MB, will now attempt to upload to Amazon S3.`, sizeBytes/1_000_000.0)
		S3ZipName = fmt.Sprintf("benchmarking%d.zip", sizeBytes)

		zipFile, err := os.Open(localZipName)
		if err != nil {
			log.Fatalf("Failed to open zip file %q: %v", localZipName, err)
		}

		sessionInstance := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(awsRegion),
		}))
		uploader := s3manager.NewUploader(sessionInstance)
		uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("benchmarking-aws"),
			Key:    aws.String(S3ZipName),
			Body:   zipFile,
		})
		if err != nil {
			log.Fatalf("Unable to upload %q to %q, %v", S3ZipName, S3Bucket, err.Error())
		}

		log.Infof("Successfully uploaded %q to bucket %q (%s)", S3ZipName, S3Bucket, uploadOutput.Location)
		return uploadOutput.Location
	}
	return fmt.Sprintf("fileb://%s", localZipName)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func generateRandomFile(sizeBytes int) {
	if fileExists(randomFileName) {
		log.Infof("Random file `%s` already exists, removing...", randomFileName)
		if err := os.Remove(randomFileName); err != nil {
			log.Fatalf("Failed to remove random file `%s`", randomFileName)
		}
	}

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("Failed to fill buffer with random bytes: `%s`", err.Error())
	}

	if err := ioutil.WriteFile(randomFileName, buffer, 0666); err != nil {
		log.Fatalf("Could not generate random file with size %d bytes", sizeBytes)
	}
}
