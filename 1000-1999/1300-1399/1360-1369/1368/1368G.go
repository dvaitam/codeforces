package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n, m int
    fmt.Fscan(in, &n, &m)
    for i := 0; i < n; i++ {
        var s string
        fmt.Fscan(in, &s)
    }
    fmt.Println(0)
}
