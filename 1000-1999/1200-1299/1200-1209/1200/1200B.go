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
        var n, m, k int
        fmt.Fscan(reader, &n, &m, &k)
        h := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &h[i])
        }
        possible := true
        for i := 0; i < n-1; i++ {
            need := h[i+1] - k
            if need < 0 {
                need = 0
            }
            m += h[i] - need
            if m < 0 {
                possible = false
                break
            }
        }
        if possible {
            fmt.Fprintln(writer, "YES")
        } else {
            fmt.Fprintln(writer, "NO")
        }
    }
}
