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
    var ans int64
    var ones int64
    for i := 0; i < n; i++ {
        var x int
        fmt.Fscan(reader, &x)
        if x == 1 {
            ones++
        } else {
            ans += ones
        }
    }
    fmt.Fprint(writer, ans)
}
