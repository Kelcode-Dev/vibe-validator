# Vibe Validator

[![GitHub release](https://img.shields.io/github/v/release/Kelcode-Dev/vibe-validator?style=flat-square)](https://github.com/Kelcode-Dev/vibe-validator/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kelcode-Dev/vibe-validator?style=flat-square)](https://goreportcard.com/report/github.com/Kelcode-Dev/vibe-validator)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Kelcode-Dev/vibe-validator)](https://pkg.go.dev/github.com/Kelcode-Dev/vibe-validator)

[oo] Vibe check your dependencies.

`vibe-validator` is a cross-platform CLI tool that scans your project for suspicious or shady dependencies across Python (PyPI), JavaScript (npm), and Go modules.

It flags packages that:

- `[x]` Don't exist in public registries
- `[~]` Are recently changed (less than 30 days old)
- `[✓]` Pass the vibe check

## 🧪 Supported Ecosystems

- **Python**: `requirements.txt`, `Pipfile.lock` (lockfile support via `--include-lockfiles`)
- **Node.js**: `package.json` (includes `dependencies` & `devDependencies`), `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml` (lockfile support via `--include-lockfiles`)
- **Go**: `go.mod`

More to come: Dockerfiles, source import scanning.

## 📦 Installation

Build from source (requires Go 1.24+):

```bash
git clone https://github.com/Kelcode-Dev/vibe-validator.git
cd vibe-validator
go build -o vibe-validator
````

Or install globally (once released):

```bash
go install github.com/Kelcode-Dev/vibe-validator@0.2.0
```

## 🚀 Usage

```bash
vibe-validator <path-to-project> [--include-lockfiles]
```

Examples:

```bash
vibe-validator .
vibe-validator ./tests/npm --include-lockfiles # includes package lock files
vibe-validator ~/code/my-cool-app --include-vendor # includes vendor specific package files
vibe-validator ./tests -vv # max verbosity
```

## ✅ Output Format

Terminal-friendly output:

```
[oo] Scanning: ./my-app

[vibe-validator] Dependency Vibe Report

pypi:
  [✓] requests             -                  tests/pypi/requirements.txt
  [✗] shady-lib            Not found on PyPI  tests/pypi/Pipfile.lock

npm:
  [✓] express              -                  tests/npm/package-lock.json, tests/npm/package.json
  [✗] weird-package        Not found on npm   tests/npm/package-lock.json

go:
  [~] github.com/sus/module     Recently added (3 days ago)  tests/go/go.mod
```

### Verbosity Levels

* By default (no verbosity flags), only packages needing attention are shown: [✗] (not found) and [~] (investigate)
* `-v` adds all [✓] (safe) packages to the output
* `-vv` includes a count and detailed scanning logs of all dependencies found (including duplicates)

## 🛠️ Roadmap

* [ ] GitHub repo validation (e.g. missing README, license, stars)
* [ ] Source file import scanning (`import`, `require`)
* [ ] Output options: `--json`, `--yaml`, `--markdown`
* [ ] CI-friendly exit codes (`--strict`)
* [ ] Package risk scores / badges
* [ ] New validators for PHP Composer, Ruby Gemfiles, extensions to existing validators for things like poetry etc.

## 📜 License

MIT — but if you use this to vibe-check your production stack, please consider buying your devs coffee ☕

## 📝 Changelog

### 0.2.0 - 2025-06-30
- Refactored scanner architecture to driver-based pattern for easier extensibility
- Added PHP, Ruby, and Rust dependency scanners and validators
- Improved CLI with verbosity flags for progressive output

### 0.1.1 - 2025-06-30
- Support for vendor folders and lockfile scanning (`--include-vendor`, `--include-lockfiles`)

### 0.1.0 - 2025-06-30
- Initial release supporting Python, JavaScript, and Go dependencies
- Basic validation of package existence and recent publish dates
- Terminal-friendly dependency vibe report output
