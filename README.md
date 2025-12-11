# radarlance

`radarlance` is a reconnaissance utility that monitors remote web assets - JavaScript & HTML - for changes over time. It fetches, beautifies, hashes, and archives each version of a target resource and reports only when meaningful changes occur. Designed for bug bounty hunters, red teams, and security researchers tracking attack surfaces.

> **WARNING**: This tool is for **authorized security testing only**.
> Unauthorized use may violate laws and regulations.
> The author and contributors are not responsible for misuse.
> Always obtain explicit permission before testing any system.

---

## Features

- **Change Detection**: Tracks JS and HTML endpoints for content modifications using SHA-1 hashing.
- **Beautified Archiving**: Automatically formats and saves readable versions of changed resources for analysis.
- **Noise-Free Output**: Quiet mode (`-q`) ensures no output unless changes are detected—ideal for cron + notify pipelines.
- **Concurrent Fetching**: Configurable worker threads for high-performance monitoring.
- **Lightweight & Portable**: Single binary with no runtime dependencies beyond Go's standard library.

---

## Installation

### Prerequisites

- **Go** (1.21 or later)
- **Make** (for optional multi-platform builds)
- **Git**

### Steps

Build for all platforms:

```
$ make all
```

Or build a specific target:

```
$ make linux-amd64
$ make windows-amd64
$ make darwin-arm64
```

Binaries will appear under the `build/` directory.

---

## Usage

### Command-Line Flags

```
$ ./radarlance
Usage of ./radarlance
  -d string
    	base directory for saved files and hashes.json (default "data")
  -i string
    	input file containing URLs (one per line)
  -o string
    	output file to store JS file hashes (default "hashes.json")
  -q	quiet mode (no completion output)
  -t int
    	number of concurrent threads (default 10)
  -type string
    	content type: js or html (default "js")
  -u string
    	single URL to check
  -v	enable verbose output
```

---

## Examples

### Monitor a list of JS files in quiet mode (cron + notify friendly)

```
$ katana -u https://target.com | grep -E '\.js(\?|$)' > js-files.txt
$ radarlance -q -type js -i js-files.txt
```

### Monitor HTML endpoints from a discovery tool

```
$ katana -u https://target.com | grep -E '\.html?$' > html.txt
$ radarlance -q -type html -i html.txt
```

### Check a single URL

```
radarlance -u https://target.com/app.js
```

Example output when changes are detected:

```
[CHANGED] https://target.com/app.js
    Old file: data/target.com/2025-12-11_13/9bc109.../app.js
    New file: data/target.com/2025-12-11_14/62dfe2.../app.js
```

The two paths can be diffed using GNU diff tools:

```
diff -u <old> <new>
```

---

## Disclaimer

`radarlance` is provided "as is" without warranties.
The authors assume no responsibility for misuse.
Use only for research or authorized security assessments.

---

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE.
See the [LICENSE](LICENSE) file for more details.
