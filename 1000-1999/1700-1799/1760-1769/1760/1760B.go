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
        fmt.Fscan(in, &n)
        fmt.Fscan(in, &s)
        maxCh := byte('a')
        for i := 0; i < n; i++ {
            if s[i] > maxCh {
                maxCh = s[i]
            }
        }
        fmt.Println(int(maxCh - 'a' + 1))
    }
}
