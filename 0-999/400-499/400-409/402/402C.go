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
    fmt.Fscan(reader, &t)
    for ; t > 0; t-- {
        var n, p int
        fmt.Fscan(reader, &n, &p)
        // first cycle edges (distance 1)
        for i := 0; i < n; i++ {
            j := (i + 1) % n
            fmt.Fprintf(writer, "%d %d\n", i+1, j+1)
        }
        // second cycle edges (distance 2)
        for i := 0; i < n; i++ {
            j := (i + 2) % n
            fmt.Fprintf(writer, "%d %d\n", i+1, j+1)
        }
        // extra p edges with increasing distances >=3
        for d := 3; d < n && p > 0; d++ {
            for i := 0; i < n && p > 0; i++ {
                j := (i + d) % n
                if i < j {
                    fmt.Fprintf(writer, "%d %d\n", i+1, j+1)
                    p--
                }
            }
        }
    }
}
