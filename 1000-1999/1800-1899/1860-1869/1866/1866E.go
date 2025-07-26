package main

import (
    "bufio"
    "fmt"
    "os"
    "math"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var N, Q int
    if _, err := fmt.Fscan(in, &N, &Q); err != nil {
        return
    }
    fees := make([]int, Q)
    for i := 0; i < Q; i++ {
        fmt.Fscan(in, &fees[i])
    }

    type Event struct {
        kind int
        x, y int
    }
    events := make([]Event, Q)
    for i := 0; i < Q; i++ {
        var t int
        fmt.Fscan(in, &t)
        if t == 1 {
            var x, y int
            fmt.Fscan(in, &x, &y)
            events[i] = Event{kind: 1, x: x, y: y}
        } else {
            var p int
            fmt.Fscan(in, &p)
            events[i] = Event{kind: 2, x: p}
        }
    }

    floors := []int{1, 1, 1}
    on := []bool{true, true, true}
    var total int64

    for day := 0; day < Q; day++ {
        fee := int64(fees[day])
        ev := events[day]
        if ev.kind == 2 {
            p := ev.x - 1
            on[p] = !on[p]
            continue
        }
        x, y := ev.x, ev.y
        best := -1
        bestDist := math.MaxInt64
        for i := 0; i < 3; i++ {
            if !on[i] {
                continue
            }
            dist := abs(floors[i]-x) + abs(x-y)
            if dist < bestDist {
                bestDist = dist
                best = i
            }
        }
        if best == -1 {
            // should not happen according to problem statement
            continue
        }
        total += fee * int64(bestDist)
        floors[best] = y
    }
    fmt.Println(total)
}

func abs(a int) int {
    if a < 0 {
        return -a
    }
    return a
}

