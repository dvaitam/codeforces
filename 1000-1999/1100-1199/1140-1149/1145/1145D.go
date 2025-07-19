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

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &a[i])
    }
    m := a[0]
    for i := 1; i < n; i++ {
        if a[i] < m {
            m = a[i]
        }
    }
    // Compute result: 2 + (m XOR a[2])
    res := 2 + (m ^ a[2])
    fmt.Fprintln(writer, res)
}
