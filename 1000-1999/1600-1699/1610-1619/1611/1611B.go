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

    var T int
    fmt.Fscan(reader, &T)
    for ; T > 0; T-- {
        var a, b int
        fmt.Fscan(reader, &a, &b)
        if a > b {
            a, b = b, a
        }
        total := (a + b) / 4
        ans := a
        if b < ans {
            ans = b
        }
        if total < ans {
            ans = total
        }
        fmt.Fprintln(writer, ans)
    }
}
