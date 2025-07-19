package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

// Node represents a point with coordinates (x, y), time t, and value p.
type Node struct {
    x, y, t int
    p       float64
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    nodes := make([]Node, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &nodes[i].x, &nodes[i].y, &nodes[i].t, &nodes[i].p)
    }
    sort.Slice(nodes, func(i, j int) bool {
        return nodes[i].t < nodes[j].t
    })
    dp := make([]float64, n)
    res := 0.0
    const eps = 1e-9
    for i := 0; i < n; i++ {
        dp[i] = nodes[i].p
        for j := 0; j < i; j++ {
            dt := float64(nodes[i].t - nodes[j].t)
            dx := float64(nodes[i].x - nodes[j].x)
            dy := float64(nodes[i].y - nodes[j].y)
            // reachable if time difference squared >= distance squared
            if dt*dt-dx*dx-dy*dy > -eps {
                if dp[j]+nodes[i].p > dp[i] {
                    dp[i] = dp[j] + nodes[i].p
                }
            }
        }
        if dp[i] > res {
            res = dp[i]
        }
    }
    fmt.Printf("%.9f\n", res)
}
