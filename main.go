package main

import (
	"context"
	"image"
	"image/draw"
	"log"
	"os"
	"strconv"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "count-frames",
				Usage: "Count frames",
				Action: func(ctx context.Context, cmd *cli.Command) error {

					return count_video_frames(cmd.Args().First())
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

					first_frame, _ := getImageFromFilePath(cmd.StringArg("frame1"))
					second_frame, _ := getImageFromFilePath(cmd.StringArg("frame2"))

					first_rgba := imageToRGBA(first_frame)
					second_rgba := imageToRGBA(second_frame)
					return compare_frames(first_rgba, second_rgba)
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

					first_frame, _ := getImageFromFilePath(cmd.StringArg("frame1"))
					second_frame, _ := getImageFromFilePath(cmd.StringArg("frame2"))

					first_rgba := imageToRGBA(first_frame)
					second_rgba := imageToRGBA(second_frame)
					return compare_frames_alt(first_rgba, second_rgba)
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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// count_video_frames
// Prints out the total ammount of frames within `video`
//
// Parameters:
//   - video string - path to video file
//
// Returns:
//   - error
func count_video_frames(video string) error {
	log.Default().Print("Trying to open video at: " + video)
	video_file, _ := vidio.NewVideo(video)
	count := 0
	for video_file.Read() {
		count++
	}
	log.Default().Println("Video total frames: " + strconv.Itoa(count))
	return nil
}

func compare_frames(frame1 *image.RGBA, frame2 *image.RGBA) error {
	accumError := int64(0)
	for i := 0; i < len(frame1.Pix); i++ {

		if isDiffUInt8WithTolerance(frame1.Pix[i], frame2.Pix[i], 0) { // Set tolerance to 0
			accumError++
		}
	}
	log.Default().Println("Total differing pixels: " + strconv.FormatInt(accumError, 10))
	return nil
}

func compare_frames_alt(frame1 *image.RGBA, frame2 *image.RGBA) error {
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
	} else {
		return false
	}
}

func isDiffUInt8WithTolerance(x, y uint8, tolerance uint64) bool {
	d := uint64(x) - uint64(y)
	sq := d * d
	if sq > tolerance {
		return true
	} else {
		return false
	}
}

func countUniqueVideoFrames(video_path1 string, video_path2 string, min_diff uint64, use_sq_diff bool) error {
	video1, _ := vidio.NewVideo(video_path1)
	video2, _ := vidio.NewVideo(video_path2)
	video1_frame := image.NewRGBA(image.Rect(0, 0, video1.Width(), video1.Height()))
	video2_frame := image.NewRGBA(image.Rect(0, 0, video2.Width(), video2.Height()))
	video1.SetFrameBuffer(video1_frame.Pix)
	video2.SetFrameBuffer(video2_frame.Pix)
	total_frames := 0
	unique_frames := 0
	for video1.Read() {
		total_frames++
		video2.Read()
		accumError := uint64(0)
		for i := 0; i < len(video1_frame.Pix); i++ {
			if use_sq_diff {
				if isDiffUInt8WithTolerance(video1_frame.Pix[i], video2_frame.Pix[i], min_diff) {
					accumError++
				}
			} else {
				if isDiffUInt8(video1_frame.Pix[i], video2_frame.Pix[i]) {
					accumError++
				}
			}
		}
		if min_diff <= accumError {
			unique_frames++
			log.Default().Println("[" + strconv.Itoa(total_frames) + "]Unique frame")
		} else {
			log.Default().Println("[" + strconv.Itoa(total_frames) + "]Non-unique frame")
		}
	}
	video1.Close()
	video2.Close()
	log.Default().Println(strconv.Itoa(unique_frames) + "/" + strconv.Itoa(total_frames) + " are unique!")
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
