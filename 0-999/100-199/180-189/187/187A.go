package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    pos := make([]int, n+1)
    for i := 0; i < n; i++ {
        var x int
        fmt.Fscan(in, &x)
        pos[x] = i
    }
    b := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &b[i])
    }
    // Find longest suffix of b with increasing positions in a
    length := 1
    cur := pos[b[n-1]]
    for i := n - 2; i >= 0; i-- {
        p := pos[b[i]]
        if p < cur {
            length++
            cur = p
        } else {
            break
        }
    }
    fmt.Println(n - length)
}
