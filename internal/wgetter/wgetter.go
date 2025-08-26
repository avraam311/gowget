package wgetter

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type WGetter struct {
}

func New() *WGetter {
	return &WGetter{}
}

func (wg *WGetter) WGet() {

}

func getRows(fileName string) []string {
	var data []byte
	var err error

	if fileName == "" {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("error reading from stdin:", err)
			os.Exit(1)
		}
	} else {
		data, err = os.ReadFile(fileName)
		if err != nil {
			fmt.Println("error reading from file:", err)
			os.Exit(1)
		}
	}

	rows := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	return rows
}
