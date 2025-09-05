package app

import (
	"log"

	"github.com/avraam311/gowget/internal/flags"
	"github.com/avraam311/gowget/internal/wgetter"
)

type App struct {
	wgetter *wgetter.WGetter
	flags   *flags.Flags
}

func New(wget *wgetter.WGetter, flags *flags.Flags) *App {
	return &App{
		wgetter: wget,
		flags:   flags,
	}
}

func (a *App) Run() {
	err := a.wgetter.WGet(a.flags.URL)
	if err != nil {
		log.Fatalf("error downloading site: %v", err)
	}
}
