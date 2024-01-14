package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

type EventType byte

const (
    EventTypeInvalid EventType = iota
    EventTypeKeyDown
    EventTypeKeyUp
)

const (
  kCGEventKeyDown = 10
  kCGEventKeyUp = 11
)

type Parser struct {
    reader io.Reader
}

type Record struct {
    timestamp time.Time
    eventType EventType
    flags uint64
    keyCode uint16
    keys string
}

func NewParser(r io.Reader) *Parser {
    return &Parser{reader: r}
}

func readTimestamp(r io.Reader) (time.Time, error) {
    var sec int64
    err := binary.Read(r, binary.LittleEndian, &sec)
    if err != nil {
        return time.Time{}, err
    }

    var usec int32
    err = binary.Read(r, binary.LittleEndian, &usec)
    if err != nil {
        return time.Time{}, err
    }

    ts := time.Unix(int64(sec), int64(usec))
    return ts, nil
}

func readEventType(r io.Reader) (EventType, error) {
    var raw uint32
    err := binary.Read(r, binary.LittleEndian, &raw)
    if err != nil {
        return EventTypeInvalid, err
    }
    switch raw {
    case kCGEventKeyDown:
        return EventTypeKeyDown, nil
    case kCGEventKeyUp:
        return EventTypeKeyUp, nil
    default:
        return EventTypeInvalid, fmt.Errorf("invalid event type: %d", raw)
    }
}

func readUTF8Rune(reader io.Reader) (rune, error) {
    var runeBuf bytes.Buffer
    var b byte

    for {
        err := binary.Read(reader, binary.LittleEndian, &b)
        if err != nil {
            return 0, err
        }

        err = binary.Write(&runeBuf, binary.LittleEndian, b)
        if err != nil {
            return 0, err
        }

        if utf8.FullRune(runeBuf.Bytes()) {
            break
        }
    }

    r, _ := utf8.DecodeRune(runeBuf.Bytes())
    return r, nil
}

func readUTF8String(reader io.Reader) (string, error) {
    var sb strings.Builder
    sb.Grow(13)

    for {
        r, err := readUTF8Rune(reader)
        if err != nil {
            return "", err
        }
        if r == '\n' {
            break
        }
        sb.WriteRune(r)
    }
    return sb.String(), nil
}

func (p *Parser) ParseRecord() (*Record, error) {
    // A record contains
    // 1. An i64 (8 bytes) representing the timestamp in seconds
    // 1. An i32 (4 bytes) representing microseconds
    // 2. An i32 (4 bytes) representing the event type
    // 3. An i64 (8 bytes) representing the event flags
    // 4. An i16 (2 bytes) representing the key code
    // 5. A series of UTF-8 characters terminated by a newline
    timestamp, err := readTimestamp(p.reader)
    if err != nil {
        return nil, err
    }

    eventType, err := readEventType(p.reader)
    if err != nil {
        return nil, err
    }

    var flags uint64
    err = binary.Read(p.reader, binary.LittleEndian, &flags)
    if err != nil {
        return nil, err
    }
    
    var keyCode uint16
    err = binary.Read(p.reader, binary.LittleEndian, &keyCode)
    if err != nil {
        return nil, err
    }

    keys, err := readUTF8String(p.reader)
    if err != nil {
        return nil, err
    }

    r := Record{timestamp, eventType, flags, keyCode, keys}

    return &r, nil
}
