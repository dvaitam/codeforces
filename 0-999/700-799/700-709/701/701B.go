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

    var n, m int
    if _, err := fmt.Fscan(reader, &n, &m); err != nil {
        return
    }
    rows := make([]bool, n+1)
    cols := make([]bool, n+1)
    var remRows, remCols int64 = int64(n), int64(n)
    for i := 0; i < m; i++ {
        var x, y int
        fmt.Fscan(reader, &x, &y)
        if !rows[x] {
            rows[x] = true
            remRows--
        }
        if !cols[y] {
            cols[y] = true
            remCols--
        }
        fmt.Fprint(writer, remRows*remCols)
        if i != m-1 {
            fmt.Fprint(writer, " ")
        }
    }
    fmt.Fprintln(writer)
}
