package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		Usage()
	}

	switch os.Args[1] {
	case "info":
		ExecInfo(os.Args[2:])
	default:
		Usage()
	}
}

func Usage() {
	fmt.Println("Usage: jubamodel command [options...] [args...]")
	fmt.Println()
	fmt.Println("jubamodel: A manipulation tool for Jubatus model files")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    info    Show information of model files")
	os.Exit(1)
}
