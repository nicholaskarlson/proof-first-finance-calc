package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	sub := os.Args[1]
	switch sub {
	case "version":
		fmt.Println(version)
		return
	case "demo":
		cmdDemo(os.Args[2:])
		return
	case "serve":
		cmdServe(os.Args[2:])
		return
	case "-h", "--help", "help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", sub)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("proof-first-finance-calc")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  fincalc demo  --out ./out/demo")
	fmt.Println("  fincalc serve --addr :8080")
	fmt.Println("  fincalc version")
	fmt.Println()
}

func cmdDemo(args []string) {
	fs := flag.NewFlagSet("demo", flag.ExitOnError)
	outDir := fs.String("out", "./out/demo", "output directory")
	_ = fs.Parse(args)

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fatal(err)
	}

	// Placeholder output; calculator contract + goldens land in the next PR.
	p := filepath.Join(*outDir, "demo_placeholder.txt")
	if err := os.WriteFile(p, []byte("demo placeholder (calculator contract lands in the next PR)\n"), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("Demo complete. Wrote placeholder to %s\n", p)
}

func cmdServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", ":8080", "listen address")
	_ = fs.Parse(args)

	fmt.Printf("serve placeholder (HTTP API lands in a later PR). Would listen on %s\n", *addr)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
