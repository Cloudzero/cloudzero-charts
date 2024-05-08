# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

### Added
-  updated crape configs to only scrape metrics that are needed

## [0.0.2]

### Fixed
- `existingSecretName` is used in the `secretName` named template so that the deployment uses the correct Secret
- `cloudAccountId` is coerced to a string to prevent malformed `cloud_account_id` query params

## [0.0.1]

### Added
- Initial release
