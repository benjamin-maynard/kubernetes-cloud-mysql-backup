package config

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sethvargo/go-envconfig"
)

// BackupDestinationConfig is used for storing the configuration for backup storage providers
type BackupDestinationConfig struct {
	// The name of the backup provider to use, either "gcp" or "aws"
	BackupProvider string `env:"BACKUP_PROVIDER,default=aws" validate:"required,oneof=aws gcp"`

	//
	// Google Cloud GCS Config
	//
	// A base64 encoded GCP Service Account Key. This should almost never be set. You should use Workload Identity
	// which allows the container to get short lived credentials from the Metadata Server which is much more
	// secure
	GCPServiceAccountKey string `env:"GCP_SERVICE_ACCOUNT_KEY" validate:"omitempty,base64"`
	// The name of the Google Cloud Storage Bucket where backups should be stored.
	//  You should not include the gs:// prefix. Only the bucket name
	GCPBucketName string `env:"GCS_BUCKET_NAME" validate:"required_if=BackupProvider gcp"`
	// The path to store the backup file in GCS. This should not contain a trailing slash or file name
	GCPBucketBackupPath string `env:"GCS_BUCKET_BACKUP_PATH" validate:"required_if=BackupProvider gcp"`

	//
	// AWS S3 Config
	//
	// The AWS Access Key ID which will be used for authenticating to the bucket
	AWSAccessKeyID string `env:"AWS_ACCESS_KEY_ID" validate:"required_if=BackupProvider aws"`
	// The AWS Secret Access Key which will be used for authenticating to the bucket
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"  validate:"required_if=BackupProvider aws"`
	// The region that the AWS S3 bucket is located
	AWSDefaultRegion string `env:"AWS_DEFAULT_REGION"  validate:"required_if=BackupProvider aws"`
	// The name of the S3 Bucket where the backup will be stored
	AWSBucketName string `env:"AWS_BUCKET_NAME"  validate:"required_if=BackupProvider aws"`
	// The path to store the backup file in s3. This should not contain a trailing slash or file name
	AWSBucketBackupPath string `env:"AWS_BUCKET_BACKUP_PATH"  validate:"required_if=BackupProvider aws"`
	// The S3 Endpoint to use if using a non AWS S3 compatible bucket
	AWSS3Endpoint string `env:"AWS_S3_ENDPOINT"`
}

// DatabaseBackupConfig is used for storing the configuration for the databases that need to be
// backed up
type DatabaseBackupConfig struct {
	// The hostname or IP addess of the MySQL host to backup
	TargetDatabaseHost string `env:"TARGET_DATABASE_HOST" validate:"required"`
	// The port that the MySQL host is listening on
	TargetDatabasePort int `env:"TARGET_DATABASE_PORT,default=3306" validate:"gte=0,lte=65535"`
	// The username to use for authenticating to the MySQL Host
	TargetDatabaseUser string `env:"TARGET_DATABASE_USER" validate:"required"`
	// The password for authenticating to the MySQL Host
	TargetDatabasePassword string `env:"TARGET_DATABASE_PASSWORD" validate:"required"`
	// The Name(s) of the databases to dump. Should be comma separated
	TargetDatabaseNames string `env:"TARGET_DATABASE_NAMES" validate:"required_if=TargetAllDatabases false"`
	// If set to true, all databases will be backed up
	TargetAllDatabases bool `env:"TARGET_ALL_DATABASES,default=false"`
	// If set to true, the "CREATE DATABASE" and "USE" statements will be added to the MySQL backup
	BackupCreateDatabaseStatement bool `env:"BACKUP_CREATE_DATABASE_STATEMENT,default=false"`
	// Custom additional flags to add to the mysql dump command
	BackupAdditionalParams string `env:"BACKUP_ADDITIONAL_PARAMS"`
	// If set to true, the backup will be compressed
	BackupCompress bool `env:"BACKUP_COMPRESS,default=false"` // Enable gzip backup compression
	// The GZIP compression level to use for the backup
	BackupCompressLevel int `env:"BACKUP_COMPRESS_LEVEL,default=9" validate:"gte=1,lte=9"`
	// The Public Key to use for the Age encryption
	// Public key used to encrypt the backup with age
	AgePublicKey string `env:"AGE_PUBLIC_KEY"`
	// Golang time formatting string used to prefix to the backup file name
	BackupTimestamp string `env:"BACKUP_TIMESTAMP"`
}

type SlackNotificationConfig struct {
	// If set to true, backups to slack are enabled
	SlackEnabled bool `env:"SLACK_ENABLED,default=false"`
	// The username to use for posting Slack messages
	SlackUsername string `env:"SLACK_USERNAME" validate:"required_if=SlackEnabled true"`
	// The Slack channel name to publish the notification to
	SlackChannel string `env:"SLACK_CHANNEL" validate:"required_if=SlackEnabled true"`
	// The Webhook URL where the Slack message should be POST'ed to
	SlackWebhookURL string `env:"SLACK_WEBHOOK_URL" validate:"required_if=SlackEnabled true,omitempty,url"`
	// The Proxy URL to use for the Slack network call
	SlackProxy string `env:"SLACK_PROXY"` // Proxy url to use for Slack network call
}

// The overall configuration for the application
type AppConfig struct {
	BackupDestinationConfig
	DatabaseBackupConfig
	SlackNotificationConfig
}

// NewConfigFromEnvironment returns a new application configuration, loading the config
// from environment variables, and validating it. If any configuration errors are detected
// an error is returned
func NewConfigFromEnvironment() (AppConfig, error) {

	// Create an empty AppConfig
	config := AppConfig{}

	// Load the environment variables
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		return AppConfig{}, err
	}

	// Run additional validation
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return AppConfig{}, err
	}

	return config, nil

}
