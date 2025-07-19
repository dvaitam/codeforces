package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &a[i])
    }
    var pos, neg []int
    zero := false
    for _, v := range a {
        if v > 0 {
            pos = append(pos, v)
        } else if v < 0 {
            neg = append(neg, v)
        } else {
            zero = true
        }
    }
    sort.Ints(neg)
    out := false
    // print positives
    for _, v := range pos {
        fmt.Fprintf(writer, "%d ", v)
        out = true
    }
    // print pairs of negatives
    for i := 0; i+1 < len(neg); i += 2 {
        fmt.Fprintf(writer, "%d %d ", neg[i], neg[i+1])
        out = true
    }
    if !out {
        if zero {
            fmt.Fprint(writer, "0")
        } else if len(neg) > 0 {
            // only one negative
            fmt.Fprintf(writer, "%d", neg[0])
        }
    }
}
