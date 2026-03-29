package flags

import (
	"github.com/spf13/pflag"
)

type Flags struct {
	URL   string
	Depth int
}

func New() *Flags {
	pflag.Parse()
	url := pflag.Arg(0)

	var depth int
	pflag.IntVarP(&depth, "depth", "d", 1, "максимальная глубина рекурсии скачивания")
	pflag.Parse()

	flags := Flags{
		URL:   url,
		Depth: depth,
	}

	return &flags
}
