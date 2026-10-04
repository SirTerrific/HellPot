# Changelog

This fork starts from upstream [yunginnanet/HellPot](https://github.com/yunginnanet/HellPot) at commit `0ba62c9`. No version has been tagged yet, so everything below is unreleased.

## Unreleased

### Compatibility

Configuration keys and defaults, the log structure (JSON keys, levels, messages, field order), the HTTP behavior, and the Docker image layout (`/app`, `/config`, `/logs`, port `8080`, entrypoint, user) are unchanged.

Two differences in log *text*, both at debug or trace level and both coming from the updated fasthttp library, not from HellPot: `END_ON_ERR` now reports `fasthttputil: connection closed` instead of `connection closed`, and the "non-GET request" line now reads `fasthttp: non-get request received`.

### Added

- `performance.max_conns_per_ip` (environment variable `HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP`): maximum simultaneous connections per remote IP. The default is `10`, the value that was hard-coded, and `0` means unlimited. Behind a reverse proxy all clients share the proxy's IP, so the old fixed limit capped the whole instance and refused the extra clients with HTTP 429.
- Docker image published to `ghcr.io/sirterrific/hellpot` (`linux/amd64`, `linux/arm64`) by a new workflow, rebuilt on every push to `main`, on `v*` tags and weekly.
- `govulncheck` step in CI.
- First unit test (`heffalump/heffalump_test.go`).
- `README.fr.md` (French translation), this changelog, and a rewritten README that documents every configuration key, the command line, the Docker image, the log format and the default-handling behavior of `-c`.

### Changed

- Go modules updated: fasthttp 1.63.0 to 1.74.0, zerolog 1.34.0 to 1.35.1, koanf/v2 2.2.1 to 2.3.7, koanf file provider 1.2.0 to 1.2.1, golang.org/x/term 0.32.0 to 0.46.0, golang.org/x/sys 0.33.0 to 0.48.0, plus their indirect dependencies. The `go` directive is now 1.26.0 (required by the updated dependencies) and the `toolchain` line was removed.
- Go module renamed to `github.com/SirTerrific/HellPot` (imports and README links updated). Use `go install github.com/SirTerrific/HellPot/cmd/HellPot@latest`.
- Dockerfile: build image `golang:1.27`, runtime image `distroless/static-debian13` (was `static-debian11`), native cross-compilation for multi-arch builds, optional `VERSION` build argument, `org.opencontainers.image.source` label points to this repository.
- CI: runs on the latest stable Go (Go 1.26.0 has known standard-library vulnerabilities reported by `govulncheck`), `actions/checkout` v5, `actions/setup-go` v6. Dependabot now also tracks Docker base images (weekly).
- Release binaries now embed the tag name (for example `v1.2.3`, shown as `1.2.3` by `--banner`). The workflow used to embed the full git ref.

### Fixed

- `robots.txt` handler passed generated content as a `Fprintf` format string (`go vet` error): a `%` in a configured path would have corrupted the output.

### Removed

- The unused 100 KB buffer pool in `heffalump`: `io.CopyBuffer` never used it because `*bufio.Writer` implements `io.ReaderFrom`. Output is unchanged.
- The `trigger.yml` workflow (it dispatched to the upstream repository) and the Docker Hub publishing job of the release workflow (it pushed to the upstream's Docker Hub and GHCR namespaces).
