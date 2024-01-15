package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

var (
    helpFlag bool
)

func printUsage() error {
    _, err := fmt.Println("Usage: key-analyzer [options] <name>")
    if err != nil {
        return err
    }
    flag.PrintDefaults()
    return nil
}

func run() error {
    flag.BoolVar(&helpFlag, "help", false, "Show this help")
    flag.Parse()
    args := flag.Args()
    if len(args) == 0 || helpFlag {
        return printUsage()
    }

    filename := args[0]
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    p := NewParser(file)
    prev := &Record{timestamp: time.Unix(math.MaxInt64, 0)}

    for {
        r, err := p.ParseRecord()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        if r.eventType == EventTypeKeyDown {
            if r.timestamp.Sub(prev.timestamp) > 500*time.Millisecond {
                fmt.Printf("\n%s", r.keys)
            } else {
                fmt.Print(r.keys)
            }
            prev = r
        }
    }
}

func main() {
    if err := run(); err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}

