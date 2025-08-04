package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	// "gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	// var scroll gesture.Scroll

	// Scroll list state
	list := layout.List{Axis: layout.Vertical}

	// Simulated list items
	items := make([]string, 100)
	for i := range items {
		items[i] = "Item #" + string(rune('A'+i%26))
	}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// scroll.Add(gtx.Ops)
			//
			// if e, ok := scroll.State(); ok {
			// 	list.Position.Offset += e.Scroll.Y
			// }

			layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return list.Layout(gtx, len(items), func(gtx layout.Context, i int) layout.Dimensions {
						// Each list item box
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							paint.FillShape(gtx.Ops, color.NRGBA{R: 200, G: 220, B: 255, A: 255},
								clip.Rect{Max: gtx.Constraints.Max}.Op())
							lbl := material.Body1(th, items[i])
							return lbl.Layout(gtx)
						})
					})
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
