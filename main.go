package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
    helpFlag bool
)

func printUsage() {
    fmt.Println("Usage: key-analyzer [options] <name>")
    flag.PrintDefaults()
}

func run() error {
    flag.BoolVar(&helpFlag, "help", false, "Show this help")
    flag.Parse()
    args := flag.Args()
    if len(args) == 0 || helpFlag {
        printUsage()
        return nil
    }

    filename := args[0]
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    p := NewParser(file)

    for {
        err := p.ParseRecord()
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
        // Sleep for 1 second
        time.Sleep(1 * time.Second)
    }

    return nil
}

func main() {
    if err := run(); err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}

