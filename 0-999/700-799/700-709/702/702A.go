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
    fmt.Fscan(reader, &n)
    if n <= 0 {
        fmt.Fprint(writer, 0)
        return
    }

    var prev, curr int
    fmt.Fscan(reader, &prev)
    maxLen, curLen := 1, 1
    for i := 1; i < n; i++ {
        fmt.Fscan(reader, &curr)
        if curr > prev {
            curLen++
        } else {
            if curLen > maxLen {
                maxLen = curLen
            }
            curLen = 1
        }
        prev = curr
    }
    if curLen > maxLen {
        maxLen = curLen
    }

    fmt.Fprint(writer, maxLen)
}
