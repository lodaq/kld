package main

import (
	"image/color"
	"os"
	"os/exec"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	blue  = color.NRGBA{R: 0x00, G: 0x78, B: 0xD4, A: 0xFF}
	bgCol = color.NRGBA{R: 0xF3, G: 0xF3, B: 0xF3, A: 0xFF}
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("kld"), app.Size(unit.Dp(500), unit.Dp(400)))
		if err := run(w); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	all := getApps()
	var filtered []App

	// SingleLine in the struct literal — guaranteed before first layout.
	editor := widget.Editor{SingleLine: true, Submit: true}

	var list widget.List
	list.Axis = layout.Vertical

	var clicks []widget.Clickable
	sel := -1
	focused := false

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if !focused {
				gtx.Execute(key.FocusCmd{Tag: &editor})
				focused = true
			}

			// Filter.
			q := strings.ToLower(editor.Text())
			filtered = filtered[:0]
			for _, a := range all {
				if strings.Contains(strings.ToLower(a.Name), q) {
					filtered = append(filtered, a)
				}
			}
			for len(clicks) < len(filtered) {
				clicks = append(clicks, widget.Clickable{})
			}
			if sel >= len(filtered) {
				sel = len(filtered) - 1
			}

			// Keyboard nav.
			for {
				ev, ok := gtx.Event(
					key.Filter{Name: key.NameEscape},
					key.Filter{Name: key.NameReturn},
					key.Filter{Name: key.NameDownArrow},
					key.Filter{Name: key.NameUpArrow},
					key.Filter{Name: "J", Required: key.ModCtrl},
					key.Filter{Name: "K", Required: key.ModCtrl},
				)
				if !ok {
					break
				}
				ke, ok := ev.(key.Event)
				if !ok || ke.State != key.Press {
					continue
				}
				switch {
				case ke.Name == key.NameEscape:
					return nil
				case ke.Name == key.NameReturn:
					launch(filtered, max(sel, 0))
					return nil
				case ke.Name == key.NameDownArrow || ke.Name == "J":
					if sel < len(filtered)-1 {
						sel++
					}
				case ke.Name == key.NameUpArrow || ke.Name == "K":
					if sel > 0 {
						sel--
					}
				}
			}

			// Mouse clicks.
			for i := range filtered {
				if clicks[i].Clicked(gtx) {
					launch(filtered, i)
					return nil
				}
			}

			paint.Fill(gtx.Ops, bgCol)

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// Input row.
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						ed := material.Editor(th, &editor, "Search apps...")
						ed.TextSize = unit.Sp(16)
						return ed.Layout(gtx)
					})
				}),
				// List.
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.List(th, &list).Layout(gtx, len(filtered),
						func(gtx layout.Context, i int) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							bg := color.NRGBA{A: 0}
							if i == sel {
								bg = color.NRGBA{R: 0xE3, G: 0xF2, B: 0xFF, A: 0xFF}
							}
							return clicks[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawRow(gtx, th, filtered[i].Name, bg, i == sel)
							})
						})
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func drawRow(gtx layout.Context, th *material.Theme, name string, bg color.NRGBA, selected bool) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body1(th, name)
		lbl.TextSize = unit.Sp(15)
		if selected {
			lbl.Color = blue
		} else {
			lbl.Color = color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
		}
		return lbl.Layout(gtx)
	})
	call := macro.Stop()
	if bg.A > 0 {
		paint.FillShape(gtx.Ops, bg, clip.Rect{Max: dims.Size}.Op())
	}
	call.Add(gtx.Ops)
	return dims
}

func launch(apps []App, i int) {
	if i < len(apps) {
		exec.Command("sh", "-c", apps[i].Exec).Start()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
