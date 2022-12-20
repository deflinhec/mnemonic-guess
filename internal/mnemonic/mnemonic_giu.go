package mnemonic

import (
	"fmt"

	g "github.com/AllenDang/giu"
)

func (v *FetcherType) Build() {
	g.Align(g.AlignCenter).To(
		g.Row(
			g.Label("有效組合"),
			g.Label(fmt.Sprint(v.Iterates.Load())),
		),
	).Build()
	g.Align(g.AlignCenter).To(
		g.Row(
			g.ProgressBar(float32(v.manager[MATCH].Progress())).
				Size(128, 24).
				Overlayf("MATCH %v", v.manager[MATCH].Jobs()),
			g.ProgressBar(float32(v.manager[EXPAND].Progress())).
				Size(128, 24).
				Overlayf("EXPAND %v", v.manager[EXPAND].Jobs()),
		),
	).Build()
}
