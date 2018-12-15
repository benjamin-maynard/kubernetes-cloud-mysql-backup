# aws-database-backup

aws-database-backup is a container image based from Alpine Linux. This container is designed to run in Kubernetes as a cronjob to perform automatic backups of MySQL databases to Amazon S3.

This container was made to suit my own specific needs, and so is fairly limited in terms of configuration options. As of now, it performs a full database dump using the `mysqldump` command, and uploads it to an S3 Bucket specificed via environment variables. A full list of configuration environment variables are listed below.

Over time, this will likely be updated to support more features and functionality. You can read my blog post about my Kubernetes Architecture [here](https://benjamin.maynard.io/this-blog-now-runs-on-kubernetes-heres-the-architecture/).


[All changes are documented in the changelog](CHANGELOG.md)

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
| TARGET_DATABASE_PORT        | Port MySQL is listening on.                              |
| TARGET_DATABASE_NAME        | Name of the database to dump.                            |
| TARGET_DATABASE_USER        | Username to authenticate to the database with.           |
| TARGET_DATABASE_PASSWORD    | Password to authenticate to the database with. Should be configured using a Secret in Kubernetes. |

## Configuring the S3 Bucket & AWS IAM User

aws-database-backup performs a backup to the same path, with the same filename each time it runs (unless you change the environment variables each time). It therefore assumes that you have Versioning enabled on your S3 Bucket. A typical setup would involve S3 Versioning, with a Lifecycle Policy.

An IAM Users should be created, with API Credentials. An example Policy to attach to the IAM User (for a minimal permissions set) is as follows:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "s3:ListBucket",
            "Resource": "arn:aws:s3:::<BUCKET NAME>"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "s3:PutObject"
            ],
            "Resource": "arn:aws:s3:::<BUCKET NAME>/*"
        }
    ]
}
```

## Example Kubernetes Cronjob

An example of how to schedule this container in Kubernetes as a cronjob is below. This would configure a database backup to run each day at 01:00am. The AWS Secret Access Key, and Target Database Password are stored in secrets.

```
apiVersion: v1
kind: Secret
metadata:
  name: AWS_SECRET_ACCESS_KEY
type: Opaque
data:
  aws_secret_access_key: <AWS Secret Access Key>
--
apiVersion: v1
kind: Secret
metadata:
  name: TARGET_DATABASE_PASSWORD
type: Opaque
data:
  database_password: <Your Database Password>
--
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: my-database-backup
spec:
  schedule: "0 01 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: my-database-backup
            image: benjaminmaynard/aws-database-backup
            imagePullPolicy: Always
            env:
              - name: AWS_ACCESS_KEY_ID
                value: "<Your Access Key>"
              - name: AWS_SECRET_ACCESS_KEY
                valueFrom:
                   secretKeyRef:
                     name: AWS_SECRET_ACCESS_KEY
                     key: aws_secret_access_key
              - name: AWS_DEFAULT_REGION
                value: "<Your S3 Bucket Region>"
              - name: AWS_BUCKET_NAME
                value: "<Your S3 Bucket Name>"
              - name: AWS_BUCKET_BACKUP_PATH
                value: "<Your S3 Bucket Backup Path>"
              - name: AWS_BUCCKET_BACKUP_NAME
                value: "<Your Backup File Name.sql>"
              - name: TARGET_DATABASE_HOST
                value: "<Your Target Database Host>"
              - name: TARGET_DATABASE_PORT
                value: "<Your Target Database Port>"
              - name: TARGET_DATABASE_NAME
                value: "<Your Target Database Name>"
              - name: TARGET_DATABASE_USER
                value: "<Your Target Database Username>"
              - name: TARGET_DATABASE_PASSWORD
                valueFrom:
                   secretKeyRef:
                     name: TARGET_DATABASE_PASSWORD
                     key: database_password
          restartPolicy: OnFailure
```
