package main

import (
    "bufio"
    "fmt"
    "os"
)

// This is a placeholder implementation for problem E.
// The expected solution involves a monotonic stack and dynamic
// programming with a segment tree, but that algorithm is omitted.
// The current version only reads the input and outputs 0.
func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    fmt.Fscan(reader, &n)
    h := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &h[i])
    }
    b := make([]int64, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &b[i])
    }

    // TODO: implement the correct algorithm.
    fmt.Fprintln(writer, 0)
}

