package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: jubamodel /path/to/file [/path/to/file...]")
		os.Exit(1)
	}

	ms, err := Info(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	var js []byte
	if len(ms) == 1 {
		js, err = json.Marshal(ms[0])
	} else { // including 0
		js, err = json.Marshal(ms)
	}
	if err == nil {
		fmt.Print(string(js))
	}
}
