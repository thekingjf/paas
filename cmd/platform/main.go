package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thekingjf/paas/internal/tarutil"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: platform deploy <dir>")
		os.Exit(1)
	}

	if os.Args[1] != "deploy" {
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(1)
	}

	buf, err := tarutil.StreamIn(os.Args[2])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("tar bytes:", buf.Len())
}
