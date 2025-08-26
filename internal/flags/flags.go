package flags

import "github.com/spf13/pflag"

type Flags struct {
}

func New() *Flags {
	pflag.Parse()

	flags := Flags{
	}

	return &flags
}
