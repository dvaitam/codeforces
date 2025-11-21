package main

import (
    "bufio"
    "fmt"
    "os"
)

func max64(a, b int64) int64 {
    if a > b {
        return a
    }
    return b
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        var c int64
        fmt.Fscan(in, &n, &c)
        a := make([]int64, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }
        adj := make([][]int, n)
        for i := 0; i < n-1; i++ {
            var u, v int
            fmt.Fscan(in, &u, &v)
            u--
            v--
            adj[u] = append(adj[u], v)
            adj[v] = append(adj[v], u)
        }
        parent := make([]int, n)
        for i := range parent {
            parent[i] = -1
        }
        order := make([]int, 0, n)
        stack := []int{0}
        parent[0] = 0
        for len(stack) > 0 {
            v := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            order = append(order, v)
            for _, to := range adj[v] {
                if parent[to] == -1 {
                    parent[to] = v
                    stack = append(stack, to)
                }
            }
        }
        dp0 := make([]int64, n)
        dp1 := make([]int64, n)
        twoC := 2 * c
        for i := len(order) - 1; i >= 0; i-- {
            v := order[i]
            sum0 := int64(0)
            sum1 := a[v]
            for _, to := range adj[v] {
                if to == parent[v] {
                    continue
                }
                sum0 += max64(dp0[to], dp1[to])
                sum1 += max64(dp0[to], dp1[to]-twoC)
            }
            dp0[v] = sum0
            dp1[v] = sum1
        }
        ans := max64(dp0[0], dp1[0])
        fmt.Fprintln(out, ans)
    }
}
