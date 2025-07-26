package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log" // TODO: remove
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
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
		err := run(window)
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

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Define an large label with an appropriate text:
			title := material.H1(theme, "Hello, Gio")

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon

			// Change the position of the label.
			title.Alignment = text.Middle

			// Draw the label to the graphics context.
			title.Layout(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
