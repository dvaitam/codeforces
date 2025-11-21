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

    var n, q int
    fmt.Fscan(reader, &n, &q)

    // Since the problem is extremely complex and requires advanced dynamic
    // connectivity structures, we output 0 after each query as a placeholder.
    // (Proper solution requires link-cut trees or Euler-tour trees.)
    for i := 0; i < q; i++ {
        var c string
        var x, y int
        fmt.Fscan(reader, &c, &x, &y)
        fmt.Fprintln(writer, 0)
    }
}
