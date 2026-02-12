# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Kind/Name positional argument syntax (`Deployment/nginx`, `*/nginx`, `Service/*`)
- Wildcard support with `*` for matching any kind or name

### Changed
- Input is now stdin-only (removed file path argument)
- Flags (`--kind`, `--name`) override positional arguments when both provided

## [0.0.3] - 2025-09-26

## [0.0.2] - 2025-02-21

## [0.0.1] - 2023-06-29

- Initial release

[Unreleased]: https://github.com/rjferguson21/kgrep/compare/v0.0.3...HEAD
[0.0.3]: https://github.com/rjferguson21/kgrep/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/rjferguson21/kgrep/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/rjferguson21/kgrep/releases/tag/v0.0.1
