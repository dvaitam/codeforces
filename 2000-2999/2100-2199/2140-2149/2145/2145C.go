package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }

    for ; t > 0; t-- {
        var n int
        var s string
        fmt.Fscan(in, &n)
        fmt.Fscan(in, &s)

        diff := 0
        for _, ch := range s {
            if ch == 'a' {
                diff++
            } else {
                diff--
            }
        }

        if diff == 0 {
            fmt.Fprintln(out, 0)
            continue
        }

        prefix := 0
        earliest := map[int]int{0: 0}
        best := n + 1

        for i, ch := range s {
            if ch == 'a' {
                prefix++
            } else {
                prefix--
            }

            if idx, ok := earliest[prefix-diff]; ok {
                if length := i + 1 - idx; length < best {
                    best = length
                }
            }

            if _, ok := earliest[prefix]; !ok {
                earliest[prefix] = i + 1
            }
        }

        if best == n {
            fmt.Fprintln(out, -1)
        } else {
            fmt.Fprintln(out, best)
        }
    }
}

