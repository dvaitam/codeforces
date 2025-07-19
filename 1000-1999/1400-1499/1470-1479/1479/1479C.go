package main

import (
    "fmt"
)

func main() {
    var l, r int64
    if _, err := fmt.Scan(&l, &r); err != nil {
        return
    }
    type Edge struct {
        u, v int
        cost int64
    }
    var es []Edge
    addEdge := func(u, v int, cost int64) {
        es = append(es, Edge{u, v, cost})
    }
    for i := 2; i <= 21; i++ {
        for j := 1; j < i; j++ {
            cost := int64(1)
            if j-2 > 0 {
                cost = 1 << uint(j-2)
            }
            addEdge(j, i, cost)
        }
    }
    off := l - 1
    l -= off
    r -= off
    const n = 22
    var s int64
    addEdge(1, n, 1)
    s += 1
    for i := 2; i <= 21; i++ {
        if ((r - s) >> uint(i-2)) & 1 == 1 {
            addEdge(i, n, s)
            s += 1 << uint(i-2)
        }
    }
    m := len(es)
    fmt.Println("YES")
    fmt.Printf("%d %d\n", n, m)
    for _, e := range es {
        cost := e.cost
        if e.u == 1 {
            cost += off
        }
        fmt.Printf("%d %d %d\n", e.u, e.v, cost)
    }
}
