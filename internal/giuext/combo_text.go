package giuext

import (
	"log"

	g "github.com/AllenDang/giu"
)

type ComboTextWidget struct {
	g.Widget
	selected   int32
	list       []string
	label      string
	onSelected func(string)
}

func (w *ComboTextWidget) Build() {
	preview := ""
	if len(w.list) > 0 {
		preview = w.list[0]
	}
	g.Row(
		g.Label(w.label),
		g.Combo("", preview, w.list[1:], &w.selected).
			OnChange(func() {
				if w.onSelected == nil {
					return
				} else if w.selected < 0 {
					return
				}
				i := w.selected + 1
				value := w.list[i]
				w.list[i] = w.list[0]
				w.list[0] = value
				log.Print(value)
				w.onSelected(value)
			}),
	).Build()
}

func (w *ComboTextWidget) OnSelected(fn func(string)) *ComboTextWidget {
	w.onSelected = fn
	return w
}

func ComboText(label string, list []string) *ComboTextWidget {
	return &ComboTextWidget{
		selected: int32(len(list) - 1),
		label:    label,
		list:     list,
	}
}
