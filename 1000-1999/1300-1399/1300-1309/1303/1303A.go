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
    for i := 0; i < t; i++ {
        var s string
        fmt.Fscan(reader, &s)
        l, r := -1, -1
        for j, ch := range s {
            if ch == '1' {
                if l == -1 {
                    l = j
                }
                r = j
            }
        }
        if l < 0 || l == r {
            fmt.Fprintln(writer, 0)
            continue
        }
        cnt := 0
        for j := l; j <= r; j++ {
            if s[j] == '0' {
                cnt++
            }
        }
        fmt.Fprintln(writer, cnt)
    }
}
