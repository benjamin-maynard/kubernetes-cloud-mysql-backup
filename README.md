# aws-database-backup
This container takes a backup of a specified database, and uploads it to a specified location in an AWS S3 bucket. It is configured through Environment Variables.

## Environment Variables (with examples)

- AWS_ACCESS_KEY_ID=
- AWS_SECRET_ACCESS_KEY=
- AWS_DEFAULT_REGION=
- AWS_BUCKET_NAME=
- AWS_BUCKET_BACKUP_PATH=
- AWS_BUCCKET_BACKUP_NAME=
- TARGET_DATABASE_HOST=
- TARGET_DATABASE_NAME=
- TARGET_DATABASE_USER=root
- TARGET_DATABASE_PASSWORD=