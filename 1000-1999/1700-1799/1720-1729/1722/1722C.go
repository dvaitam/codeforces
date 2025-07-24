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

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(reader, &n)
        a := make([]string, n)
        b := make([]string, n)
        c := make([]string, n)
        counts := make(map[string]int)
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &a[i])
            counts[a[i]]++
        }
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &b[i])
            counts[b[i]]++
        }
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &c[i])
            counts[c[i]]++
        }
        scores := []int{0, 0, 0}
        for i := 0; i < n; i++ {
            if counts[a[i]] == 1 {
                scores[0] += 3
            } else if counts[a[i]] == 2 {
                scores[0] += 1
            }
        }
        for i := 0; i < n; i++ {
            if counts[b[i]] == 1 {
                scores[1] += 3
            } else if counts[b[i]] == 2 {
                scores[1] += 1
            }
        }
        for i := 0; i < n; i++ {
            if counts[c[i]] == 1 {
                scores[2] += 3
            } else if counts[c[i]] == 2 {
                scores[2] += 1
            }
        }
        fmt.Fprintln(writer, scores[0], scores[1], scores[2])
    }
}
