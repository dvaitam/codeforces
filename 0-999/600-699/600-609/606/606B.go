package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var x, y, x0, y0 int
    if _, err := fmt.Fscan(reader, &x, &y, &x0, &y0); err != nil {
        return
    }
    var s string
    if _, err := fmt.Fscan(reader, &s); err != nil {
        return
    }

    n := len(s)
    visited := make([][]bool, x)
    for i := range visited {
        visited[i] = make([]bool, y)
    }

    cx, cy := x0-1, y0-1
    visited[cx][cy] = true
    visitedCount := 1

    counts := make([]int, n+1)
    counts[0] = 1

    for i, ch := range s {
        switch ch {
        case 'L':
            if cy > 0 {
                cy--
            }
        case 'R':
            if cy+1 < y {
                cy++
            }
        case 'U':
            if cx > 0 {
                cx--
            }
        case 'D':
            if cx+1 < x {
                cx++
            }
        }
        if !visited[cx][cy] {
            visited[cx][cy] = true
            counts[i+1]++
            visitedCount++
        }
    }

    counts[n] += x*y - visitedCount

    for i, v := range counts {
        if i > 0 {
            fmt.Fprint(writer, " ")
        }
        fmt.Fprint(writer, v)
    }
    fmt.Fprintln(writer)
}
