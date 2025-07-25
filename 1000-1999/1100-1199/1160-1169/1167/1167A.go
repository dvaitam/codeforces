package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n int
        var s string
        fmt.Fscan(in, &n, &s)
        ok := false
        // find an '8' such that remaining length >= 11
        for i := 0; i < n; i++ {
            if s[i] == '8' && n-i >= 11 {
                ok = true
                break
            }
        }
        if ok {
            fmt.Println("YES")
        } else {
            fmt.Println("NO")
        }
    }
}
