package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    t := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Fscan(reader, &t[i])
    }
    parent := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Fscan(reader, &parent[i])
    }
    // depth for nodes in hotel chains
    depth := make([]int, n+1)
    // process each hotel
    for i := 1; i <= n; i++ {
        if t[i] != 1 {
            continue
        }
        if depth[i] != 0 {
            continue
        }
        // collect path from hotel up until known depth or root
        path := make([]int, 0)
        u := i
        for u != 0 && depth[u] == 0 {
            path = append(path, u)
            u = parent[u]
        }
        base := 0
        if u != 0 {
            base = depth[u]
        }
        // assign depths
        for j := len(path) - 1; j >= 0; j-- {
            base++
            depth[path[j]] = base
        }
    }
    // find best hotel
    best := 1
    maxd := 0
    for i := 1; i <= n; i++ {
        if t[i] == 1 && depth[i] > maxd {
            maxd = depth[i]
            best = i
        }
    }
    // reconstruct path from root to hotel
    res := make([]int, 0, maxd)
    u := best
    for u != 0 {
        res = append(res, u)
        u = parent[u]
    }
    // reverse res
    for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
        res[i], res[j] = res[j], res[i]
    }
    // output
    fmt.Fprintln(writer, len(res))
    for i, v := range res {
        if i > 0 {
            writer.WriteByte(' ')
        }
        fmt.Fprint(writer, v)
    }
    fmt.Fprintln(writer)
}
