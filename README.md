# genQ — QR Code Generator

A cross-platform desktop QR code generator built with Go and [Fyne](https://fyne.io).

## Features

- **Customisable QR codes** — generate from any URL or text
- **Data dot styles** — Square, Circular, or Rounded
- **Eye styles** — Square, Rounded, or Circle for the finder patterns
- **Foreground & background colours** — full colour picker with live preview
- **Size presets** — Screen (200px–1200px) and Print (1″–6″ @ 300 DPI)
- **Error correction levels** — Low, Medium, Quartile, High
- **Save as PNG** or copy to clipboard
- **Dark / Light mode** toggle

## Screenshot

_TODO: Add screenshot_

## Requirements

- **Go 1.21+**
- C compiler (CGO is required by Fyne)
- Linux: `libgl1-mesa-dev`, `xorg-dev` (or Wayland equivalents)
- macOS: Xcode command-line tools
- Windows: MinGW-w64 (for cross-compilation from Linux)

## Getting Started

```bash
# Clone
git clone git@github.com:thesle/genQ.git
cd genQ

# Run
go run .

# Or build and run
make build-linux
./build/genQ-linux-amd64
```

## Build Targets

| Command              | Description                          |
|----------------------|--------------------------------------|
| `make run`           | Run in dev mode                      |
| `make build-linux`   | Build Linux amd64 binary             |
| `make build-windows` | Build Windows amd64 binary           |
| `make build-macos`   | Build macOS amd64 + arm64 binaries   |
| `make build-all`     | Build for all platforms              |
| `make flatpak`       | Package as Flatpak (GNOME 49)        |
| `make clean`         | Remove build artifacts               |

## Flatpak

The Flatpak manifest uses the GNOME 49 runtime and packages the pre-compiled Linux binary — no recompilation during the Flatpak build.

```bash
make flatpak
```

Requires `flatpak-builder` to be installed.

## Project Structure

```
genQ/
├── main.go              # App entry point, UI wiring
├── qr/
│   └── render.go        # QR matrix generation + image rendering
├── ui/
│   ├── controls.go      # Input controls (text, styles, colours, sizes)
│   ├── preview.go       # QR code preview widget
│   └── theme.go         # Dark/light mode support
├── flatpak/
│   ├── com.genq.GenQ.yml          # Flatpak manifest
│   ├── com.genq.GenQ.desktop      # Desktop entry
│   ├── com.genq.GenQ.metainfo.xml # AppStream metadata
│   └── com.genq.GenQ.svg          # App icon
├── Makefile
├── go.mod
└── go.sum
```

## License

MIT
