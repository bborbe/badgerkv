# Changelog

All notable changes to this project will be documented in this file.

## v1.9.10

- Update go-git/go-git to v5.17.1 (fix security vulnerabilities)

## v1.9.9

- bump bborbe/errors, bborbe/kv, bborbe/run dependencies
- bump golangci-lint v2.11.4, osv-scanner v2.3.5
- bump docker, containerd, moby toolchain deps
- add runtime-spec replace directive for opencontainers/runtime-spec

## v1.9.8

- Update bborbe/collection to v1.20.7, bborbe/errors to v1.5.7, bborbe/kv to v1.19.2
- Update shoenig/go-modtool to v0.6.0
- Update bbolt to v1.4.3, go-yaml/v3 to v3.0.4
- Remove replace/exclude directives from go.mod

## v1.9.7

- chore: enable golangci-lint in Makefile check target and update .golangci.yml to standard config with nestif, errname, unparam, bodyclose, forcetypeassert, asasalint, prealloc linters
- refactor: extract runTx helper in badgerdb to eliminate dupl violation between Update and View
- fix: simplify bool comparisons and use bytes.Equal in badgerkv_tx.go to resolve staticcheck violations

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.9.6

- standardize Makefile: multiline trivy format

## v1.9.5

- chore: fix Go module cache corruption to restore passing precommit

## v1.9.4

- go mod update

## v1.9.3

- Update Go to 1.26.0

## v1.9.2

- Update Go to 1.25.7
- Update github.com/dgraph-io/badger/v4 to v4.9.1
- Update github.com/onsi/ginkgo/v2 to v2.28.1 and gomega to v1.39.1
- Update github.com/google/osv-scanner/v2 to v2.3.2
- Update numerous indirect dependencies

## v1.9.1

- Update Go to 1.25.5
- Update golang.org/x/crypto to v0.47.0
- Update dependencies

## v1.9.0

- update go and deps

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
