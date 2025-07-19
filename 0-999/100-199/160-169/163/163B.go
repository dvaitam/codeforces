package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

// Lem represents a lemming with weight, speed, and original index
type Lem struct {
    weight int
    speed  int
    idx    int
}

var (
    n, k, h int
    lems    []Lem
    ans     []int
)

// check determines if it's possible to select k lemmings
// such that each can climb to its assigned ledge within time x.
func check(x float64) bool {
    t := 0
    for i := 0; i < n; i++ {
        // ledge height is (t+1) * h
        if float64((t+1)*h) <= float64(lems[i].speed)*x {
            ans[t] = lems[i].idx
            t++
            if t >= k {
                return true
            }
        }
    }
    return false
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    fmt.Fscan(reader, &n, &k, &h)
    lems = make([]Lem, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &lems[i].weight)
    }
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &lems[i].speed)
    }
    for i := 0; i < n; i++ {
        lems[i].idx = i + 1
    }
    sort.Slice(lems, func(i, j int) bool {
        if lems[i].weight != lems[j].weight {
            return lems[i].weight < lems[j].weight
        }
        return lems[i].speed < lems[j].speed
    })
    ans = make([]int, k)

    // binary search on time
    l, r := 0.0, 1e9
    for it := 0; it < 100; it++ {
        m := (l + r) / 2
        if check(m) {
            r = m
        } else {
            l = m
        }
    }
    // final assignment
    check(r)
    for i := 0; i < k; i++ {
        fmt.Fprint(writer, ans[i])
        if i+1 < k {
            fmt.Fprint(writer, " ")
        }
    }
    fmt.Fprintln(writer)
}
