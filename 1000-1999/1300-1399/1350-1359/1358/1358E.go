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
    // read initial array of size (n+1)/2
    m := (n + 1) / 2
    a := make([]int64, m)
    for i := 0; i*2 < n; i++ {
        fmt.Fscan(in, &a[i])
    }
    var x int64
    fmt.Fscan(in, &x)
    // compute initial sum
    var s int64
    for i := 0; i < m; i++ {
        s += a[i]
    }
    // if even taking all yields non-positive and x >= 0, impossible
    if s+int64(n/2)*x <= 0 && x >= 0 {
        fmt.Print(-1)
        return
    }
    le := n
    // half offset for calculation: ~n/2 == -((n+1)/2)
    half := m
    // slide window
    for i := 0; le*2 > n && i+le <= n; i++ {
        for le*2 > n && s+x*int64(le+i-half) <= 0 {
            le--
        }
        s -= a[i]
    }
    if le*2 > n {
        fmt.Print(le)
    } else {
        fmt.Print(-1)
    }
}
