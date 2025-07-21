package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n, t int
    fmt.Fscan(in, &n, &t)
    a := make([]int, n+1)
    for i := 1; i < n; i++ {
        fmt.Fscan(in, &a[i])
    }
    pos := 1
    for pos < t {
        pos += a[pos]
    }
    if pos == t {
        fmt.Println("YES")
    } else {
        fmt.Println("NO")
    }
}
