package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	g "github.com/AllenDang/giu"
	"github.com/jessevdk/go-flags"
	"mnemonic.deflinhec.dev/internal/mnemonic"
)

var (
	Version = "0.0.0"
	Build   = "-"
)

var opts struct {
	Profiler func() `long:"pprof" short:"p" description:"Enable profiler"`

	Version func() `long:"version" short:"v" description:"檢視建置版號"`
}

var placeholder struct {
	mnemonic string
	address  string
	message  string
}

var (
	phrases    = 0
	maxphrases = 12
	process    atomic.Bool
	done       = make(chan bool)
	fetcher    = mnemonic.Fetcher()
	parser     = flags.NewParser(&opts, flags.Default)
)

func init() {
	opts.Version = func() {
		fmt.Printf("Version: %v", Version)
		fmt.Printf("\tBuild: %v", Build)
		os.Exit(0)
	}
	opts.Profiler = func() {
		go func() {
			fmt.Println("Profiler listening on port 6060")
			fmt.Println(http.ListenAndServe(":6060", nil))
		}()
	}
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
	placeholder.mnemonic = `range sheriff try enroll deer over ten level bring display stamp *`
	placeholder.address = `TXaMXTQgtdV6iqxtmQ7HNnqzXRoJKXfFAz`
	phrases = mnemonic.Pharse(placeholder.mnemonic).Len()
	process.Store(false)
}

func maxInt(value, max int) int {
	return int(math.Max(float64(value), float64(max)))
}

func loop() {
	g.SingleWindowWithMenuBar().Layout(
		g.PrepareMsgbox(),
		g.MenuBar().Layout(
			g.MenuItem("說明").OnClick(func() {
				g.Msgbox("說明", `複製貼上註記詞，將不確定的助記詞以＊代替．`).
					Buttons(g.MsgboxButtonsOk)
			}),
			g.Align(g.AlignCenter).To(
				g.Row(
					g.Label("助記詞"),
					g.Label("未知"),
					g.Label(fmt.Sprint(maxInt(maxphrases-int(phrases), 0))),
					g.Label("已知"),
					g.Label(fmt.Sprint(phrases)),
				),
			),
		),
		g.Align(g.AlignCenter).To(
			g.Row(
				g.Label("助記詞："),
				g.InputTextMultiline(&placeholder.mnemonic).
					Size(760, 24).
					Flags(g.InputTextFlagsEnterReturnsTrue).
					Flags(g.InputTextFlagsCtrlEnterForNewLine).
					OnChange(func() {
						phrases = mnemonic.Pharse(placeholder.mnemonic).Len()
					}),
			),
		),
		g.Align(g.AlignCenter).To(
			g.Row(
				g.Label("錢包地址："),
				g.InputTextMultiline(&placeholder.address).
					Size(747, 24).
					Flags(g.InputTextFlagsEnterReturnsTrue).
					Flags(g.InputTextFlagsCtrlEnterForNewLine),
			),
		),
		g.Align(g.AlignCenter).To(
			g.Button("匹配").
				OnClick(func() {
					placeholder.address = strings.TrimSpace(placeholder.address)
					if len(placeholder.address) != 34 {
						g.Msgbox("警告", "無效的 TRC20 USDT地址").
							Buttons(g.MsgboxButtonsOk)
						return
					} else if !strings.HasPrefix(placeholder.address, "T") {
						g.Msgbox("警告", "無效的 TRC20 USDT地址").
							Buttons(g.MsgboxButtonsOk)
						return
					} else if mnemonic.Pharse(placeholder.mnemonic).Len() > 12 {
						g.Msgbox("警告", "不支援 12 組以上的助記詞").
							Buttons(g.MsgboxButtonsOk)
						return
					} else {
						defer process.Store(true)
						placeholder.message = "匹配中"
						phrases := mnemonic.Pharse(placeholder.mnemonic)
						go func() {
							defer process.Store(false)
							fetcher.Fetch(placeholder.address, phrases).Wait()
							if found := fetcher.Found(); !found {
								placeholder.message = "無法匹配助記詞"
							} else {
								placeholder.message = fetcher.Result().String()
							}
						}()
					}
				}).Disabled(process.Load()),
		),
		g.Align(g.AlignCenter).To(
			g.Row(
				g.Label("有效組合"),
				g.Label(fmt.Sprint(fetcher.Iterates.Load())),
			),
			g.Label(placeholder.message),
		),
	)
	g.Update()
}

func main() {
	title := "Mnemonic Guess "
	title += fmt.Sprintf("Version:%v Build:%v", Version, Build)
	wnd := g.NewMasterWindow(title, 860, 156, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
