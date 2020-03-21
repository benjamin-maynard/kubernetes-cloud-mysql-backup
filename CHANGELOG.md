# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v2.3.0] - 21-03-2020
### Implement GCS Backend, and rename to kubernetes-cloud-mysql-backup
- Added the ability to GZIP compress backup files. Thanks & credit: @LucasBG0
- Added the ability to use custom S3 compatible storage endpoints. Thanks & credit: @mwienk
- Bumped Google Cloud SDK version to 285.0.1
- Bumped Alpine Linux version to 3.11
- Corrected log filenames that were not correctly updated as part of the v2.2.0 rename


## [v2.2.0] - 28-11-2019
### Implement GCS Backend, and rename to kubernetes-cloud-mysql-backup
- Added the ability to use Google Cloud Storage (GCS) as a backend storage provider (backwards compatible)
- Renamed to kubernetes-cloud-mysql-backup to better reflect the function of the application
- Improved environment variable processing (removed case sensitivity of Slack environment variable)
- Upgraded Alpine version to 3.10 base
- Switched to Python3
- Documentation improvements

## [v2.1.0] - 28-08-2019
### Added the ability to add a timestamp to the backup file name
- Ability to append timestamp to the database dump via the BACKUP_TIMESTAMP environment variable added. Thanks & credit: @kuzm1ch

## [v2.0.1] - 21-05-2019
### Renamed to kubernetes-mysql-backup
- Renamed to kubernetes-s3-mysql-backup from aws-database-backup to better describe function

## [v2.0.0] - 16-12-2018
### Fix issue with Slack Alerts
- Implemented the ability to backup multiple databases from a single host
- Updated the Variable Name of TARGET_DATABASE_NAME to TARGET_DATABASE_NAMES
- Updated the format of AWS_BUCKET_BACKUP_PATH so that the trailing / is not required
- Removed $AWS_BUCCKET_BACKUP_NAME variable (which had a typo). Database backups are now saved using their database names

## [v1.1.1] - 16-12-2018
### Fix issue with Slack Alerts
- Fixed issue with failed Slack alerts when log messages contained special characters
- Fixed /bin/ash error when evaluating if the log files are empty or not
- Fixed an error message about the log file not existing when the backup runs successfully
- Suppressed CURL output for Slack alerts

## [v1.1.0] - 15-12-2018
### Slack Integration & Error Handling
- Added Slack Integration
- Introduced Error Handling to make sure the container exits with the correct status, and provides useful debug information
- Fixed a bug where a failure of `mysqldump` would lead to a blank database backup being uploaded to S3
- Introduced default Environment Variables for non-essential values
- Improved the README.md

## [v1.0.0] - 02-12-2018
### Initial Release
- Initial Release