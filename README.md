# Sovorem.am CLI

Official command-line tool for [Sovorem.am](https://sovorem.am/) — the Armenian backend development learning platform. Students use this CLI to run lesson tests locally and submit cryptographically signed results to the Sovorem.am backend.

## Table of Contents

- [Installation](#installation)
  - [1. Install Go](#1-install-go)
  - [2. Install the Sovorem CLI](#2-install-the-sovorem-cli)
  - [3. Login to the CLI](#3-login-to-the-cli)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Base URL for HTTP tests](#base-url-for-http-tests)
  - [CLI colors](#cli-colors)
  - [Troubleshooting the Config](#troubleshooting-the-config)
- [Upgrading](#upgrading)
  - [Troubleshooting Upgrading](#troubleshooting-upgrading)

## Installation

### 1. Install Go

To use the Sovorem CLI, you need an up-to-date Golang toolchain installed on your system.

Most courses are designed for Linux or macOS — or Linux-in-Windows via WSL. If you're on Windows, usually you'll want WSL and the Linux instructions below. Some courses also support Windows/PowerShell natively.

**Option 1 (Linux/WSL/macOS):** The [Webi installer](https://webinstall.dev/golang/) is the simplest way for most people:

```sh
curl -sS https://webi.sh/golang | sh
```

**Option 2 (any platform, including Windows/PowerShell):** Use the [official Golang installation instructions](https://go.dev/doc/install).

After installing Go, open a new shell session and run `go version` to confirm it works.

### 2. Install the Sovorem CLI

```sh
go install github.com/sovorem/sovorem@latest
```

Run `sovorem --version` to verify the installation.

If you get "command not found", add Go's bin directory to your `PATH` (usually `$HOME/go/bin`):

```sh
# Linux/WSL
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

```sh
# macOS
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc
```

```sh
# fish
fish_add_path $HOME/go/bin
```

### 3. Login to the CLI

Run `sovorem login` to authenticate with your Sovorem.am account. After authenticating, you're ready to run and submit lessons.

## Usage

| Command | Description |
|---------|-------------|
| `sovorem login` | Authenticate with your Sovorem.am account |
| `sovorem logout` | Disconnect the CLI from your account |
| `sovorem status` | Show login and version status |
| `sovorem run UUID` | Run a lesson's tests without submitting |
| `sovorem run UUID -s` | Run tests and submit in one step |
| `sovorem submit UUID` | Run tests and submit results to Sovorem.am |
| `sovorem config base_url URL` | Override the base URL for HTTP tests |
| `sovorem upgrade` | Install the latest CLI version |

Lesson UUIDs are shown on each CLI lesson page at [sovorem.am](https://sovorem.am).

## Configuration

The CLI stores settings in `~/.sovorem.yaml`, or `$XDG_CONFIG_HOME/sovorem/config.yaml` when `XDG_CONFIG_HOME` is set.

All commands support `-h`/`--help`.

### Base URL for HTTP tests

For lessons with HTTP tests, you can set a base URL that overrides the lesson default. This is useful when your local server runs on a different port.

```sh
sovorem config base_url http://localhost:8080/
sovorem config base_url
sovorem config base_url --reset
```

Include the protocol scheme (`http://` or `https://`) in the URL.

### CLI colors

Customize terminal output colors (success, error, secondary text):

```sh
sovorem config colors --red VALUE --green VALUE --gray VALUE
sovorem config colors
sovorem config colors --reset
```

Use an [ANSI color code](https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit) or hex string as `VALUE`.

### Troubleshooting the Config

To reset configuration completely, delete your config file. The CLI will recreate a fresh one on next run. Then run `sovorem login` again.

## Upgrading

The CLI checks for updates on login and before running authenticated commands.

```sh
sovorem upgrade
```

Or install a specific version:

```sh
go install github.com/sovorem/sovorem@v0.1.0
```

### Troubleshooting Upgrading

**Bypass the proxy** if you keep seeing the same upgrade message:

```sh
GOPROXY=direct go install github.com/sovorem/sovorem@latest
```

**Reinstall** if that doesn't work:

```sh
rm "$(which sovorem)"
go install github.com/sovorem/sovorem@latest
sovorem login
```

## Development

```sh
git clone https://github.com/sovorem/sovorem.git
cd sovorem
go test ./...
go build -o sovorem .
```

## License

See [LICENSE](LICENSE).
