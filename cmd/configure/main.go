package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log" // TODO: remove
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
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
	var ops op.Ops
	var imagesList = layout.List{
		Axis: layout.Vertical,
	}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// dark gray background
			paint.Fill(gtx.Ops, color.NRGBA{R: 30, G: 30, B: 30, A: 255})

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Flexed(.5, func(gtx layout.Context) layout.Dimensions {
					return imagesTable(gtx, imagesList, 3)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func imagesTable(gtx layout.Context, list layout.List, colums int) layout.Dimensions {
	configPath := filepath.Join(exeDir(), "config.json")
	config := parseConfig(configPath)
	// imagesLen := len(config.Changes[0].Inputs)
	var row = layout.List{}
	const radius = 20

	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return list.Layout(gtx, 6, func(gtx layout.Context, i int) layout.Dimensions {
			return row.Layout(gtx, colums, func(gtx layout.Context, j int) layout.Dimensions {
				// fmt.Println(i, j, i*(colums)+j)
				img := pngFromMp4Dir(config.Changes[0].Inputs[i*(colums)+j])
				size := img.Bounds().Size()

				imgOp := paint.NewImageOp(img)
				imgOp.Add(gtx.Ops)

				defer clip.RRect{
					Rect: image.Rectangle{Max: size},
					SE:   radius,
					SW:   radius,
					NW:   radius,
					NE:   radius,
				}.Push(gtx.Ops).Pop()

				paint.PaintOp{}.Add(gtx.Ops)

				return layout.Dimensions{Size: size}
			})
		})
	})
}

func pngFromMp4Dir(mp4Dir string) image.Image {
	base := filepath.Base(mp4Dir)
	cacheDirBase := base[0:len(base)-3] + "png"
	cacheDir := filepath.Join(exeDir(), "cache", cacheDirBase)

	file, err := os.Open(cacheDir)

	if err != nil {
		cachePngFromMp4Dir(mp4Dir)
		return pngFromMp4Dir(mp4Dir)
	}

	img, err := png.Decode(file)

	if err != nil {
		fmt.Println("Failed to decode image:", err)
		return nil
	}

	return img
}

func cachePngFromMp4Dir(mp4Dir string) error {
	base := filepath.Base(mp4Dir)
	cacheDirBase := base[0:len(base)-3] + "png"
	cacheDir := filepath.Join(exeDir(), "cache", cacheDirBase)

	err := exec.Command("ffmpeg",
		"-i",
		mp4Dir,
		"-vf",
		"select=eq(n\\,0), scale=300:-1",
		"-vframes",
		"1",
		cacheDir,
	).Run()

	if err != nil {
		return err
	}

	return nil
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
