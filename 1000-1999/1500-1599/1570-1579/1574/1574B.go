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

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var a, b, c, m int64
        fmt.Fscan(reader, &a, &b, &c, &m)
        counts := []int64{a, b, c}
        sort.Slice(counts, func(i, j int) bool { return counts[i] < counts[j] })
        x, y, z := counts[2], counts[1], counts[0]
        minPairs := int64(0)
        if x > y+z+1 {
            minPairs = x - (y + z + 1)
        }
        maxPairs := (a - 1) + (b - 1) + (c - 1)
        if m >= minPairs && m <= maxPairs {
            fmt.Fprintln(writer, "YES")
        } else {
            fmt.Fprintln(writer, "NO")
        }
    }
}
