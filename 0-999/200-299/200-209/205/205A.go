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
    var minTime int64 = 1<<63 - 1
    idx := -1
    count := 0
    for i := 1; i <= n; i++ {
        var t int64
        fmt.Fscan(reader, &t)
        if t < minTime {
            minTime = t
            idx = i
            count = 1
        } else if t == minTime {
            count++
        }
    }
    if count > 1 {
        fmt.Fprint(writer, "Still Rozdil")
    } else {
        fmt.Fprint(writer, idx)
    }
}
