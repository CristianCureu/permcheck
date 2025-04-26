# Permcheck

[![Go CI](https://github.com/cristiancureu/permcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/cristiancureu/permcheck/actions) [![Go Version](https://img.shields.io/badge/Go-1.24-blue)](https://golang.org/dl/)

**Permcheck** is a lightweight CLI tool for scanning directories and identifying insecure file permissions.  
It can optionally fix permissions automatically according to recommended security standards.

## Features

- Detects world-writable files (e.g., `0666`)
- Flags sensitive files (`.env`, `id_rsa`, `config.yaml`) with incorrect permissions
- Optional automatic fixing of insecure permissions
- Parallel scanning using worker pools
- Minimal, readable terminal output

## Installation

Clone and build the project:

```bash
git clone https://github.com/cristiancureu/permcheck.git
cd permcheck
make build
```

(Optional) Install the binary globally:

```bash
sudo cp permcheck /usr/local/bin/
```

## Usage

```bash
permcheck scan [directory] [flags]
```

## Examples

Scan the current directory:

```bash
permcheck scan .
```

Scan a directory and show only insecure files:

```bash
permcheck scan ./mock_files --insecure-only
```

Scan and fix insecure file permissions:

```bash
permcheck scan ./mock_files --fix
```

## Flags

| Flag            | Description                                           |
| --------------- | ----------------------------------------------------- |
| -w, --workers   | Number of concurrent workers (default: CPU cores × 2) |
| --insecure-only | Show only files with insecure permissions             |
| --fix           | Automatically fix detected insecure permissions       |

## Development

Run all tests:

```bash
make test
```

Format the code:

```bash
make fmt
```

Tidy the Go modules:

```bash
make tidy
```
