# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.8]

### Fixed
- `prometheusConfig.configOverride` only replaces the scrape_configs in the Prometheus configuration.

## [0.0.7]

### Added
- Arbitrary annotations can be added to the configmap and secret.

### Fixed
- Removed duplicate `resources` in deployment
- Moved `nodeSelector` to `server` block in values.yaml
- Updated scrape configs to properly drop unused metrics

## [0.0.6]

### Fixed
- Failure to scrape for container metrics due to improper joining of kubeMetrics and containerMetrics.

## [0.0.4]

### Added
- updated scrape configs to only scrape metrics that are needed
- added CPU and memory requests and limits to `deploy.yaml` and updated `values.yaml` with suggested values
- refactored `PVC.yaml` and updated `values.yaml` to match

### Fixed
- Updated `Exporting Pod Labels` README.md for clarity
- Fixed typos in templates
- removed extra `seperator` in cm.yaml

## [0.0.3]

### Fixed
- Fixed typos in readme

## [0.0.2]

### Fixed
- `existingSecretName` is used in the `secretName` named template so that the deployment uses the correct Secret
- `cloudAccountId` is coerced to a string to prevent malformed `cloud_account_id` query params

## [0.0.1]

### Added
- Initial release
