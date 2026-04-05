# Villain Couch

[![Banner Image](https://raw.githubusercontent.com/fezcode/Villain-Couch/main/banner.png)](https://github.com/fezcode/Villain-Couch/)

**Villain Couch** is a lightweight, cross-platform progress tracker and binge-watching assistant for VLC Media Player. It automatically saves your playback position and provides a sleek Terminal User Interface (TUI) to control your viewing experience.

## Key Features

- **Automated Tracking**: Automatically saves your progress in movies and TV shows to a local SQLite database.
- **Interactive TUI**: A beautiful terminal interface built with **Bubble Tea** and **Lip Gloss**.
- **Cross-Platform**: Supports Windows, Linux, and macOS (AMD64 & ARM64).
- **Binge Watching**: Automatically finds and plays the next episode in a series.
- **Watching History**: Easily browse and resume your recently watched media.
- **VLC Integration**: Direct control over VLC playback, volume, and seeking without terminal clutter.

## TUI Controls

| Key | Action |
| --- | --- |
| `Space` | Toggle Pause/Play |
| `n` | Play Next Episode |
| `h` / `l` | Seek Backward/Forward (10s) |
| `+` / `-` | Increase/Decrease Volume |
| `Tab` | Toggle between **Now Playing** and **History** |
| `Enter` | Resume selected file from History |
| `q` | Quit |

## Building

Villain Couch uses **[gobake](https://github.com/fezcode/gobake)** for its build system.

### Prerequisites

- Go 1.25+
- `gobake` installed: `go install github.com/fezcode/gobake/cmd/gobake@latest`

### Build Tasks

To build for all supported platforms (Windows, Linux, macOS):

```powershell
gobake build
```

The binaries will be generated in the `build/` directory.

To build and run locally:

```powershell
go run ./agent/src/
```

## Configuration

For advanced configuration (HTTP ports, passwords, custom VLC paths), please refer to the `config.json` file generated in the application's storage directory.

## Directory Structure

- `agent/src`: Core application logic and TUI.
- `common`: Shared utility packages.
- `build`: Compiled binaries (after running `gobake build`).
- `Recipe.go`: Build recipe for `gobake`.

## License

This project is licensed under the MIT License. See the `recipe.piml` for details.
