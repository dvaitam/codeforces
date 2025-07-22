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

    var n, k int
    if _, err := fmt.Fscan(reader, &n, &k); err != nil {
        return
    }
    x := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &x[i])
    }

    cnt := 0
    i := 0
    for i < n-1 {
        j := i
        for j+1 < n && x[j+1]-x[i] <= k {
            j++
        }
        if j == i {
            fmt.Fprintln(writer, -1)
            return
        }
        cnt++
        i = j
    }
    fmt.Fprintln(writer, cnt)
}
