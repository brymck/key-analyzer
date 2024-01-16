package main

import (
	"fmt"
	"sort"
)

type Evaluator struct {
    fingers [10]int32
    total int32
    sfbs int32
    dsfbs int32
    frequencies map[rune]int32
    unmapped map[rune]int32
}

type Coordinates struct {
    x byte
    y byte
    finger Finger
    shift bool
}

type FullCoordinates struct {
    x byte
    y byte
    finger Finger
    shift bool
    ctrl bool
    opt bool
    cmd bool
}

type Finger byte

const (
    LeftPinky Finger = iota
    LeftRing
    LeftMiddle
    LeftIndex
    LeftThumb
    RightThumb
    RightIndex
    RightMiddle
    RightRing
    RightPinky
)

var (
    qwerty = map[rune]Coordinates{
        '`': {0, 0, LeftPinky, false},
        '~': {0, 0, LeftPinky, true},
        '1': {0, 1, LeftPinky, false},
        '!': {0, 1, LeftPinky, true},
        '2': {0, 2, LeftRing, false},
        '@': {0, 2, LeftRing, true},
        '3': {0, 3, LeftMiddle, false},
        '#': {0, 3, LeftMiddle, true},
        '4': {0, 4, LeftIndex, false},
        '$': {0, 4, LeftIndex, true},
        '5': {0, 5, LeftIndex, false},
        '%': {0, 5, LeftIndex, true},
        '6': {0, 6, RightIndex, false},
        '^': {0, 6, RightIndex, true},
        '7': {0, 7, RightIndex, false},
        '&': {0, 7, RightIndex, true},
        '8': {0, 8, RightMiddle, false},
        '*': {0, 8, RightMiddle, true},
        '9': {0, 9, RightRing, false},
        '(': {0, 9, RightRing, true},
        '0': {0, 10, RightIndex, false},
        ')': {0, 10, RightIndex, true},
        '-': {0, 11, RightIndex, false},
        '_': {0, 11, RightIndex, true},
        '=': {0, 12, RightIndex, false},
        '+': {0, 12, RightIndex, true},
        '⌫': {0, 13, RightIndex, false},
        '⇥': {1, 0, LeftPinky, false},
        'q': {1, 1, LeftPinky, false},
        'Q': {1, 1, LeftPinky, true},
        'w': {1, 2, LeftRing, false},
        'W': {1, 2, LeftRing, true},
        'e': {1, 3, LeftMiddle, false},
        'E': {1, 3, LeftMiddle, true},
        'r': {1, 4, LeftIndex, false},
        'R': {1, 4, LeftIndex, true},
        't': {1, 5, LeftIndex, false},
        'T': {1, 5, LeftIndex, true},
        'y': {1, 6, RightIndex, false},
        'Y': {1, 6, RightIndex, true},
        'u': {1, 7, RightIndex, false},
        'U': {1, 7, RightIndex, true},
        'i': {1, 8, RightMiddle, false},
        'I': {1, 8, RightMiddle, true},
        'o': {1, 9, RightRing, false},
        'O': {1, 9, RightRing, true},
        'p': {1, 10, RightPinky, false},
        'P': {1, 10, RightPinky, true},
        '[': {1, 11, RightPinky, false},
        '{': {1, 11, RightPinky, true},
        ']': {1, 12, RightPinky, false},
        '}': {1, 12, RightPinky, true},
        '\\': {1, 13, RightPinky, false},
        '|': {1, 13, RightPinky, true},
        'a': {2, 0, LeftPinky, false},
        'A': {2, 0, LeftPinky, true},
        's': {2, 1, LeftRing, false},
        'S': {2, 1, LeftRing, true},
        'd': {2, 2, LeftMiddle, false},
        'D': {2, 2, LeftMiddle, true},
        'f': {2, 3, LeftIndex, false},
        'F': {2, 3, LeftIndex, true},
        'g': {2, 4, LeftIndex, false},
        'G': {2, 4, LeftIndex, true},
        'h': {2, 5, RightIndex, false},
        'H': {2, 5, RightIndex, true},
        'j': {2, 6, RightIndex, false},
        'J': {2, 6, RightIndex, true},
        'k': {2, 7, RightMiddle, false},
        'K': {2, 7, RightMiddle, true},
        'l': {2, 8, RightRing, false},
        'L': {2, 8, RightRing, true},
        ';': {2, 9, RightRing, false},
        ':': {2, 9, RightRing, true},
        '\'': {2, 10, RightRing, false},
        '"': {2, 10, RightRing, true},
        '↩': {2, 11, RightRing, false},
        'z': {3, 0, LeftPinky, false},
        'Z': {3, 0, LeftPinky, true},
        'x': {3, 1, LeftRing, false},
        'X': {3, 1, LeftRing, true},
        'c': {3, 2, LeftMiddle, false},
        'C': {3, 2, LeftMiddle, true},
        'v': {3, 3, LeftIndex, false},
        'V': {3, 3, LeftIndex, true},
        'b': {3, 4, LeftIndex, false},
        'B': {3, 4, LeftIndex, true},
        'n': {3, 5, RightIndex, false},
        'N': {3, 5, RightIndex, true},
        'm': {3, 6, RightIndex, false},
        'M': {3, 6, RightIndex, true},
        ',': {3, 7, RightMiddle, false},
        '<': {3, 7, RightMiddle, true},
        '.': {3, 8, RightRing, false},
        '>': {3, 8, RightRing, true},
        '/': {3, 9, RightIndex, false},
        '?': {3, 9, RightIndex, true},
        '␣': {15, 15, LeftThumb, false},
    }
)

func NewEvaluator() *Evaluator {
    return &Evaluator{frequencies: make(map[rune]int32), unmapped: make(map[rune]int32)}
}

func (e *Evaluator) Evaluate(s string) {
    var ctrl bool
    var opt bool
    var cmd bool
    var shift bool
    var fingers [10]bool
    fmt.Print(s)
    prev := FullCoordinates{}
    // prevPrev := Coordinates{}
    for i, ch := range s {
        switch ch {
        case '⌃':
            ctrl = true
            continue
        case '⌥':
            opt = true
            continue
        case '⌘':
            cmd = true
            continue
        case '⇧':
            shift = true
            continue
        }
        e.total++

        coords, ok := qwerty[ch]
        if !ok {
            e.unmapped[ch]++
            ctrl = false
            opt = false
            cmd = false
            shift = false
            continue
        }

        if coords.shift {
            shift = true
        }

        curr := FullCoordinates{coords.x, coords.y, coords.finger, shift, ctrl, opt, cmd}
        if prev.x == curr.x && prev.y == curr.y {
            ctrl = false
            opt = false
            cmd = false
            shift = false
            continue
        }

        e.frequencies[ch]++

        fingers[curr.finger] = true
        if curr.ctrl {
            if fingers[LeftPinky] {
                fingers[LeftRing] = true
            } else {
                fingers[LeftPinky] = true
            }
        }
        if curr.shift {
            if fingers[LeftPinky] || curr.finger <= LeftThumb {
                if fingers[RightPinky] {
                    fingers[RightRing] = true
                } else {
                    fingers[RightPinky] = true
                }
            } else {
                fingers[LeftPinky] = true
            }
        }
        if curr.cmd {
            if curr.finger <= LeftThumb {
                fingers[RightThumb] = true
            } else {
                fingers[LeftThumb] = true
            }
        }
        if curr.opt {
            if curr.finger <= LeftThumb {
                fingers[RightThumb] = true
            } else {
                fingers[LeftThumb] = true
            }
        }

        for f, down := range fingers {
            if down {
                e.fingers[f]++
            }
        }

        if i > 0 {
            if curr.finger == prev.finger {
                e.sfbs++
            }
        }
        if i > 1 {
            xDiff := curr.x - prev.x
            if xDiff < 0 {
                xDiff = -xDiff
            }
            yDiff := curr.y - prev.y
            if yDiff < 0 {
                yDiff = -yDiff
            }
            if curr.finger == prev.finger {
                if xDiff > 1 {
                    e.dsfbs++
                } else if yDiff > 1 {
                    e.dsfbs++
                }
            }
        }

        // prevPrev = prev
        prev = curr
        ctrl = false
        opt = false
        cmd = false
        shift = false

        for i := range fingers {
            fingers[i] = false
        }
    }
}

func (e *Evaluator) Print() {
    for i := 0; i < 10; i++ {
        fmt.Printf("%d: %d (%.1f%%)\n", i, e.fingers[i], float32(e.fingers[i])*100.0/float32(e.total))
    }
    fmt.Printf("total: %d\n", e.total)
    fmt.Printf("sfbs: %d (%.1f%%)\n", e.sfbs, float32(e.sfbs)*100.0/float32(e.total))
    fmt.Printf("dsfbs: %d (%.1f%%)\n", e.dsfbs, float32(e.dsfbs)*100.0/float32(e.total))
    fmt.Println()
    fmt.Println("frequencies:")
    e.printDescending(e.frequencies)
    fmt.Println()
    fmt.Println("unmapped:")
    e.printDescending(e.unmapped)
}

type keyValue struct {
    key rune
    value int32
}

func (e *Evaluator) printDescending(m map[rune]int32) {
    // Create a slice of key-value pairs
    pairs := make([]keyValue, 0, len(m))
    for k, v := range m {
        pairs = append(pairs, keyValue{k, v})
    }

    // Sort the slice by value in descending order
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].value > pairs[j].value
    })

    // Print the sorted pairs
    for _, pair := range pairs {
        fmt.Printf("%c: %d\n", pair.key, pair.value)
    }
}
