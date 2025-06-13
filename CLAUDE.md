# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI application for video frame analysis and comparison. The project provides tools to count frames in videos, compare individual frames, analyze differences between two videos, and perform frame persistence analysis for single videos.

## Core Architecture

The application is built using:
- **CLI Framework**: urfave/cli/v3 for command-line interface
- **Video Processing**: AlexEidt/Vidio library for video file handling with FPS detection
- **Image Processing**: Standard Go image libraries for frame comparison

### Main Components

- **CLI Commands**: Five main commands for comprehensive frame analysis operations
- **Frame Comparison**: Pixel-level comparison with configurable tolerance using squared difference
- **Video Processing**: Frame-by-frame video analysis with streaming support and memory-efficient processing
- **Frame Persistence Analysis**: Detects consecutive duplicate frames and calculates persistence duration

### Key Functions

- `count_video_frames()`: Counts total frames in a video file
- `compare_frames()`: Compares two frames with tolerance-based difference detection
- `compare_frames_alt()`: Alternative frame comparison using exact pixel matching
- `countUniqueVideoFrames()`: Analyzes differences between corresponding frames in two videos
- `analyzeFramePersistence()`: **Main feature** - Analyzes frame persistence in single video with per-second statistics
- `isDiffUInt8WithTolerance()`: Pixel comparison with configurable tolerance threshold
- `imageToRGBA()`: Converts images to RGBA format for consistent processing

## Development Commands

### Build and Run
```bash
go build -o fps-go-brr .
./fps-go-brr <command> [args]
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

## CLI Usage

Available commands:
- `count-frames <video>` - Count frames in a video
- `compare-frames <frame1> <frame2>` - Compare two image frames
- `count-frames-differing-pixels <frame1> <frame2>` - Count pixel differences between frames
- `count-unique-video-frames <video1> <video2>` - Compare corresponding frames between two videos
- `analyze-frame-persistence [--tolerance float] <video>` - **Main feature**: Analyze frame persistence and unique frames per second

### Frame Persistence Analysis

The main feature provides:
- Real-time FPS detection from video metadata
- Frame-by-frame comparison with previous frame
- Detection of consecutive duplicate frame sequences (3+ identical frames)
- Per-second unique frame counting
- Persistence duration calculation in milliseconds
- Configurable pixel difference tolerance (0-255)

## Implementation Notes

- The application processes video frames in memory using RGBA format
- Pixel comparison uses squared difference for tolerance-based matching
- Video processing is done frame-by-frame to handle large files efficiently
- Frame persistence detection only reports sequences of 3+ consecutive identical frames
- All image formats supported by Go's image package can be used for frame comparison
- The `analyze-frame-persistence` command is the primary tool for video quality analysis