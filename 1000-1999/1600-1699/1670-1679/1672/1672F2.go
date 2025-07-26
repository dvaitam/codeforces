package main

import (
    "bufio"
    "fmt"
    "os"
)

// TODO: implement a correct solution. This placeholder just prints WA for all cases.
func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    var t int
    fmt.Fscan(reader, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(reader, &n)
        for i := 0; i < n; i++ {
            var x int
            fmt.Fscan(reader, &x)
        }
        for i := 0; i < n; i++ {
            var x int
            fmt.Fscan(reader, &x)
        }
        fmt.Fprintln(writer, "WA")
    }
}

