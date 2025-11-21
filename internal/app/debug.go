package app

import "github.com/k0kubun/pp/v3"

var printer *pp.PrettyPrinter

func init() {
	p := pp.New()
	p.SetColorScheme(pp.ColorScheme{
		String:          pp.Magenta,
		StringQuotation: pp.Cyan,
	})
	p.SetColoringEnabled(true)
	printer = p
}

func Dump(msg string, v any) {
	_, _ = printer.Println(msg+" ➡️➡️➡️ ", v)
}
