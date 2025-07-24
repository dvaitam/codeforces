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
        var w, h int64
        fmt.Fscan(reader, &n, &w, &h)
        a := make([]int64, n)
        b := make([]int64, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &a[i])
        }
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &b[i])
        }
        lower := int64(-1 << 60)
        upper := int64(1 << 60)
        for i := 0; i < n; i++ {
            l := (b[i] + h) - (a[i] + w)
            r := (b[i] - h) - (a[i] - w)
            if l > lower {
                lower = l
            }
            if r < upper {
                upper = r
            }
        }
        if lower <= upper {
            fmt.Fprintln(writer, "YES")
        } else {
            fmt.Fprintln(writer, "NO")
        }
    }
}
