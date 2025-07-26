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
    deg := make([]int, n+1)
    for i := 0; i < n-1; i++ {
        var u, v int
        fmt.Fscan(in, &u, &v)
        deg[u]++
        deg[v]++
    }
    leaves := 0
    for i := 1; i <= n; i++ {
        if deg[i] == 1 {
            leaves++
        }
    }
    ans := n - 1
    if leaves*(leaves-1)/2 > ans {
        ans = leaves * (leaves - 1) / 2
    }
    fmt.Println(ans)
}

