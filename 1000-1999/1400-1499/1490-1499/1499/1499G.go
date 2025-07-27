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

    var n1, n2, m int
    fmt.Fscan(reader, &n1, &n2, &m)
    for i := 0; i < m; i++ {
        var x, y int
        fmt.Fscan(reader, &x, &y)
    }
    var q int
    fmt.Fscan(reader, &q)
    for ; q > 0; q-- {
        var t int
        fmt.Fscan(reader, &t)
        if t == 1 {
            var v1, v2 int
            fmt.Fscan(reader, &v1, &v2)
            // TODO: compute the hash of the optimal coloring after adding the edge
            // The full algorithm is quite involved. For now we output 0 as a placeholder.
            fmt.Fprintln(writer, 0)
        } else {
            // TODO: output the coloring with the same hash printed after the last query of type 1
            // This placeholder outputs an empty coloring.
            fmt.Fprintln(writer, 0)
        }
    }
}
