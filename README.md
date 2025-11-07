# Janitor – tidy your folders by file type

Janitor is a tiny CLI that organizes files in a directory into subfolders based on their extensions. It reads simple YAML rules (see `config.sample.yml`) to decide where each file should go, creates folders if they’re missing, and skips anything it doesn’t recognize.

Works great for quickly cleaning your Downloads or Desktop.

## Highlights

- Organizes files in the target directory (non-recursive) by extension
- Creates destination folders on demand; leaves unknown files untouched
- Safe moves: won’t overwrite existing files (it logs a conflict and skips)
- Skips the Janitor binary itself and `config.yml`
- Configurable via a small YAML file; sane defaults if no config is found

## Install / Build

You can either use package managers (Scoop, Homebrew) or build from source. This repo uses Go modules.

### Scoop (Windows)

Add the bucket and install Janitor:

```bash
scoop bucket add biohazard786 https://github.com/BioHazard786/scoop-bucket.git
scoop install biohazard786/janitor
```

Upgrade when new versions are released:

```bash
scoop update
scoop update janitor
```

### Homebrew (macOS / Linux with Homebrew)

Tap the repository and install the cask:

```bash
brew tap BioHazard786/tap
brew install --cask janitor
```

Upgrade:

```bash
brew upgrade janitor
```

If you prefer building from source instead of using the cask:

```bash
git clone https://github.com/BioHazard786/Janitor.git
cd Janitor
go build -o janitor ./cmd/janitor
./janitor -v
```

### Build from source (generic)

- Prerequisite: Go (the module declares `go 1.25.3`; building with the latest stable Go is fine).

```bash
# from the repository root
go build -o janitor ./cmd/janitor
./janitor -v
```

Run without building a binary first:

```bash
go run ./cmd/janitor -i
```

Note: `go install` with a module path requires the module to be published; since the module path here is local (`janitor`), prefer the build commands above.

## Usage

```text
janitor [flags] [DIRECTORY]

Flags:
  -v   Print Janitor's version
  -i   Print developer info
  -h   Show help (provided by the flag package)
```

- If `DIRECTORY` is omitted, Janitor uses the current working directory.
- Janitor only organizes top-level files in the target directory (doesn’t recurse into subfolders).

Examples (Windows, Git Bash):

```bash
# Clean the current folder
./janitor

# Clean your Downloads folder (quotes handle spaces)
./janitor "C:\\Users\\<you>\\Downloads"

# Show version / info
./janitor -v
./janitor -i
```

You’ll see logs such as which config was used and which files were moved.

## Configuration (config.yml)

Janitor looks for configuration in this order:

1. `./config.yml` (in the directory you run against)
2. `~/.config/janitor/config.yml` (per-user, global)
3. Built-in defaults (used if no config file is found)

The YAML schema is a single map of folder names to lists of file extensions. Extensions are case-insensitive, but must include the leading dot (e.g., `.mp4`).

Minimal example (from `config.sample.yml`):

```yaml
folders:
  Videos:
    - .mp4
    - .mov
    - .avi
    - .mkv
  Music:
    - .mp3
    - .wav
    - .flac
  Pictures:
    - .jpg
    - .jpeg
    - .png
    - .gif
  Illustrations: # custom folder
    - .psd
    - .ai
  Documents:
    - .pdf
    - .docx
    - .doc
    - .txt
    - .md
  Archives:
    - .zip
    - .tar
    - .gz
  Other:
    - .exe
    - .dmg
```

How it works:

- Janitor inverts this mapping so it can quickly look up a destination folder by a file’s extension.
- If a file’s extension matches a key in your config, the file is moved into the corresponding folder under the target directory (the folder is created if missing).
- If there’s no match, the file is left where it is.
- If a same-named file already exists in the destination, Janitor logs a conflict and skips moving that file.

### Where to put the config

- Project-specific rules: put `config.yml` directly in the directory you want to organize (e.g., your Downloads folder). This takes priority for that run.
- Global rules: put it at `~/.config/janitor/config.yml`.
  - On Windows with Git Bash, that expands to `C:\Users\<you>\.config\janitor\config.yml`.

### Using the sample config

1. Copy `config.sample.yml` to `config.yml`.
2. Adjust folder names and extension lists to your preference (keep the leading dot in extensions).
3. Place it either next to where you’ll run Janitor (per-folder rules) or in your user config directory (global).

## Built-in defaults (when no config file is found)

Janitor ships with a sensible default mapping that covers common categories, including (truncated list):

- Videos: .mp4, .mov, .avi, .mkv, .wmv, .webm, .flv, .mpg, .mpeg, .av1, .opus, .ts
- Music: .mp3, .wav, .flac, .aac, .ogg, .mpa, .m4a, .wma, .midi
- Pictures: .jpg, .jpeg, .png, .gif, .bmp, .svg
- Documents: .pdf, .docx, .doc, .txt, .md, .xls, .xlsx, .ppt, .pptx, .cbz, .cbr
- Archives: .zip, .tar, .gz, .rar, .7z, .xz, .bz2
- Applications: .exe, .msi, .apk, .dmg, .deb, .rpm, .appx, .msix

You can fully override or extend these via your own `config.yml`.

## Behavior and limitations

- Scope: Only files in the target directory’s top level are considered (no recursion into subfolders).
- Moves within the same directory tree using `os.Rename` (fast and atomic on the same volume).
- Conflicts: If the destination file already exists, the move is skipped with a log message.
- Safety: Janitor skips moving itself and `config.yml`.

## Development

Project layout:

- `cmd/janitor/main.go` – CLI entrypoint and flags
- `internal/config/` – config struct, defaults, and loading logic
- `internal/organizer/` – core move logic
- `config.sample.yml` – example configuration to copy and tweak

Run all builds:

```bash
go build ./...
```

## FAQ

- Can I map the same extension to multiple folders? No—each extension resolves to a single destination. If you list it twice, the last one wins.
- Do extensions require a leading dot? Yes. Use `.mp4`, not `mp4`.
- Will it search subfolders? No—only the top-level files in the target directory.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
