# vibe-validator

[oo] Vibe check your dependencies.

`vibe-validator` is a cross-platform CLI tool that scans your project for suspicious or shady dependencies across Python (PyPI), JavaScript (npm), and Go modules.

It flags packages that:

- ❌ Don’t exist in public registries
- [~] Are recently published (less than 30 days old)
- [✓] Pass the vibe check

## 🧪 Supported Ecosystems

- **Python**: `requirements.txt`
- **Node.js**: `package.json`
- **Go**: `go.mod`

More to come: Dockerfiles, import statements, source scanning.

## 📦 Installation

Build from source (requires Go 1.20+):

```bash
git clone https://github.com/Kelcode-Dev/vibe-validator.git
cd vibe-validator
go build -o vibe-validator
````

Or install globally (once released):

```bash
go install github.com/Kelcode-Dev/vibe-validator@latest
```

## 🚀 Usage

```bash
vibe-validator <path-to-project>
```

Examples:

```bash
vibe-validator .
vibe-validator ./tests/npm
vibe-validator ~/code/my-cool-app
```

## ✅ Output Format

Terminal-friendly output:

```
[oo] Scanning: ./my-app

[vibe-validator] Dependency Vibe Report

pypi:
  [✓] requests             -
  [✗] shady-lib            Not found on PyPI

npm:
  [✓] express              -
  [✗] weird-package        Not found on npm

go:
  [~] github.com/sus/module     Recently added (fetched 4h ago)
```

## 🛠️ Roadmap

* [ ] GitHub repo validation (e.g. missing README, license, stars)
* [ ] Source file import scanning (`import`, `require`)
* [ ] Output options: `--json`, `--yaml`, `--markdown`
* [ ] CI-friendly exit codes (`--strict`)
* [ ] Package risk scores / badges
* [ ] Plugin system for custom registries (e.g. internal Artifactory)

## 👀 Logo

Coming soon: GOGO the dancer 🕺 with shades, clipboard, and heavy judgement.

## 📜 License

MIT — but if you use this to vibe-check your production stack, please consider buying your devs coffee ☕
