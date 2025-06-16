// A Go CLI tool for video frame analysis and comparison. Analyze frame persistence, detect dropped frames, and export data for visualization tools like those used by Digital Foundry.
package main

import (
	"cmp"
	"context"
	"encoding/csv"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/cheggaaa/pb"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "count-frames",
				Usage: "Count frames",
				Action: func(ctx context.Context, cmd *cli.Command) error {

					return countVideoFrames(cmd.Args().First())
				},
			},
			{
				Name:  "compare-frames",
				Usage: "Are two frames different?",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "frame1",
					},
					&cli.StringArg{
						Name: "frame2",
					},
				},

				Action: func(ctx context.Context, cmd *cli.Command) error {

					firstFrame, _ := getImageFromFilePath(cmd.StringArg("frame1"))
					secondFrame, _ := getImageFromFilePath(cmd.StringArg("frame2"))

					firstRGBA := imageToRGBA(firstFrame)
					secondRGBA := imageToRGBA(secondFrame)
					return compareFrames(firstRGBA, secondRGBA)
				},
			},
			{
				Name:  "count-frames-differing-pixels",
				Usage: "Are two frames different?",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "frame1",
					},
					&cli.StringArg{
						Name: "frame2",
					},
				},

				Action: func(ctx context.Context, cmd *cli.Command) error {

					firstFrame, _ := getImageFromFilePath(cmd.StringArg("frame1"))
					secondFrame, _ := getImageFromFilePath(cmd.StringArg("frame2"))

					firstRGBA := imageToRGBA(firstFrame)
					secondRGBA := imageToRGBA(secondFrame)
					return compareFramesAlt(firstRGBA, secondRGBA)
				},
			},
			{
				Name:  "count-unique-video-frames",
				Usage: "Are two frames different?",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "video1",
					},
					&cli.StringArg{
						Name: "video2",
					},
				},

				Action: func(ctx context.Context, cmd *cli.Command) error {
					return countUniqueVideoFrames(cmd.StringArg("video1"), cmd.StringArg("video2"), 1, false)
				},
			},
			{
				Name:  "analyze-frame-persistence",
				Usage: "Analyze frame persistence in a single video",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "video",
					},
				},
				Flags: []cli.Flag{
					&cli.Float64Flag{
						Name:  "tolerance",
						Usage: "Pixel difference tolerance (0-255)",
						Value: 0,
					},
					&cli.StringFlag{
						Name:  "csv-output",
						Usage: "Path to CSV file for frame data output",
						Value: "",
					},
					&cli.BoolFlag{
						Name:  "resdet",
						Usage: "use the resdet cli to measure each frame's resoltion\nWARNING: This will slow the process down by a LOT",
						Value: false,
					},

					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "print out total unique frames for every second of measurements",
						Value: false,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					tolerance := uint64(cmd.Float64("tolerance"))
					csvOutput := cmd.String("csv-output")
					return analyzeFramePersistence(cmd.StringArg("video"), tolerance, csvOutput, cmd.Bool("resdet"), cmd.Bool("verbose"))
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// countVideoFrames
// Prints out the total ammount of frames within `video`
//
// Parameters:
//   - video string - path to video file
//
// Returns:
//   - error
func countVideoFrames(video string) error {
	log.Default().Print("Trying to open video at: " + video)
	videoFile, _ := vidio.NewVideo(video)
	count := 0
	for videoFile.Read() {
		count++
	}
	log.Default().Println("Video total frames: " + strconv.Itoa(count))
	return nil
}

func compareFrames(frame1 *image.RGBA, frame2 *image.RGBA) error {
	accumError := int64(0)
	for i := 0; i < len(frame1.Pix); i++ {

		if isDiffUInt8WithTolerance(frame1.Pix[i], frame2.Pix[i], 0) { // Set tolerance to 0
			accumError++
		}
	}
	log.Default().Println("Total differing pixels: " + strconv.FormatInt(accumError, 10))
	return nil
}

func compareFramesAlt(frame1 *image.RGBA, frame2 *image.RGBA) error {
	// diff_frame := image.NewRGBA(frame1.Rect)
	accumError := int64(0)
	for i := 0; i < len(frame1.Pix); i++ {
		if isDiffUInt8(frame1.Pix[i], frame2.Pix[i]) {
			accumError++
		}
	}
	log.Default().Println("Total differing pixels: " + strconv.FormatInt(accumError, 10))
	return nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func isDiffUInt8(x, y uint8) bool {
	d := uint64(x) - uint64(y)
	sq := d * d
	if sq > 0 {
		return true
	}

	return false

}

func isDiffUInt8WithTolerance(x, y uint8, tolerance uint64) bool {
	d := uint64(x) - uint64(y)
	sq := d * d
	if sq > tolerance {
		return true
	}

	return false

}

func countUniqueVideoFrames(videoPath1 string, videoPath2 string, minDiff uint64, useSqDiff bool) error {
	video1, _ := vidio.NewVideo(videoPath1)
	video2, _ := vidio.NewVideo(videoPath2)
	video1Frame := image.NewRGBA(image.Rect(0, 0, video1.Width(), video1.Height()))
	video2Frame := image.NewRGBA(image.Rect(0, 0, video2.Width(), video2.Height()))
	video1.SetFrameBuffer(video1Frame.Pix)
	video2.SetFrameBuffer(video2Frame.Pix)
	totalFrames := 0
	uniqueFrames := 0
	for video1.Read() {
		totalFrames++
		video2.Read()
		accumError := uint64(0)
		for i := 0; i < len(video1Frame.Pix); i++ {
			if useSqDiff {
				if isDiffUInt8WithTolerance(video1Frame.Pix[i], video2Frame.Pix[i], minDiff) {
					accumError++
				}
			} else {
				if isDiffUInt8(video1Frame.Pix[i], video2Frame.Pix[i]) {
					accumError++
				}
			}
		}
		if minDiff <= accumError {
			uniqueFrames++
			log.Default().Println("[" + strconv.Itoa(totalFrames) + "]Unique frame")
		} else {
			log.Default().Println("[" + strconv.Itoa(totalFrames) + "]Non-unique frame")
		}
	}
	video1.Close()
	video2.Close()
	log.Default().Println(strconv.Itoa(uniqueFrames) + "/" + strconv.Itoa(totalFrames) + " are unique!")
	return nil
}

func imageToRGBA(src image.Image) *image.RGBA {

	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func analyzeFramePersistence(videoPath string, tolerance uint64, csvOutput string, toggleResdet bool, verbose bool) error {
	video, err := vidio.NewVideo(videoPath)
	if err != nil {
		return err
	}
	defer video.Close()

	fps := video.FPS()
	frameTimeMs := 1000.0 / fps

	log.Default().Printf("Video FPS: %.2f, Frame time: %.2f ms", fps, frameTimeMs)

	var csvWriter *csv.Writer
	var csvFile *os.File
	if csvOutput != "" {
		csvFile, err = os.Create(csvOutput)
		if err != nil {
			return fmt.Errorf("failed to create CSV file: %v", err)
		}
		defer csvFile.Close()

		csvWriter = csv.NewWriter(csvFile)
		defer csvWriter.Flush()

		err = csvWriter.Write([]string{"frame", "average_fps", "frame_time", "unique_frame_count", "real_frame_time", "frame_width", "frame_height"})
		if err != nil {
			return fmt.Errorf("failed to write CSV header: %v", err)
		}
	}

	// Data structures for frame analysis
	type FrameData struct {
		frameNumber      int
		uniqueFrameCount int
		effectiveFPS     float64
		currentFrameTime float64
		realFrameTime    float64
		frameWidth       int
		frameHeight      int
	}

	var frameAnalysisData []FrameData
	var uniqueFrameDurations []int // Duration of each unique frame

	currentFrame := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
	previousFrame := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
	video.SetFrameBuffer(currentFrame.Pix)

	// FIRST PASS: Analyze frame durations
	var frameNumber int
	var uniqueFramesPerSecond []int
	var framePersistenceDurations []float64
	var frameWidthMeasurements []int
	var frameHeightMeasurements []int

	currentSecond := 0
	uniqueFramesInCurrentSecond := 0
	consecutiveDuplicateCount := 0
	totalUniqueFrames := 0
	currentUniqueFrameDuration := 1

	hasFirstFrame := false

	bar := pb.StartNew(video.Frames())

	for video.Read() {
		frameNumber++

		// frame colum will be full of 0s normally, not the worst compromise
		frameWidth := 0
		frameHeight := 0

		// mesure resoltion
		if toggleResdet {
			frameFile, err0 := os.Create("/tmp/frame.png")

			err1 := png.Encode(frameFile, currentFrame)

			out, err2 := exec.Command("resdet", "-v", "1", frameFile.Name()).Output()

			err3 := frameFile.Close()

			err4 := os.Remove(frameFile.Name())

			formattedOutput := strings.Split(string(out), " ")

			frameWidthOut, err5 := strconv.Atoi(formattedOutput[0])

			frameHeightOut, err6 := strconv.Atoi(strings.TrimSuffix(formattedOutput[1], "\n"))

			if err := cmp.Or(err0, err1, err2, err3, err4, err5, err6); err != nil {
				log.Fatal(err)
			}

			frameWidth = frameWidthOut
			frameHeight = frameHeightOut
		}

		frameWidthMeasurements = append(frameWidthMeasurements, frameWidth)
		frameHeightMeasurements = append(frameHeightMeasurements, frameHeight)

		if !hasFirstFrame {
			copy(previousFrame.Pix, currentFrame.Pix)
			hasFirstFrame = true
			uniqueFramesInCurrentSecond = 1
			totalUniqueFrames = 1
			currentUniqueFrameDuration = 1

			// Store data for first frame
			currentTime := float64(frameNumber) / fps
			effectiveFPS := float64(totalUniqueFrames) / currentTime
			actualFrameTimeMs := float64(currentUniqueFrameDuration) * frameTimeMs
			frameAnalysisData = append(frameAnalysisData, FrameData{
				frameNumber:      frameNumber,
				uniqueFrameCount: totalUniqueFrames,
				effectiveFPS:     effectiveFPS,
				currentFrameTime: actualFrameTimeMs,
				realFrameTime:    0, // Will be calculated in second pass
				frameWidth:       frameWidth,
				frameHeight:      frameHeight,
			})
			continue
		}

		isFrameDifferent := false
		pixelDifferences := uint64(0)

		for i := 0; i < len(currentFrame.Pix); i++ {
			if isDiffUInt8WithTolerance(currentFrame.Pix[i], previousFrame.Pix[i], tolerance) {
				pixelDifferences++
			}
		}

		if pixelDifferences > 0 {
			isFrameDifferent = true
		}

		if !isFrameDifferent {
			consecutiveDuplicateCount++
			currentUniqueFrameDuration++
		} else {
			// Record the duration of the previous unique frame
			if totalUniqueFrames > 0 {
				if len(uniqueFrameDurations) < totalUniqueFrames {
					uniqueFrameDurations = append(uniqueFrameDurations, currentUniqueFrameDuration)
				} else {
					uniqueFrameDurations[totalUniqueFrames-1] = currentUniqueFrameDuration
				}
			}

			if consecutiveDuplicateCount > 1 {
				persistenceMs := float64(consecutiveDuplicateCount+1) * frameTimeMs
				framePersistenceDurations = append(framePersistenceDurations, persistenceMs)
				log.Default().Printf("Frame persisted for %.2f ms (%d consecutive duplicates)", persistenceMs, consecutiveDuplicateCount)
			}
			consecutiveDuplicateCount = 0

			uniqueFramesInCurrentSecond++
			totalUniqueFrames++
			copy(previousFrame.Pix, currentFrame.Pix)

			// Start tracking new unique frame
			currentUniqueFrameDuration = 1
		}

		// Store data for EVERY frame
		currentTime := float64(frameNumber) / fps
		effectiveFPS := float64(totalUniqueFrames) / currentTime
		actualFrameTimeMs := float64(currentUniqueFrameDuration) * frameTimeMs
		frameAnalysisData = append(frameAnalysisData, FrameData{
			frameNumber:      frameNumber,
			uniqueFrameCount: totalUniqueFrames,
			effectiveFPS:     effectiveFPS,
			currentFrameTime: actualFrameTimeMs,
			realFrameTime:    0, // Will be calculated in second pass
			frameWidth:       frameWidth,
			frameHeight:      frameHeight,
		})

		if verbose {
			newSecond := int(float64(frameNumber-1) / fps)
			if newSecond > currentSecond {
				uniqueFramesPerSecond = append(uniqueFramesPerSecond, uniqueFramesInCurrentSecond)
				log.Default().Printf("Second %d: %d unique frames", currentSecond+1, uniqueFramesInCurrentSecond)
				currentSecond = newSecond
				uniqueFramesInCurrentSecond = 0
			}
		}

		bar.Increment()
	}

	bar.Finish()

	// Record the final unique frame duration
	if totalUniqueFrames > 0 {
		if len(uniqueFrameDurations) < totalUniqueFrames {
			uniqueFrameDurations = append(uniqueFrameDurations, currentUniqueFrameDuration)
		} else {
			uniqueFrameDurations[totalUniqueFrames-1] = currentUniqueFrameDuration
		}
	}

	// SECOND PASS: Calculate real frame times and write CSV
	if csvWriter != nil {
		for i, frameData := range frameAnalysisData {
			realFrameTimeMs := float64(uniqueFrameDurations[frameData.uniqueFrameCount-1]) * frameTimeMs
			err := csvWriter.Write([]string{
				strconv.Itoa(frameData.frameNumber),
				fmt.Sprintf("%.2f", frameData.effectiveFPS),
				fmt.Sprintf("%.2f", frameData.currentFrameTime),
				strconv.Itoa(frameData.uniqueFrameCount),
				fmt.Sprintf("%.2f", realFrameTimeMs),
				strconv.Itoa(frameData.frameWidth),
				strconv.Itoa(frameData.frameHeight),
			})
			if err != nil {
				log.Default().Printf("Warning: failed to write CSV row %d: %v", i+1, err)
			}
		}
	}

	if consecutiveDuplicateCount > 1 {
		persistenceMs := float64(consecutiveDuplicateCount+1) * frameTimeMs
		framePersistenceDurations = append(framePersistenceDurations, persistenceMs)
		log.Default().Printf("Final frame persisted for %.2f ms (%d consecutive duplicates)", persistenceMs, consecutiveDuplicateCount)
	}

	if uniqueFramesInCurrentSecond > 0 {
		uniqueFramesPerSecond = append(uniqueFramesPerSecond, uniqueFramesInCurrentSecond)
		log.Default().Printf("Second %d: %d unique frames", currentSecond+1, uniqueFramesInCurrentSecond)
	}

	log.Default().Printf("\n=== SUMMARY ===")
	log.Default().Printf("Total frames analyzed: %d", frameNumber)
	log.Default().Printf("Video duration: %.2f seconds", float64(frameNumber)/fps)

	summaryUniqueFrames := 0
	for i, count := range uniqueFramesPerSecond {
		summaryUniqueFrames += count
		log.Default().Printf("Second %d: %d unique frames", i+1, count)
	}

	log.Default().Printf("Total unique frames: %d", summaryUniqueFrames)
	if len(uniqueFramesPerSecond) > 0 {
		log.Default().Printf("Average unique frames per second: %.2f", float64(summaryUniqueFrames)/float64(len(uniqueFramesPerSecond)))
	}

	if len(framePersistenceDurations) > 0 {
		log.Default().Printf("\nFrame persistence durations:")
		totalPersistence := 0.0
		for _, duration := range framePersistenceDurations {
			totalPersistence += duration
		}
		avgPersistence := totalPersistence / float64(len(framePersistenceDurations))
		log.Default().Printf("Average frame persistence: %.2f ms", avgPersistence)
		log.Default().Printf("Number of persistence events: %d", len(framePersistenceDurations))
	} else {
		log.Default().Printf("No frame persistence detected (all frames are unique)")
	}

	if len(frameWidthMeasurements) > 0 && len(frameHeightMeasurements) > 0 {
		sumWidth := 0
		sumHeight := 0

		for _, width := range frameWidthMeasurements {
			sumWidth += width
		}

		if sumWidth != 0 {

			for _, height := range frameHeightMeasurements {
				sumHeight += height
			}

			avgWidth := float64(sumWidth) / float64(len(frameWidthMeasurements))
			avgHeight := float64(sumHeight) / float64(len(frameHeightMeasurements))
			log.Default().Printf("Average Width: %.2f", avgWidth)
			log.Default().Printf("Average Height: %.2f", avgHeight)
		}
	}

	return nil
}
