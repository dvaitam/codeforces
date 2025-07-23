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
    g := make([][]int, n)
    for i := 0; i < m; i++ {
        var u, v int
        fmt.Fscan(in, &u, &v)
        u--
        v--
        g[u] = append(g[u], v)
        g[v] = append(g[v], u)
    }

    dist := make([]int, n)
    q := make([]int, n)
    ans := make([]int, n)

    for s := 0; s < n; s++ {
        for i := 0; i < n; i++ {
            dist[i] = -1
        }
        head, tail := 0, 0
        dist[s] = 0
        q[tail] = s
        tail++
        for head < tail {
            u := q[head]
            head++
            for _, v := range g[u] {
                if dist[v] == -1 {
                    dist[v] = dist[u] + 1
                    q[tail] = v
                    tail++
                }
            }
        }
        maxd := 0
        for i := 0; i < n; i++ {
            if dist[i] > maxd {
                maxd = dist[i]
            }
        }
        ans[s] = maxd
    }

    out := bufio.NewWriter(os.Stdout)
    for i := 0; i < n; i++ {
        if i > 0 {
            out.WriteByte(' ')
        }
        fmt.Fprint(out, ans[i])
    }
    out.WriteByte('\n')
    out.Flush()
}

