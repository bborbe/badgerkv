---
status: completed
summary: Enabled golangci-lint in Makefile check target, updated .golangci.yml to standard config with new linters, and fixed all lint violations by refactoring Update/View into a shared runTx helper and simplifying bool comparisons in badgerkv_tx.go
container: badgerkv-002-enable-lint
dark-factory-version: v0.59.5-dirty
created: "2026-03-21T10:00:00Z"
queued: "2026-03-21T10:12:34Z"
started: "2026-03-21T10:12:41Z"
completed: "2026-03-21T10:18:08Z"
---

<summary>
- The Makefile check target includes lint alongside vet, errcheck, and other checks
- The TODO comment about enabling lint is removed from the Makefile
- The golangci-lint configuration matches the standard config used across bborbe repos
- All golangci-lint violations are fixed in the source code
- The full precommit pipeline passes including the newly enabled lint step
</summary>

<objective>
Enable the golangci-lint linter in the Makefile check target for the badgerkv project. The lint target already exists but is excluded from the check dependency chain. Activate it, update the golangci-lint config to the current standard, and fix all resulting lint violations so make precommit passes.
</objective>

<context>
Read CLAUDE.md for project conventions and build commands.

Read these files before making changes:
- `Makefile` — current check target with lint commented out
- `.golangci.yml` — current golangci-lint config (outdated, missing linters)

Key observations:
- The Makefile has a `lint` target that already runs `golangci-lint v2` with `--config .golangci.yml`
- The `check` target currently lists: `vet errcheck vulncheck osv-scanner gosec trivy` (no lint)
- The TODO comment and commented-out check target with lint exist above the current check target
- The `.golangci.yml` is missing several linters and exclusion rules vs the standard

Reference standard `.golangci.yml` config to match:
```yaml
version: "2"

run:
  timeout: 5m
  tests: true

linters:
  enable:
    - govet
    - errcheck
    - staticcheck
    - unused
    - revive
    - gosec
    - gocyclo
    - depguard
    - dupl
    - nestif
    - errname
    - unparam
    - bodyclose
    - forcetypeassert
    - asasalint
    - prealloc
  settings:
    depguard:
      rules:
        Main:
          deny:
            - pkg: "github.com/pkg/errors"
              desc: "use github.com/bborbe/errors instead"
            - pkg: "github.com/bborbe/argument"
              desc: "use github.com/bborbe/argument/v2 instead"
            - pkg: "golang.org/x/net/context"
              desc: "use context from standard library instead"
            - pkg: "golang.org/x/lint/golint"
              desc: "deprecated, use revive or staticcheck instead"
            - pkg: "io/ioutil"
              desc: "deprecated since Go 1.16, use io and os packages instead"
    funlen:
      lines: 80
      statements: 50
    gocognit:
      min-complexity: 20
    nestif:
      min-complexity: 4
    maintidx:
      min-maintainability-index: 20
  exclusions:
    presets:
      - comments
      - std-error-handling
      - common-false-positives
    rules:
      - linters:
          - staticcheck
        text: "SA1019"
      - linters:
          - errname
        text: "(KeyNotFoundError|TransactionAlreadyOpenError|BucketNotFoundError|BucketAlreadyExistsError)"
      - linters:
          - revive
        path: "_test\\.go$"
        text: "dot-imports"
      - linters:
          - revive
        text: "unused-parameter"
      - linters:
          - revive
        text: "exported"
      - linters:
          - dupl
        path: "_test\\.go$"
      - linters:
          - unparam
        path: "_test\\.go$"
      - linters:
          - dupl
        path: "-test-suite\\.go$"
      - linters:
          - revive
        path: "-test-suite\\.go$"
        text: "dot-imports"

formatters:
  enable:
    - gofmt
    - goimports
```
</context>

<requirements>
1. Update `.golangci.yml` to match the standard config from `the reference standard config shown in the context section`:
   - Add missing linters under `linters.enable`: `nestif`, `errname`, `unparam`, `bodyclose`, `forcetypeassert`, `asasalint`, `prealloc`
   - Add missing `linters.settings` blocks: `funlen`, `gocognit`, `nestif`, `maintidx`
   - Fix the typo in the existing depguard deny rules: `"github.com/pkg/erros"` should be `"github.com/pkg/errors"` and `"github.com/bborbe/erros"` should be `"github.com/bborbe/errors"`
   - Add the additional depguard deny entries from the kv config (argument/v2, golang.org/x/net/context, golint, io/ioutil)
   - Add missing exclusion rules: `errname` text exclusion for badgerkv-specific error types (KeyNotFoundError, TransactionAlreadyOpenError, BucketNotFoundError, BucketAlreadyExistsError), `unparam` path exclusion for test files, `dupl` path exclusion for `-test-suite.go$` files, `revive` path exclusion for `-test-suite.go$` files with dot-imports text
   - Preserve any badgerkv-specific settings that are correct

2. Update `Makefile`:
   - Remove the two comment lines (the `# TODO: enable lint` line and the `# check: lint vet errcheck vulncheck osv-scanner gosec trivy` line)
   - Change the `check:` dependency line from `check: vet errcheck vulncheck osv-scanner gosec trivy` to `check: lint vet errcheck vulncheck osv-scanner gosec trivy`

3. Run `make lint`  to identify all lint violations

4. Fix all lint violations in the Go source files. Common fixes include:
   - Adding error checks for unchecked returns
   - Fixing type assertion safety (use comma-ok pattern instead of bare assertions)
   - Reducing function complexity or nesting depth
   - Removing unused parameters
   - Preallocating slices where appropriate
   - Closing response bodies
   - Do NOT suppress violations by adding nolint directives unless the violation is a clear false positive

5. Run `make lint` again to confirm zero violations
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT change the behavior of any existing functions — only fix lint violations
- Do NOT add nolint directives unless the lint finding is a clear false positive that cannot be fixed in code
- Do NOT modify test expectations or remove test coverage
- Existing tests must still pass
- Keep changes minimal — fix only what the linter flags
</constraints>

<verification>
Run `make precommit`  — must pass with exit code 0.
Run `make lint`  — must pass with zero violations.
</verification>
