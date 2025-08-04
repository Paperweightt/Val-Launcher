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
	"syscall"
	"unsafe"

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
		windowName := "Val Editor"
		window := new(app.Window)
		window.Option(app.Title("Val Editor"))

		hwnd, err := findWindow("GioWindow", windowName)

		if err != nil {
			fmt.Println("failed to find window", err)
		}

		dll := syscall.NewLazyDLL("../drag_drop/dragdrop.dll")
		setup := dll.NewProc("SetupDragAndDrop")
		setup.Call(uintptr(hwnd))

		err = draw(window)
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

	configPath := filepath.Join(exeDir(), "config.json")
	config := parseConfig(configPath)

	var row = layout.List{}
	const image_round = 20
	const outside_inset = unit.Dp(20)
	const inside_inset = unit.Dp(8)

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Step 1: Get the available width in pixels
			availableWidthPx := gtx.Constraints.Max.X - (int(outside_inset) * 2)
			imageWidthDp := unit.Dp(200) + inside_inset
			imageWidthInt := int(imageWidthDp)

			pictures := len(config.Changes[0].Inputs)
			columns := availableWidthPx / imageWidthInt
			rows := (pictures + columns - 1) / columns // rows = pictures / colums

			fmt.Println("col: ", columns, " row: ", rows)

			// dark gray background
			paint.Fill(gtx.Ops, color.NRGBA{R: 30, G: 30, B: 30, A: 255})

			layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(outside_inset).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return imagesList.Layout(gtx, rows, func(gtx layout.Context, i int) layout.Dimensions {
							if i+1 == rows {
								columns = pictures % columns
							}

							return row.Layout(gtx, columns, func(gtx layout.Context, j int) layout.Dimensions {
								return layout.UniformInset(inside_inset).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									// fmt.Println(i, j, i*(colums)+j)
									img := pngFromMp4Dir(config.Changes[0].Inputs[i*(columns)+j])
									size := img.Bounds().Size()

									imgOp := paint.NewImageOp(img)
									imgOp.Add(gtx.Ops)

									defer clip.RRect{
										Rect: image.Rectangle{Max: size},
										SE:   image_round,
										SW:   image_round,
										NW:   image_round,
										NE:   image_round,
									}.Push(gtx.Ops).Pop()

									paint.PaintOp{}.Add(gtx.Ops)

									return layout.Dimensions{Size: size}
								})
							})
						})
					})
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
		return list.Layout(gtx, 18, func(gtx layout.Context, i int) layout.Dimensions {
			return row.Layout(gtx, colums, func(gtx layout.Context, j int) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
		"select=eq(n\\,0), scale=200:-1",
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

	if strings.Contains(exeDir, "go-build") {
		wd, _ := os.Getwd()
		return wd
	}

	return exeDir
}

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procFindWindowW = user32.NewProc("FindWindowW")
)

func findWindow(className, windowName string) (hwnd syscall.Handle, err error) {
	cn, _ := syscall.UTF16PtrFromString(className)
	wn, _ := syscall.UTF16PtrFromString(windowName)

	ret, _, err := procFindWindowW.Call(
		uintptr(unsafe.Pointer(cn)),
		uintptr(unsafe.Pointer(wn)),
	)
	if ret == 0 {
		return 0, err
	}
	return syscall.Handle(ret), nil
}
