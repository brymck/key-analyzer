package main

import (
	"fmt"
	"sort"
	"unicode"
)

type Evaluator struct {
    fingers [10]int32
    total int32
    sfbs int32
    dsfbs int32
    unmapped map[rune]int32
}

type Coordinates struct {
    x byte
    y byte
    finger Finger
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
        'q': {0, 0, LeftPinky},
        'w': {0, 1, LeftRing},
        'e': {0, 2, LeftMiddle},
        'r': {0, 3, LeftIndex},
        't': {0, 4, LeftIndex},
        'y': {0, 5, RightIndex},
        'u': {0, 6, RightIndex},
        'i': {0, 7, RightMiddle},
        'o': {0, 8, RightRing},
        'p': {0, 9, RightPinky},
        'a': {1, 0, LeftPinky},
        's': {1, 1, LeftRing},
        'd': {1, 2, LeftMiddle},
        'f': {1, 3, LeftIndex},
        'g': {1, 4, LeftIndex},
        'h': {1, 5, RightIndex},
        'j': {1, 6, RightIndex},
        'k': {1, 7, RightMiddle},
        'l': {1, 8, RightRing},
        'z': {2, 0, LeftPinky},
        'x': {2, 1, LeftRing},
        'c': {2, 2, LeftMiddle},
        'v': {2, 3, LeftIndex},
        'b': {2, 4, LeftIndex},
        'n': {2, 5, RightIndex},
        'm': {2, 6, RightIndex},
        '␣': {15, 15, LeftThumb},
    }
)

func NewEvaluator() *Evaluator {
    return &Evaluator{unmapped: make(map[rune]int32)}
}

func (e *Evaluator) Evaluate(s string) {
    fmt.Print(s)
    prev := Coordinates{}
    // prevPrev := Coordinates{}
    for i, ch := range s {
        e.total++

        lc := unicode.ToLower(ch)
        if lc != ch {
            // TODO handle shift
            ch = lc
        }
        curr, ok := qwerty[unicode.ToLower(ch)]
        if !ok {
            e.unmapped[ch]++
            continue
        }

        if prev.x == curr.x && prev.y == curr.y {
            continue
        }

        e.fingers[curr.finger]++

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
    fmt.Println("unmapped:")
    e.printUnmapped()
}

type keyValue struct {
    key rune
    value int32
}

func (e *Evaluator) printUnmapped() {
    // Create a slice of key-value pairs
    pairs := make([]keyValue, 0, len(e.unmapped))
    for k, v := range e.unmapped {
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
