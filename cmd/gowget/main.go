package main

import (
	"github.com/avraam311/gowget/cmd/app"
	"github.com/avraam311/gowget/internal/flags"
	"github.com/avraam311/gowget/internal/wgetter"
)

func main() {
	flags := flags.New()
	wgetter := wgetter.New()
	app := app.New(wgetter, flags)
	app.Run()
}
