# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.8.4

- add golangci-lint configuration
- enhance CI with Trivy security scanning
- update Makefile with additional security tools (osv-scanner, gosec, trivy)
- update Go version to 1.25.2
- improve code formatting for long function signatures
- go mod update

## v1.8.3

- improve README with usage example and installation instructions
- go mod update

## v1.8.2

- add github workflow
- go mod update

## v1.8.1

- add tests
- go mod update

## v1.8.0

- OpenDB and OpenMemory return badgerkv.DB
- remove vendor
- go mod update

## v1.7.3

- go mod update

## v1.7.2

- go mod update

## v1.7.1

- fix ListBucketNames

## v1.7.0

- implement ListBucketNames
- go mod update

## v1.6.0

- add remove db files
- go mod update

## v1.5.2

- go mod update

## v1.5.1

- go mod update

## v1.5.0

- cache buckets per tx
- go mod update

## v1.4.2

- fix bucket name problem
- go mod update

## v1.4.1

- add interface to access bolt db, tx, bucket if needed

## v1.4.0

- prevent transaction open second transaction

## v1.3.0

- fulfill bucket testsuite

## v1.2.0

- use new testsuite

## v1.1.1

- update libkv

## v1.1.0

- Add context to update and view

## v1.0.0

- Initial Version
