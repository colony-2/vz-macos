# VM API validation

Fresh checks used Go 1.25.10 on Linux arm64 with GOWORK disabled. These results cover the configured-Backend revision and replace historical candidate passes. The review's `test-results.md` maps API-001–035 individually; `evidence/commands.json` records exact commands, environments, exit statuses and log paths.

| Check | Result |
|---|---|
| `go test -v ./vmapi` | Pass, including examples, error matching, full intent mismatch matrix, configuration and restart/domain fixtures |
| `go vet ./vmapi` | Pass |
| Fresh external consumer `go test -v ./...` | Pass with ordinary non-test vmapi import and copied external fixtures |
| Consumer `go list -deps -test -json ./...` | Only vmapi, consumer/test harness and standard-library packages; no Apple dependency |
| Package and consumer `go test -c`, CGO disabled, linux/arm64 | Pass; compiled, not cross-executed |
| Package and consumer `go test -c`, CGO disabled, darwin/arm64 | Pass; compiled, not executed or signed |
| `gofmt -l vmapi`, `git diff --check` | Pass after removing an extra blank line at EOF |
| `go build ./...` | Fails: configured Apple replacement directory is absent |
| `go test ./...` | Fails: configured Apple replacement directory is absent; some independent packages pass |

The isolated scratch module requires `github.com/tmc/vz-macos v0.0.0`, replaced by this checkout. It uses GOWORK=off and no Apple sources; `/Users/tmc/go/src/github.com/tmc/apple` is absent. Root go.mod is unchanged. Consumer source, dependency inventory and individual logs are retained with the review. Its local replacement proves source usability, not published module resolution.

Compatible Mac root execution and signing remain unavailable separately: the worker is Linux, the configured Apple sources are absent, and codesign is absent. No native executable was produced. Darwin compilation cannot establish virtualization entitlements or runtime behavior; later Mac builds must follow AGENTS.md signing requirements. No Docker or deployed VM checks are required for this source contract.

## Publication and consumer records

Fresh `git tag --list` and `git ls-remote --tags origin` both return no tags. The production-import search outside vmapi finds no imports (rg exit 1 means no matches). Supplied original audit and pulse-tart owner records describe future consumption, not an existing dependency on this candidate. The original audit job is recorded completed/success. `c2j cells` cannot enumerate more cells: `JobDB runtime URL is required`. Available records identify no released dependency; arbitrary external consumers and release systems beyond these records remain unknown.

This removes Placement and its fields from the unreleased candidate while preserving all 17 Backend method signatures. No compatibility alias, constructor, production adapter, CLI/wire change or go.mod change is included. A later discovery of a released dependency requires an explicit deprecation/migration decision before removing its API.

## Source and runtime evidence

The independent Tart-, Orchard- and VZ-shaped fixtures satisfy Backend. The consumer records a stable backend key and immutable configuration revision, rejects missing/rebound bindings before dispatch, preserves source/endpoint reservations, and validates the full normalized effective intent plus known UID. Persisted uncertain journals reopen with new objects and observe the same configured domain without Create replay. Two domains deliberately share names and UIDs; source bindings, cursor rejection, endpoint selection and originating recovery records remain distinct in the fixtures.

Declaration/documentation review establishes complete observational inventory, immutable UID mutation requirements, suspended-clone compatibility, bounded sensitive Exec, typed bootstrap evidence, exact endpoint ownership, eligible at-least-once recovery and explicit fallbacks. Fixtures demonstrate consumer decisions and representability. Their in-memory stores and JSON journal round trip do not establish native durability, secrecy, identity race safety, worker affinity, suspended cloning or endpoint connectivity. Those remain later backend qualification obligations.
