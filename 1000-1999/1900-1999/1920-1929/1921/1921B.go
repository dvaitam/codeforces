package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        var s, f string
        fmt.Fscan(in, &s)
        fmt.Fscan(in, &f)
        cnt01, cnt10 := 0, 0
        for i := 0; i < n; i++ {
            if s[i] != f[i] {
                if s[i] == '0' {
                    cnt01++
                } else {
                    cnt10++
                }
            }
        }
        if cnt01 > cnt10 {
            fmt.Fprintln(out, cnt01)
        } else {
            fmt.Fprintln(out, cnt10)
        }
    }
}
