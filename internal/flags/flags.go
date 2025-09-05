package flags

import (
	"github.com/spf13/pflag"
)

type Flags struct {
	URL string
}

func New() *Flags {
	pflag.Parse()
	url := pflag.Arg(0)

	flags := Flags{
		URL: url,
	}

	return &flags
}
