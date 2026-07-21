package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	dir := os.Args[2]
	buf, err := tarutil.StreamIn(dir)

	if err != nil {
		log.Fatal(err)
	}

	name := filepath.Base(filepath.Clean(dir))
	addr := os.Getenv("PLATFORM_ADDR")
	if addr == "" {
		addr = "http://localhost:8080"
	}
	resp, err := http.Post(
		addr+"/apps/"+name+"/deploy",
		"application/x-tar",
		buf,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "deploy failed (%d): %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	fmt.Println("deployed: http://" + name + ".localhost:8080")
}
