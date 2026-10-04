<div align="center">
  <img src="https://tcp.ac/i/00ctL.gif" alt="HellPot"/>

[![Vibe Check](https://github.com/SirTerrific/HellPot/actions/workflows/go.yml/badge.svg)](https://github.com/SirTerrific/HellPot/actions/workflows/go.yml) [![Docker (GHCR)](https://github.com/SirTerrific/HellPot/actions/workflows/docker.yml/badge.svg)](https://github.com/SirTerrific/HellPot/actions/workflows/docker.yml) [![GoDoc](https://godoc.org/github.com/SirTerrific/HellPot?status.svg)](https://godoc.org/github.com/SirTerrific/HellPot) [![Go Report Card](https://goreportcard.com/badge/github.com/SirTerrific/HellPot)](https://goreportcard.com/report/github.com/SirTerrific/HellPot) [![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**English** | [Français](README.fr.md)

</div>

## About this fork

This is a maintained fork of [yunginnanet/HellPot](https://github.com/yunginnanet/HellPot), whose upstream is no longer maintained.

The goal is a drop-in replacement: **same behavior, same configuration file, same log format.** What changed is limited to updated dependencies, a refreshed build, a published Docker image, a few small fixes and one new optional setting. See [CHANGELOG.md](CHANGELOG.md) for the exact list.

If you were building the upstream image locally, switch to the published image without touching your configuration or log pipeline: see [Docker Method](#docker-method).

## Summary

HellPot is an endless honeypot based on [Heffalump](https://github.com/carlmjohnson/heffalump) that sends unruly HTTP bots to hell.

Notably it implements a [toml configuration file](https://github.com/knadh/koanf), has [JSON logging](https://github.com/rs/zerolog), and comes with significant performance gains.

## Grave Consequences

Clients (hopefully bots) that disregard `robots.txt` and connect to your instance of HellPot will **suffer eternal consequences**.

HellPot will send an infinite stream of data that is _just close enough_ to being a real website that they might just stick around until their soul is ripped apart and they cease to exist.

Under the hood of this eternal suffering is a markov engine that chucks bits and pieces of [The Birth of Tragedy (Hellenism and Pessimism)](https://www.gutenberg.org/files/51356/51356-h/51356-h.htm) by Friedrich Nietzsche at the client using [fasthttp](https://github.com/valyala/fasthttp).

## Installation

| Method | Command / location |
| --- | --- |
| Docker (multi-arch `linux/amd64` + `linux/arm64`) | `docker pull ghcr.io/sirterrific/hellpot:latest` |
| Compiled binaries (Linux, Windows, macOS, FreeBSD) | [GitHub Releases](https://github.com/SirTerrific/HellPot/releases/latest) |
| From source | see [Building From Source](#building-from-source) |

## Building From Source

HellPot requires **Go 1.26 or newer** (see the `go` line in [go.mod](go.mod)). It uses [go modules](https://go.dev/blog/using-go-modules), which makes it dead simple to build with a stock Go installation. A GNU Makefile is provided.

1 ) `git clone https://github.com/SirTerrific/HellPot`

2 ) `cd HellPot`

3 ) `make`

4 ) _Consider the potential grave consequences of your actions._

The Makefile targets are `deps` (`go mod tidy`), `check` (`go vet`), `build`, `run` and `format`; `make` alone runs `deps check build`. `make build` writes the `HellPot` binary in the current directory.

To install straight from the module path:

```
go install github.com/SirTerrific/HellPot/cmd/HellPot@latest
```

## Usage

### YOLO Method:

In the event of a missing configuration file, HellPot will attempt to place its default config in **$HOME/.config/HellPot/config.toml** (the user config directory of your OS). This allows irresponsible souls to begin raining hellfire with ease, **_immediately_**:

1 ) Download a [compiled release](https://github.com/SirTerrific/HellPot/releases/latest)

2 ) Run binary and immediately begin sending clients directly to hell.

---

### Reasonable Method:

1 ) Configure webserver as reverse proxy (see below)

2 ) `./HellPot --genconfig`

3 ) Edit your newly generated `config.toml` as desired.

4 ) Ponder your ~~existence~~ server's ability to handle your chosen performance values.

5 ) `./HellPot -c config.toml`

---

### Docker Method:

A multi-arch image is published to the GitHub Container Registry. It is rebuilt on every push to `main`, on every `v*` tag and every Monday (to pick up Go and base-image security patches).

```
docker pull ghcr.io/sirterrific/hellpot:latest
```

| Tag | Published when |
| --- | --- |
| `latest` | push to `main`, and the weekly rebuild |
| `sha-<short commit>` | every push to `main` or to a `v*` tag, and the weekly rebuild |
| `<version>` and `<major>.<minor>` | a `v*` tag is pushed (for example `v1.2.3` gives `1.2.3` and `1.2`) |

The image layout is the same as the upstream `Dockerfile`:

| Item | Value |
| --- | --- |
| Binary | `/app` |
| Entrypoint | `/app -c /config` |
| Default config | `/config` (a copy of [docker_config.toml](docker_config.toml)) |
| Log directory | `/logs/` |
| Port | `8080` |
| Base image | `gcr.io/distroless/static-debian13` (no shell) |
| User | root (uid 0) |

The default config binds `0.0.0.0:8080`, enables `catchall` (every GET path is answered with the tarpit, no `robots.txt` is served), sets the `Server` header to `nginx`, blacklists the user agents containing `Cloudflare-Traffic-Manager` or `curl`, reads the client IP from `X-Real-IP`, and writes its JSON logs to `/logs/` with color-free console output.

Quick start:

```
docker run -d --name hellpot -p 8080:8080 -v ./logs:/logs ghcr.io/sirterrific/hellpot:latest
```

Use your own configuration by mounting a file over `/config`:

```
docker run -d --name hellpot -p 8080:8080 -v ./config.toml:/config:ro -v ./logs:/logs ghcr.io/sirterrific/hellpot:latest
```

Or override single keys with [environment variables](#environment-variables). `docker compose` example, behind a reverse proxy that shares the external `proxy` network:

```yaml
services:
  hellpot:
    image: ghcr.io/sirterrific/hellpot:latest
    container_name: hellpot
    restart: unless-stopped
    expose:
      - 8080
    volumes:
      - ./logs:/logs
    environment:
      - TZ=America/Toronto
      # all clients share the reverse proxy's IP: raise the per-IP connection limit
      - HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP=100
    # optional hardening, verified to work with this image
    read_only: true
    cap_drop: [ALL]
    security_opt:
      - no-new-privileges:true
    networks:
      - proxy

networks:
  proxy:
    external: true
```

`TZ` sets the time zone used in the log file name and in the `time` field of every log line.

> [!NOTE]
> If you pass `-c` yourself (`docker run ... <image> -c other.toml`), the arguments are appended to the entrypoint, so `/config` is loaded first and your file is merged over it.

---

666 ) 𝙏͘͝𝙝̓̓͛𝙚͑̈́̀ 𝙨͆͠͝𝙠͑̾͌𝙮̽͌͆ 𝙞̓̔̔𝙨͒͐͝ 𝙛͑̈́̚𝙖͛͒𝙡͑͆̽𝙡̾̚̚𝙞͋̒̒𝙣̾͛͝𝙜͒̒̀.́̔͝​

## Command Line

| Flag | Effect |
| --- | --- |
| `-c`, `--config <file>` | Use this configuration file. **Built-in defaults are then not applied**, see the note below. |
| `-v`, `--debug` | Force debug level on |
| `-vv`, `--trace` | Force trace level on |
| `--nocolor` | Disable color and the banner |
| `--banner` | Print the banner and version, then exit |
| `--genconfig` | Write the default configuration to `./config.toml`, then exit |
| `-h` | Print usage (only when stdout is a terminal), then exit |

> [!NOTE]
> The usage text lists `--help`, but only `-h` is recognised: `--help` is ignored and the server starts. The usage text also says `--genconfig` writes `HellPot.toml`, but the file written is `config.toml`.

The version shown by `--banner` is the git tag the binary was built from (for example `1.2.3` for tag `v1.2.3`). Builds without a tag show the short commit hash, or `dev`.

## Configuration Reference

With no `-c` flag, HellPot loads `/etc/HellPot/config.toml` if it exists (not on Windows), otherwise `config.toml` in your user config directory, which is created from the defaults when missing. The built-in defaults fill every key your file does not set.

> [!IMPORTANT]
> With `-c`, **the built-in defaults are not applied.** Any key you leave out takes its zero value (empty string, `false`, `0`, empty list). For example a file with only the `[http]` bind settings and no `[http.router]` section registers no route at all: every path answers `404` and there is no `robots.txt`. Set `catchall = true` or `paths` explicitly. The one exception is `performance.max_conns_per_ip`, which stays at `10` when absent.

### Environment variables

Configuration values can be overridden with environment variables prefixed with `HELLPOT_`. Use one underscore for each `.` level and **two underscores** for every underscore inside a key name:

| Configuration key | Environment variable |
| --- | --- |
| `http.bind_addr` | `HELLPOT_HTTP_BIND__ADDR` |
| `http.router.catchall` | `HELLPOT_HTTP_ROUTER_CATCHALL` |
| `logger.docker_logging` | `HELLPOT_LOGGER_DOCKER__LOGGING` |
| `performance.max_conns_per_ip` | `HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP` |

### Reference

The block below is what `./HellPot --genconfig` produces (with comments added). Where the Docker image default differs, it is noted.

```toml
[deception]
  # Used as "Server" HTTP header. Note that reverse proxies may hide this.
  server_name = "nginx"

[http]
  # TCP Listener (default)
  bind_addr = "127.0.0.1"   # the Docker image uses "0.0.0.0"
  bind_port = "8080"

  # header name containing clients real IP, for reverse proxy deployments
  real_ip_header = 'X-Real-IP'

  # this contains a list of blacklisted useragent strings. (case sensitive)
  # clients with useragents containing any of these strings will receive "Not found" for any requests.
  # default: ["Cloudflare-Traffic-Manager"]; the Docker image adds "curl"
  uagent_string_blacklist = ["Cloudflare-Traffic-Manager", "curl"]

  # Unix Socket Listener (will override default). Not supported on Windows.
  unix_socket_path = "/var/run/hellpot"
  unix_socket_permissions = "0666"
  use_unix_socket = false

  [http.router]
    # Toggling this to true will cause all GET requests to match. Forces makerobots = false.
    catchall = false          # the Docker image uses true
    # Toggling this to false will prevent creation of robots.txt handler.
    makerobots = true
    # Handlers will be created for these paths, as well as robots.txt entries. Only valid if catchall = false.
    paths = ["wp-login.php", "wp-login"]

[logger]
  # verbose (-v)
  debug = true
  # extra verbose (-vv)
  trace = false
  # JSON log files will be stored in the below directory.
  # default when empty: $HOME/.local/share/HellPot/logs ; the Docker image uses "/logs/"
  directory = "/home/kayos/.local/share/HellPot/logs/"
  # disable all color in console output. when using Windows this will default to true.
  nocolor = false
  # toggles the use of the current date as the names for new log files.
  use_date_filename = true
  # log JSON to stdout only (no log file), implies nocolor. Meant for `docker logs` / log collectors.
  docker_logging = false
  # Go time layout used for the timestamp in console output (not in the JSON logs)
  console_time_format = "3:04PM"

[performance]
  # max_workers is only valid if restrict_concurrency is true
  max_workers = 256
  restrict_concurrency = false
  # max simultaneous connections from one remote IP, beyond which clients get HTTP 429.
  # Behind a reverse proxy every client has the proxy's IP, so this caps the whole instance.
  # 0 = unlimited. Added in this fork; the default, 10, is the value that used to be hard-coded.
  max_conns_per_ip = 10
```

`catchall = true` overrides `makerobots`: no `robots.txt` handler is created. With `catchall = false`, a request to a path that is not in `paths` (and is not `/robots.txt`) gets `404`.

## Behavior

- Only `GET` is served. Other methods are rejected by the HTTP server.
- A matching request receives `<html><body>` followed by an endless stream generated by a markov chain. The response only ends when the client goes away.
- The request read timeout is 5 s, the maximum request body is 1 MiB, keep-alive is disabled and a connection serves at most 2 requests.
- User agents containing a blacklisted string get `404 Not found` and are not trapped.
- `robots.txt` (when enabled) is generated from `paths`, one `Disallow:` line per path, so the bots that ignore it are the ones that get trapped.
- The `Server` header is `deception.server_name`.
- The client address that is logged is the value of `real_ip_header` when present, otherwise the TCP peer address.

## Logging

HellPot writes one JSON object per line (zerolog) and, unless `docker_logging` is on, also prints a human-readable version on the console. The log file is `<directory>/HellPot[_<start date and time>].log`, for example `HellPot_03_Oct_26_20-56_EDT.log`, named once at process start. With `docker_logging = true` the JSON goes to stdout and no file is written.

Per-request lines:

```json
{"level":"info","USERAGENT":"Mozilla/5.0","REMOTE_ADDR":"1.2.3.4","URL":"/wp-login.php","time":"2026-10-03T22:00:36Z","message":"NEW"}
{"level":"info","USERAGENT":"Mozilla/5.0","REMOTE_ADDR":"1.2.3.4","URL":"/wp-login.php","BYTES":179549,"DURATION":15.914698,"time":"2026-10-03T22:00:36Z","message":"FINISH"}
```

| `message` | Level | When | Extra fields |
| --- | --- | --- | --- |
| `NEW` | info | a client was sent to hell | `USERAGENT`, `REMOTE_ADDR`, `URL` |
| `FINISH` | info | the client left | `BYTES` sent, `DURATION` in milliseconds |
| `END_ON_ERR` | trace | the write failed (usually the client disconnected) | `error` |
| `Ignoring useragent` | trace | a blacklisted user agent was refused | |
| `SERVE_ROBOTS` | debug | `robots.txt` was served | `PATHS` |
| `The number of connections from <ip> exceeds MaxConnsPerIP=<n>` | debug | a client was refused with HTTP 429 | |

`caller` is added to request lines when `trace` is enabled. Startup lines (`config`, `logger`, `Listening and serving HTTP...`) report the config file, the log file and the listen address.

> [!NOTE]
> Known quirk, kept on purpose so that log content stays identical to upstream: the debug and trace lines are currently written even when `debug = false` and `trace = false`. These two settings can only raise the level.

## Example Reverse Proxy Configs

#### nginx

<details>
  <summary>nginx</summary>

```nginx
location '/robots.txt' {
	proxy_set_header Host $host;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_pass http://127.0.0.1:8080$request_uri;
}

location '/wp-login.php' {
	proxy_set_header Host $host;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_pass http://127.0.0.1:8080$request_uri;
}
```

</details>

#### Apache

<details>
  <summary>apache (mod_proxy + mod_proxy_http)</summary>

All nonexisting URLs are being reverse proxied to a HellPot instance on localhost, which is set to catchall. Traffic served by HellPot is rate limited to 5 KiB/s.

- Create your normal robots.txt and usual content. Also create the fake Errordocument directory and files (files can be empty). In the example, the directory is "/content/"
- A request on a URL with an existing handler (f.e. a file) will be handled by apache
- Requests on nonexisting URLs cause a HTTP Error 404, which content is served by HellPot
- URLs under the "/.well-known/" suffix are excluded.

```apache
<VirtualHost yourserver>
    ErrorDocument 400 "/content/400"
    ErrorDocument 403 "/content/403"
    ErrorDocument 404 "/content/404"
    ErrorDocument 500 "/content/405"
    <Directory "$wwwroot/.well-known/">
        ErrorDocument 400 default
        ErrorDocument 403 default
        ErrorDocument 404 default
        ErrorDocument 500 default
    </Directory>
    /* HTTP Honeypot / HellPot (need mod_proxy, mod_proxy_http) */
    ProxyPreserveHost	on
    ProxyPass         "/content/" "http://localhost:8080/"
    ProxyPassReverse  "/content/" "http://localhost:8080/"

    /* Rate Limit config, need mod_ratelimit */
    <Location "/content/">
        SetOutputFilter RATE_LIMIT
        SetEnv rate-limit 5
    </Location>

    /* Remaining config */

</VirtualHost>
```

</details>

## Security Notes

- **`real_ip_header` is trusted as is.** HellPot does not check who sent it. Only expose HellPot through a reverse proxy that overwrites this header (as the nginx example does), or keep it on an internal network, otherwise a client can forge the `REMOTE_ADDR` that ends up in your logs.
- The Docker image runs as root. It needs no added capability and no writable filesystem except `/logs`: the `read_only`, `cap_drop` and `no-new-privileges` settings of the compose example above are verified to work.
- The CI runs `go vet`, `gosec`, the tests with the race detector and `govulncheck` on every push (see below), and the Docker image is rebuilt weekly.

## Development

- **CI** ([go.yml](.github/workflows/go.yml), "Vibe Check"): on every push and on pull requests to `main`: `go vet`, `gosec`, `go test -race`, `go build`, `govulncheck`. The Go version comes from `go.mod`.
- **Docker** ([docker.yml](.github/workflows/docker.yml)): builds `linux/amd64` and `linux/arm64` and pushes to `ghcr.io/sirterrific/hellpot` (pull requests only build, they do not push). The Dockerfile runs `go vet` and `go test` before building, and accepts a `VERSION` build argument.
- **Releases** ([release-command.yml](.github/workflows/release-command.yml)): creating a GitHub release builds binaries for linux, windows, darwin and freebsd (386, amd64, arm64, except darwin/386 and windows/arm64), with SHA-256 sums.
- **Dependabot** checks Go modules and GitHub Actions daily and Docker base images weekly.

To cut a release, push a tag such as `v1.2.3`: the Docker workflow publishes the versioned image tags. Then create the GitHub release from that tag to get the binaries.

A generated knowledge base about the code base is available in the [openwiki](openwiki) directory.

## Related Suffering

- https://github.com/ginger51011/pandoras_pot
  - A HellPot inspired HTTP honeypot to punish and educate unruly web crawlers, written in Rust (🚀)

## Credits and License

HellPot was created by [yung innanet](https://github.com/yunginnanet) and is based on [Heffalump](https://github.com/carlmjohnson/heffalump) by Carl Johnson. This fork keeps the original [MIT license](LICENSE).
