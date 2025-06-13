# fps-go-brr

> **⚠️ Work in Progress**: This project is actively under development and not yet feature-complete!

A Go CLI tool for video frame analysis and comparison. Analyze frame persistence, detect dropped frames, and export data for visualization tools like those used by Digital Foundry.

## Features

- **Frame Persistence Analysis**: Detect consecutive duplicate frames and measure persistence duration
- **CSV Export**: Generate data compatible with video analysis visualization tools
- **Multi-format Support**: Works with any video format supported by FFmpeg
- **Configurable Tolerance**: Adjust pixel difference sensitivity for noisy videos
- **Real-time Analysis**: Stream processing for efficient memory usage
- **Two-pass Architecture**: Accurate frame timing calculations for smooth visualizations

## Quick Start

### Installation

Download the latest release from the [releases page](https://git.aria.coffee/aria/fps-go-brr/releases) or build from source:

```bash
# Clone the repository
git clone https://git.aria.coffee/aria/fps-go-brr.git
cd fps-go-brr

# Build normally
go build -o fps-go-brr .

# Or build compact version (requires UPX)
./build-compact.sh
```

### Basic Usage

```bash
# Analyze frame persistence with CSV export
./fps-go-brr analyze-frame-persistence video.mp4 --csv-output analysis.csv

# With tolerance for noisy videos
./fps-go-brr analyze-frame-persistence video.mp4 --tolerance 10 --csv-output analysis.csv

# Count total frames in a video
./fps-go-brr count-frames video.mp4

# Compare two individual frames
./fps-go-brr compare-frames frame1.png frame2.png
```

## CSV Output Format

The `analyze-frame-persistence` command generates CSV files with the following columns:

| Column | Description |
|--------|-------------|
| `frame` | Frame number (starts on 1) |
| `average_fps` | Running effective FPS calculation |
| `frame_time` | Current frame persistence duration (ms) |
| `unique_frame_count` | Cumulative unique frame count |
| `real_frame_time` | Total persistence time for smooth visualization |

## Use Cases

- **Game Performance Analysis**: Detect frame drops and stuttering in gameplay footage
- **Technical Reviews**: Generate data for Digital Foundry-style analysis

## Development Status

This project is under active development. Current feature wish list:

- [ ] Enhanced frame comparison algorithms
- [ ] Performance optimizations for large videos
- [ ] Additional export formats
- [ ] Cross-platform testing and compatibility
- [ ] Documentation improvements
- [ ] Graph generation from CSV

## Building

### Prerequisites

- Go 1.21 or later
- UPX (optional, for compact builds)

### Commands

```bash
# Standard build
go build -o fps-go-brr .

# Compact build with UPX compression
./build-compact.sh
```

## Repository

**Main Repository**: [https://git.aria.coffee/aria/fps-go-brr](https://git.aria.coffee/aria/fps-go-brr)  
**Mirror (GitHub)**: [https://github.com/BuyMyMojo/fps-go-brr](https://github.com/BuyMyMojo/fps-go-brr)

The main development happens on the personal Forgejo instance. The GitHub mirror also accepts pull requests and bug reports for convenience.

## Contributing

This is an early-stage project. Contributions, bug reports, and feature requests are welcome on either the main repository or the GitHub mirror!

## Technical Details

Built with:
- **CLI Framework**: [urfave/cli/v3](https://github.com/urfave/cli)
- **Video Processing**: [AlexEidt/Vidio](https://github.com/AlexEidt/Vidio)
- **Image Processing**: Go standard library

<!-- The tool uses a sophisticated two-pass analysis architecture to ensure accurate frame timing calculations for professional visualization tools. -->
<!-- ??? -->

## Inspirations

This project draws inspiration from:

- **[Digital Foundry](https://www.youtube.com/@DigitalFoundry)** - Professional video game performance analysis and technical reviews
- **[Brazil Pixel](https://www.youtube.com/@brazilpixel)** - Technical video analysis and frame rate studies
- **[TRDrop](https://github.com/cirquit/trdrop)** - Raw video analysis program for framerate estimation and tear detection
- **[Original Python implementation](https://web.archive.org/web/20250613174657/https://snippets.aria.coffee/snippets/2)** - Early proof-of-concept for frame persistence analysis

The goal is to provide similar professional-grade video analysis capabilities for the open-source community.

## License

SPDX-License-Identifier: MIT OR Apache-2.0

This project is dual-licensed under your choice of:
- MIT License - see [LICENSE.MIT](LICENSE.MIT) file for details
- Apache License 2.0 - see [LICENSE.Apache-2.0](LICENSE.Apache-2.0) file for details

Copyright (c) 2025 Aria, Wicket

---
