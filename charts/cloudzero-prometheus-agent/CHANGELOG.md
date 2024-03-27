# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

## [0.0.3] - 2024-03-27
### Updated
- Metrics exported to CloudZero are configurable by an array in the values.yaml
### Added
- Adds the `node_cpu_seconds_total` metric by default

## [0.0.2] - 2024-03-26
### Changed
- Updated cloudzero-agent installation instructions to include additional options for secret management
## Fixed
- Updated cloudzero-agent README dependency section to remove outdated subcharts

## [0.0.1] - 2024-03-22
### Added
- Add initial `cloudzero-prometheus-agent` chart implementation