package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

const (
    ctrl = '\u2303'
    opt = '\u2325'
    shift = '\u21e7'
    cmd = '\u2318'
)

type Parser struct {
    binaryReader io.Reader
    runeReader *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
    return &Parser{binaryReader: r, runeReader: bufio.NewReader(r)}
}

func (p *Parser) ParseRecord() error {
    var value uint64
    err := binary.Read(p.binaryReader, binary.LittleEndian, &value)
    if err != nil {
        return err
    }

    for {
        r, _, err :=  p.runeReader.ReadRune()
        if err != nil {
            return err
        }
        fmt.Printf("%c", r)
        if r != ctrl && r != opt && r != shift && r != cmd {
            break
        }
    }
    return nil
}
