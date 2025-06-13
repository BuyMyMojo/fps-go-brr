# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI application for professional video frame analysis and comparison. The project provides tools to count frames in videos, compare individual frames, analyze differences between two videos, and perform comprehensive frame persistence analysis for single videos with DigitalFoundry-style CSV output.

## Core Architecture

The application is built using:
- **CLI Framework**: urfave/cli/v3 for command-line interface
- **Video Processing**: AlexEidt/Vidio library for video file handling with FPS detection
- **Image Processing**: Standard Go image libraries for frame comparison
- **CSV Export**: Built-in CSV generation for professional video analysis visualization

### Main Components

- **CLI Commands**: Five main commands for comprehensive frame analysis operations
- **Frame Comparison**: Pixel-level comparison with configurable tolerance using squared difference
- **Video Processing**: Frame-by-frame video analysis with streaming support and memory-efficient processing
- **Two-Pass Analysis**: Advanced frame persistence analysis with pre-calculated total durations
- **CSV Generation**: DigitalFoundry-style data export for professional visualization tools

### Key Functions

- `count_video_frames()`: Counts total frames in a video file
- `compare_frames()`: Compares two frames with tolerance-based difference detection
- `compare_frames_alt()`: Alternative frame comparison using exact pixel matching
- `countUniqueVideoFrames()`: Analyzes differences between corresponding frames in two videos
- `analyzeFramePersistence()`: **Main feature** - Two-pass frame persistence analysis with CSV export
- `isDiffUInt8WithTolerance()`: Pixel comparison with configurable tolerance threshold
- `imageToRGBA()`: Converts images to RGBA format for consistent processing

## Development Commands

### Build and Run
```bash
# Normal build
go build -o fps-go-brr .
./fps-go-brr <command> [args]

# Optimized compact build (requires UPX)
./build-compact.sh
./fps-go-brr-compact <command> [args]
```

### Testing
```bash
go test ./...
```

### Module Management
```bash
go mod tidy
go mod download
```

### Release Builds
- Forgejo Actions automatically build and release both normal and compact binaries on tag pushes
- Uses custom runner: `9950x`
- Normal and compact builds are uploaded as separate artifacts
- UPX compression applied to compact builds for size optimization

## Repository Information

- **Main Repository**: https://git.aria.coffee/aria/fps-go-brr (Personal Forgejo instance)
- **Mirror**: https://github.com/BuyMyMojo/fps-go-brr (GitHub - accepts PRs and issues)
- **Dual Licensed**: MIT OR Apache-2.0 (SPDX-License-Identifier: MIT OR Apache-2.0)
- **Copyright**: 2025 Aria, Wicket

### Inspirations

This project draws inspiration from:
- Digital Foundry (YouTube) - Professional video game performance analysis
- Brazil Pixel (YouTube) - Technical video analysis and frame rate studies  
- TRDrop (GitHub) - Raw video analysis program for framerate estimation
- Original Python implementation - Early proof-of-concept for frame persistence analysis

## Memories

- The forgejo workflow runner is executed as root so it does not need to use root

## CLI Usage

Available commands:
- `count-frames <video>` - Count frames in a video
- `compare-frames <frame1> <frame2>` - Compare two image frames
- `count-frames-differing-pixels <frame1> <frame2>` - Count pixel differences between frames
- `count-unique-video-frames <video1> <video2>` - Compare corresponding frames between two videos
- `analyze-frame-persistence [--tolerance float] [--csv-output path] <video>` - **Main feature**: Professional video analysis with CSV export

### Frame Persistence Analysis with CSV Export

The main feature provides:
- Real-time FPS detection from video metadata
- Frame-by-frame comparison with previous frame
- Detection of consecutive duplicate frame sequences (3+ identical frames)
- Per-second unique frame counting
- Two-pass analysis for accurate total frame persistence calculation
- Configurable pixel difference tolerance (0-255)
- **Professional CSV export** with 5 columns for DigitalFoundry-style analysis

### CSV Output Format

The `--csv-output` flag generates a CSV file with these columns:
- `frame`: Frame number (1-based, no skipped frames)
- `average_fps`: Running effective FPS calculation
- `frame_time`: Current frame persistence duration (real-time)
- `unique_frame_count`: Cumulative unique frame count (stays constant during duplicates)
- `real_frame_time`: **Total persistence time for each unique frame (smooth for visualization)**

### CSV Usage Examples

```bash
# Basic analysis with CSV export
./fps-go-brr analyze-frame-persistence video.mp4 --csv-output analysis.csv

# With tolerance for noisy videos
./fps-go-brr analyze-frame-persistence video.mp4 --tolerance 10 --csv-output analysis.csv
```

## Advanced Implementation Details

### Two-Pass Analysis Architecture

The `analyzeFramePersistence()` function uses a sophisticated two-pass approach:

1. **First Pass**: Analyzes entire video to calculate total duration each unique frame will persist
2. **Second Pass**: Writes CSV with correct `real_frame_time` values for smooth visualization

This ensures:
- All instances of the same unique frame show identical `real_frame_time` values
- Creates smooth, non-jumpy graphs perfect for professional video analysis
- DigitalFoundry-style frame timing visualization compatibility

### Frame Data Structure

```go
type FrameData struct {
    frameNumber      int     // Current frame number
    uniqueFrameCount int     // Cumulative unique frames
    effectiveFPS     float64 // Running average FPS
    currentFrameTime float64 // Current persistence so far
    realFrameTime    float64 // Total persistence duration
}
```

## Implementation Notes

- The application processes video frames in memory using RGBA format
- Pixel comparison uses squared difference for tolerance-based matching
- Video processing is done frame-by-frame to handle large files efficiently
- Frame persistence detection only reports sequences of 3+ consecutive identical frames
- Two-pass analysis ensures accurate total persistence calculations for visualization
- CSV output is optimized for professional video analysis tools and graphing software
- The `analyze-frame-persistence` command is the primary tool for professional video quality analysis
- All image formats supported by Go's image package can be used for frame comparison