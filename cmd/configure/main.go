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

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procFindWindowW = user32.NewProc("FindWindowW")
	windowName      = "ValEditor"
	hwnd            = 0
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title(windowName))

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

			layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return imagesTable(gtx, imagesList)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}

func imagesTable(gtx layout.Context, list layout.List) layout.Dimensions {
	configPath := filepath.Join(exeDir(), "config.json")
	config := parseConfig(configPath)
	// imagesLen := len(config.Changes[0].Inputs)
	var row = layout.List{}
	const radius = 20
	const image_round = 20
	const outside_inset = unit.Dp(20)
	const inside_inset = unit.Dp(8)

	availableWidthPx := gtx.Constraints.Max.X - (int(outside_inset) * 2)
	imageWidthDp := unit.Dp(200) + inside_inset
	imageWidthInt := int(imageWidthDp)

	pictures := len(config.Changes[0].Inputs)
	columns := availableWidthPx / imageWidthInt
	rows := (pictures + columns - 1) / columns // rows = pictures / colums

	return list.Layout(gtx, columns, func(gtx layout.Context, j int) layout.Dimensions {
		return layout.UniformInset(outside_inset).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return row.Layout(gtx, rows, func(gtx layout.Context, i int) layout.Dimensions {
				if i+1 == rows {
					columns = pictures % columns
				}

				return row.Layout(gtx, columns, func(gtx layout.Context, j int) layout.Dimensions {
					return layout.UniformInset(inside_inset).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
	})
}

func addDragDropDll() {
	if hwnd != 0 {
		return
	}

	localHwnd, err := findWindow("GioWindow", windowName)

	if localHwnd == 0 {
		fmt.Println("failed to find window: ", err)
	}

	fmt.Println("hwnd:", localHwnd)

	dll := syscall.NewLazyDLL("../drag_drop/dragdrop.dll")

	err = dll.Load()
	if err != nil {
		log.Fatal("DLL failed to load:", err)
	}

	setup := dll.NewProc("DummyFunc")

	ret, _, callErr := setup.Call()
	if callErr != syscall.Errno(0) {
		log.Fatal("Call failed:", callErr)
	}

	fmt.Println("Return value from DummyFunc:", ret) // should print 1

	hwnd = int(localHwnd)
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
