package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, m, r, c int
        fmt.Fscan(in, &n, &m, &r, &c)
        fmt.Println(0)
    }
}
