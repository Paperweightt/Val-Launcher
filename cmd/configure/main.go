package main

import (
	"encoding/json"
	"fmt"
	// "image/color"
	"image"
	"image/color"
	"image/png"
	"log" // TODO: remove
	"os"
	"os/exec"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"path/filepath"
	// "gioui.org/text"
	// "gioui.org/unit"
	// "gioui.org/widget"
	// "gioui.org/widget/material"
)

type Config struct {
	ExeFilepath string `json:"exe_filepath"`
	Changes     []struct {
		Description string   `json:"description"`
		Inputs      []string `json:"inputs"`
		Ouput       string   `json:"ouput"`
	} `json:"changes"`
}

func main() {
	config := parseConfig("./config.json")
	fmt.Println(config)

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

func parseConfig(filepath string) Config {
	var config Config

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("parse: %v", err)
	}

	return config

}

func draw(window *app.Window) error {
	configPath := filepath.Join(exeDir(), "config.json")
	config := parseConfig(configPath)

	var images = layout.List{}
	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				// layout.Rigid(
				// 	func(gtx layout.Context) layout.Dimensions {
				//
				// 		margins := layout.Inset{
				// 			Top:    unit.Dp(15),
				// 			Bottom: unit.Dp(15),
				// 			Right:  unit.Dp(15),
				// 			Left:   unit.Dp(15),
				// 		}
				//
				// 		return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 			const radius = 20
				// 			mp4path := "C:/Users/henry/Projects/Val_Launcher/my_bin/resources/red_dress_1.mp4"
				// 			img := getPng(mp4path)
				// 			size := img.Bounds().Size()
				//
				// 			imgOp := paint.NewImageOp(img)
				// 			imgOp.Add(gtx.Ops)
				//
				// 			clip.RRect{
				// 				Rect: image.Rectangle{Max: size},
				// 				SE:   radius,
				// 				SW:   radius,
				// 				NW:   radius,
				// 				NE:   radius,
				// 			}.Push(gtx.Ops)
				//
				// 			// Paint it to the screen
				// 			paint.PaintOp{}.Add(gtx.Ops)
				//
				// 			return layout.Dimensions{Size: size}
				// 		})
				// 	},
				// ),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return images.Layout(gtx, len(config.Changes[0].Inputs), func(gtx layout.Context, i int) layout.Dimensions {
							const radius = 20
							img := getPng(config.Changes[0].Inputs[i])
							size := img.Bounds().Size()

							imgOp := paint.NewImageOp(img)
							imgOp.Add(gtx.Ops)

							clip.RRect{
								Rect: image.Rectangle{Max: size},
								SE:   radius,
								SW:   radius,
								NW:   radius,
								NE:   radius,
							}.Push(gtx.Ops)

							// Paint it to the screen
							paint.PaintOp{}.Add(gtx.Ops)

							return layout.Dimensions{Size: size}
						})
					},
				),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func getPng(filepath string) image.Image {
	cmd := exec.Command("ffmpeg",
		"-i", filepath,
		"-vf", "select=eq(n\\,0), scale=300:-1",
		"-vframes", "1",
		"-f", "image2pipe",
		"-vcodec", "png",
		"-",
	)

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
}

func exeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)

	// Detect if running from a temp folder (go run)
	if strings.Contains(exeDir, "go-build") {
		fmt.Println("Using development fallback files")
		// fallback to current working dir
		wd, _ := os.Getwd()
		return wd
	}

	return exeDir
}
