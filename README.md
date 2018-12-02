# aws-database-backup

aws-database backup is a container image based from Alpine Linux. This container is designed to run in Kubernetes as a cronjob to perform automatic backups of MySQL databases to Amazon S3.

This container was made to suit my own specific needs, and so is fairly limited in terms of configuration options. As of now, it performs a full database dump using the `mysqldump` command, and uploads it to an S3 Bucket specificed via environment variables. A full list of configuration environment variables are listed below.

Over time, this will likely be updated to support more features and functionality.

## Environment Variables

The following environment variables are required by aws-database-backup.

| Environment Variable        | Purpose                                   |
| --------------------------- |-------------------------------------------|
| AWS_ACCESS_KEY_ID           | AWS IAM Access Key ID.                                   |
| AWS_SECRET_ACCESS_KEY       | AWS IAM Secret Access Key. Should have very limited IAM permissions (see below for example) and should be configured using a Secret in Kubernetes.                                                            |
| AWS_DEFAULT_REGION          | Region of the S3 Bucket (e.g. eu-west-2).                |
| AWS_BUCKET_NAME             | The name of the S3 bucket.                               |
| AWS_BUCKET_BACKUP_PATH      | Path the backup file should be saved to in S3. E.g. `/database/myblog/backups/`. **Requires the trailing / and should not include the file name.**                                                             |
| AWS_BUCCKET_BACKUP_NAME     | File name of the backup file. E.g. `database_dump.sql`.  |
| TARGET_DATABASE_HOST        | Hostname or IP address of the MySQL Host.                |
| TARGET_DATABASE_NAME        | Name of the database to dump.                            |
| TARGET_DATABASE_USER        | Username to authenticate to the database with.           |
| TARGET_DATABASE_PASSWORD    | Password to authenticate to the database with.           |