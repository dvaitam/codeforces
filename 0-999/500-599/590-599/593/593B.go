package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

type line struct {
    y1 int64
    y2 int64
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    var x1, x2 int64
    if _, err := fmt.Fscan(in, &n, &x1, &x2); err != nil {
        return
    }
    lines := make([]line, n)
    for i := 0; i < n; i++ {
        var k, b int64
        if _, err := fmt.Fscan(in, &k, &b); err != nil {
            return
        }
        lines[i].y1 = k*x1 + b
        lines[i].y2 = k*x2 + b
    }
    sort.Slice(lines, func(i, j int) bool {
        if lines[i].y1 == lines[j].y1 {
            return lines[i].y2 < lines[j].y2
        }
        return lines[i].y1 < lines[j].y1
    })
    for i := 0; i < n-1; i++ {
        if lines[i].y2 > lines[i+1].y2 {
            fmt.Println("Yes")
            return
        }
    }
    fmt.Println("No")
}
