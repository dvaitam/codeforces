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
        var n, k int64
        fmt.Fscan(reader, &n, &k)
        div := n - 1
        q := k / div
        r := k % div
        var ans int64
        if r == 0 {
            ans = q*n - 1
        } else {
            ans = q*n + r
        }
        fmt.Fprintln(writer, ans)
    }
}
