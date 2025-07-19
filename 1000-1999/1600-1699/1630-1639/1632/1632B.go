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
    for t > 0 {
        t--
        var n int
        fmt.Fscan(reader, &n)
        m := n - 1
        // find highest power of two <= m
        d := 1
        for d<<1 <= m {
            d <<= 1
        }
        // output from m down to d
        for x := m; x >= d; x-- {
            fmt.Fprint(writer, x, " ")
        }
        // output remaining in desired order
        d--
        fmt.Fprint(writer, 0)
        for d > 0 {
            fmt.Fprint(writer, " ", d)
            d--
        }
        fmt.Fprint(writer, '\n')
    }
}
