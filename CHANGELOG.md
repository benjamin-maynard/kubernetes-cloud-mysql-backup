# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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