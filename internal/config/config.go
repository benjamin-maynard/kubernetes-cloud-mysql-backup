package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

// NewConfigFromEnvironment returns a new configuration from the environment
func NewConfigFromEnvironment(ctx context.Context) (*AppConfig, error) {

	var c AppConfig
	err := envconfig.Process(ctx, &c)

	if err != nil {
		return &c, err
	}

	err = c.Validate()
	if err != nil {
		return &c, err
	}

	return &c, nil

}

// Validate ensures that a configuration is valid
func (c *AppConfig) Validate() error {

	// Reflect the struct and determine the type
	v := reflect.ValueOf(*c)
	typeOfS := v.Type()

	// Loop through all of the fields in the struct
	for i := 0; i < v.NumField(); i++ {

		// Check if we have the requiredOnParentValue set
		val := typeOfS.Field(i).Tag.Get("requiredOnParentValue")

		fmt.Println(val)

		// If the value is not empty, then get the value of the parent environment variable
		if val != "" {

			// Determine the environment variable we are looking for
			envVar := strings.Split(val, "=")[0]      // The name of the environment variable we are looking for in the "env" struct tag
			validateVal := strings.Split(val, "=")[1] // The value of that struct value that should trigger validation

			// Find the first matching field with the struct tag "env" set to the desired environment variable
			fieldVal, err := getFieldValueByStructTagValue(*c, "env", envVar)
			if err != nil {
				return err
			}

			// Check if we should validate this field is not empty
			if fieldVal != validateVal {

				if fmt.Sprintf("%v", v.Field(i).Interface()) == "" {
					return fmt.Errorf("environment variable '%s' was expected because '%s' was set to '%s' but no value was found", typeOfS.Field(i).Tag.Get("env"), envVar, validateVal)
				}

			}

		}

	}

	return nil

}

// getFieldValueByStructTagValue returns the value of a struct field that has the specified struct tag and value
// it returns the first match only
func getFieldValueByStructTagValue(c interface{}, tag, value string) (string, error) {

	v := reflect.ValueOf(c)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {

		if typeOfS.Field(i).Tag.Get(tag) == value {
			return fmt.Sprintf("%v", v.Field(i).Interface()), nil
		}

	}

	return "", fmt.Errorf("found no matching field with struct tag '%s' and value '%s'", tag, value)

}

// The configuration of kubernetes-cloud-mysql-backup
type AppConfig struct {
	BackupProvider                string `env:"BACKUP_PROVIDER" default:"aws"`                               // The name of the backup provider to use, either "gcp" or "aws"
	GCPServiceAccountKey          string `env:"GCP_SERVICE_ACCOUNT_KEY"`                                     // Base64 encoded GCP Service Account Key. This should only be set if you are running outside of GCP (No Metadata Server)
	GCSBucketName                 string `env:"GCS_BUCKET_NAME"`                                             // The name of the GCS Bucket to store the backups in
	GCSBucketBackupPath           string `env:"GCS_BUCKET_BACKUP_PATH"`                                      // The path to store the backup file in GCS. Should not contain a trailing slash or file name
	AWSAccessKeyID                string `env:"AWS_ACCESS_KEY_ID"`                                           // AWS IAM Access Key ID for accessing the S3 Bucket
	AWSSecretAccessKey            string `env:"AWS_SECRET_ACCESS_KEY" `                                      // AWS IAM Secret Access Key for accessing the S3 Bucket
	AWSDefaultRegion              string `env:"AWS_DEFAULT_REGION""`                                         // Region of the S3 Bucket
	AWSBucketName                 string `env:"AWS_BUCKET_NAME" requiredOnParentValue:"BACKUP_PROVIDER=aws"` // Name of the S3 Bucket
	AWSBucketBackupPath           string `env:"AWS_BUCKET_BACKUP_PATH"`                                      // The path to store the backup file in S3. Should not contain a trailing slash or file name
	AWSS3Endpoint                 string `env:"AWS_S3_ENDPOINT"`                                             // The S3 Endpoint if using a non AWS S3 compatible bucket
	TargetDatabaseHost            string `env:"TARGET_DATABASE_HOST,required"`                               // Hostname or IP address of the MySQL Host
	TargetDatabasePort            string `env:"TARGET_DATABASE_PORT,default=3306"`                           // Port MySQL is listening on (Default: 3306)
	TargetDatabaseUser            string `env:"TARGET_DATABASE_USER,required"`                               // Username to use for authenticating to the MySQL Host
	TargetDatabasePassword        string `env:"TARGET_DATABASE_PASSWORD,required"`                           // Password for authenticating to the MySQL Host
	TargetDatabaseNames           string `env:"TARGET_DATABASE_NAMES"`                                       // Name(s) of the databases to dump. Should be comma separated
	TargetAllDatabases            bool   `env:"TARGET_ALL_DATABASES,default=false"`                          // Dump all databases
	BackupCreateDatabaseStatement bool   `env:"BACKUP_CREATE_DATABASE_STATEMENT,default=false"`              // Add the "CREATE DATABASE" and "USE" statements to the MySQL backup
	BackupAdditionalParams        string `env:"BACKUP_ADDITIONAL_PARAMS"`                                    // Custom additional parameters to add to the mysqldump command
	BackupCompress                bool   `env:"BACKUP_COMPRESS,default=false"`                               // Enable gzip backup compression
	BackupCompressLevel           int    `env:"BACKUP_COMPRESS_LEVEL,default=9"`                             // gzip compression level to use for the backup
	AgePublicKey                  string `env:"AGE_PUBLIC_KEY"`                                              // Public key used to encrypt the backup with age
	SlackEnabled                  bool   `env:"SLACK_ENABLED,default=false"`                                 // Enable the Slack integration
	SlackUsername                 string `env:"SLACK_USERNAME"`                                              // Username to use for the slack integration
	SlackChannel                  string `env:"SLACK_CHANNEL"`                                               // Slack channel name
	SlackWebhookURL               string `env:"SLACK_WEBHOOK_URL"`                                           // Webhook URL to publish Slack messages
	SlackProxy                    string `env:"SLACK_PROXY"`                                                 // Proxy url to use for Slack network call
}

// type BackupDestination interface {
// 	Upload(filePath string) error
// }
