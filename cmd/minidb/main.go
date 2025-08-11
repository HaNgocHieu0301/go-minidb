package main

import (
	"fmt"
	"os"

	"github.com/HaNgocHieu0301/go-minidb/pkg/minidb"
)

var version = "0.1.0"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Println("minidbd", version)
		return
	}
	db, err := minidb.Open(minidb.Options{DataDir: "data"})
	if err != nil {
		fmt.Println("failed to open db:", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("minidbd started (sprint0, no server yet)")
}
