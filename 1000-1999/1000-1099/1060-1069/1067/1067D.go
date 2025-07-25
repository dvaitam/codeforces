package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // Placeholder implementation for problem D.
    // TODO: implement optimal strategy with upgrades.
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    var t int64
    if _, err := fmt.Fscan(reader, &n, &t); err != nil {
        return
    }
    // read quests
    for i := 0; i < n; i++ {
        var a, b float64
        var p float64
        fmt.Fscan(reader, &a, &b, &p)
    }
    // Not implemented
    fmt.Fprintln(writer, "0.0")
}
