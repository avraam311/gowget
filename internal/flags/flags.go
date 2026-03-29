package flags

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type Flags struct {
	URL   string
	Depth int
}

func New() *Flags {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gowget <url> [-d depth]")
		os.Exit(1)
	}

	pflag.Usage = func() {
		fmt.Println("Usage: gowget <url> [-d depth]")
		pflag.PrintDefaults()
	}

	var depth int
	pflag.IntVarP(&depth, "depth", "d", 1, "максимальная глубина рекурсии скачивания")
	pflag.Parse()

	url := pflag.Arg(0)
	if url == "" {
		pflag.Usage()
		os.Exit(1)
	}

	flags := Flags{
		URL:   url,
		Depth: depth,
	}

	return &flags
}
