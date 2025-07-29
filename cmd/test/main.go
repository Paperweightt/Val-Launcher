package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

func main() {

	go func() {
		window := new(app.Window)
		window.Option(app.Title("Val Editor"))

		err := draw(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func draw(window *app.Window) error {
	mp4path := "C:/Users/henry/Projects/Val_Launcher/my_bin/resources/red_dress_1.mp4"
	pngImage := getPng(mp4path)

	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// gtx := app.NewContext(&ops, e)

			drawImage(&ops, pngImage)
		}
	}
}

func drawImage(ops *op.Ops, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterNearest
	imageOp.Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4))).Add(ops)
	paint.PaintOp{}.Add(ops)
}

func getPng(filepath string) image.Image {
	cmd := exec.Command("ffmpeg", "-i", filepath, "-vf", "select=eq(n\\,0)", "-vframes", "1", "-f", "image2pipe", "-vcodec", "png", "-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Failed to get stdout:", err)
		return nil
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start ffmpeg:", err)
		return nil
	}

	img, err := png.Decode(stdout)
	if err != nil {
		fmt.Println("Failed to decode image:", err)
		return nil
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("ffmpeg error:", err)
		return nil
	}

	return img

	// stderr, _ := cmd.StderrPipe()
	// stdout, _ := cmd.StdoutPipe()
	//
	// if err := cmd.Start(); err != nil {
	// 	fmt.Println("Start error:", err)
	// 	return nil
	// }
	//
	// // Print ffmpeg stderr output (shows if the input video is bad or missing codec)
	// go func() {
	// 	io.Copy(os.Stderr, stderr)
	// }()
	//
	// img, err := png.Decode(stdout)
	// if err != nil {
	// 	fmt.Println("Decode error:", err)
	// 	return nil
	// }
	//
	// if err := cmd.Wait(); err != nil {
	// 	fmt.Println("Wait error:", err)
	// 	return nil
	// }
	//
	// return img
}
