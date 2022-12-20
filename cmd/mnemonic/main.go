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
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	gext "mnemonic.deflinhec.dev/internal/giuext"
	"mnemonic.deflinhec.dev/internal/mnemonic"
	_ "mnemonic.deflinhec.dev/internal/translations"
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
	translator *message.Printer
	fetcher    = mnemonic.Fetcher()
	languages  = []string{"en-GB", "zh-TW", "zh-CN"}
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
	onTranslate(languages[0])
	process.Store(false)
}

func onTranslate(lang string) {
	l := language.MustParse(lang)
	translator = message.NewPrinter(l)
	g.Update()
}

func maxInt(value, max int) int {
	return int(math.Max(float64(value), float64(max)))
}

func loop() {
	g.SingleWindowWithMenuBar().Layout(
		g.PrepareMsgbox(),
		g.MenuBar().Layout(
			g.Menu(translator.Sprintf("設定")).Layout(
				gext.ComboText(translator.Sprintf("語系"), languages).
					OnSelected(onTranslate),
			),
			g.MenuItem(translator.Sprintf("說明")).OnClick(func() {
				g.Msgbox(translator.Sprintf("說明"),
					translator.Sprintf(`複製貼上註記詞，將不確定的助記詞以＊代替．`)).
					Buttons(g.MsgboxButtonsOk)
			}),
			g.Align(g.AlignCenter).To(

				g.Row(
					g.Label(translator.Sprintf("助記詞")),
					g.Label(translator.Sprintf("未知")),
					g.Label(fmt.Sprint(maxInt(maxphrases-int(phrases), 0))),
					g.Label(translator.Sprintf("已知")),
					g.Label(fmt.Sprint(phrases)),
				),
			),
		),
		g.Align(g.AlignCenter).To(
			g.Row(
				g.Label(translator.Sprintf("助記詞：")),
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
				g.Label(translator.Sprintf("錢包地址：")),
				g.InputTextMultiline(&placeholder.address).
					Size(747, 24).
					Flags(g.InputTextFlagsEnterReturnsTrue).
					Flags(g.InputTextFlagsCtrlEnterForNewLine),
			),
		),
		g.Align(g.AlignCenter).To(
			g.Button(translator.Sprintf("匹配")).
				OnClick(func() {
					placeholder.address = strings.TrimSpace(placeholder.address)
					if len(placeholder.address) != 34 {
						g.Msgbox(translator.Sprintf("警告"),
							translator.Sprintf("無效的 TRC20 USDT地址")).
							Buttons(g.MsgboxButtonsOk)
						return
					} else if !strings.HasPrefix(placeholder.address, "T") {
						g.Msgbox(translator.Sprintf("警告"),
							translator.Sprintf("無效的 TRC20 USDT地址")).
							Buttons(g.MsgboxButtonsOk)
						return
					} else if mnemonic.Pharse(placeholder.mnemonic).Len() > 12 {
						g.Msgbox(translator.Sprintf("警告"),
							translator.Sprintf("不支援 12 組以上的助記詞")).
							Buttons(g.MsgboxButtonsOk)
						return
					} else {
						defer process.Store(true)
						placeholder.message = translator.Sprintf("匹配中")
						phrases := mnemonic.Pharse(placeholder.mnemonic)
						go func() {
							defer process.Store(false)
							fetcher.Fetch(placeholder.address, phrases).Wait()
							if found := fetcher.Found(); !found {
								placeholder.message = translator.Sprintf("無法匹配助記詞")
							} else {
								placeholder.message = fetcher.Result().String()
							}
						}()
					}
				}).Disabled(process.Load()),
		),
		g.Align(g.AlignCenter).To(
			g.Row(
				g.Label(translator.Sprintf("有效組合")),
				g.Label(fmt.Sprint(fetcher.Iterates.Load())),
			),
		),
		g.Align(g.AlignCenter).To(
			g.Row(
				g.ProgressBar(float32(fetcher.Worker(mnemonic.MATCH).Progress())).
					Size(128, 24).
					Overlayf("MATCH %v", fetcher.Worker(mnemonic.MATCH).Jobs()),
				g.ProgressBar(float32(fetcher.Worker(mnemonic.EXPAND).Progress())).
					Size(128, 24).
					Overlayf("EXPAND %v", fetcher.Worker(mnemonic.EXPAND).Jobs()),
			),
		),
		g.Align(g.AlignCenter).To(
			g.Label(placeholder.message),
		),
	)
	g.Update()
}

func main() {
	title := "Mnemonic Guess "
	title += fmt.Sprintf("Version:%v Build:%v", Version, Build)
	wnd := g.NewMasterWindow(title, 860, 180, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
