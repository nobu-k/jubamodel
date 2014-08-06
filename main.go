package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		Usage()
	}

	commandArgs := os.Args[2:]
	switch os.Args[1] {
	case "info":
		ExecInfo(commandArgs)
	case "rewrite-version":
		ExecRewriteVersion(commandArgs)
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
	fmt.Println("    info             Show information of model files")
	fmt.Println("    rewrite-version  Rewrite a Jubatus version of the given model file")
	os.Exit(1)
}
