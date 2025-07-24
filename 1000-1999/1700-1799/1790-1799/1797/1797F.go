package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    fmt.Fscan(in, &n)

    maxN := 2*n + 5
    g := make([][]int, maxN)
    for i := 0; i < n-1; i++ {
        var u, v int
        fmt.Fscan(in, &u, &v)
        g[u] = append(g[u], v)
        g[v] = append(g[v], u)
    }

    var m int
    fmt.Fscan(in, &m)
    ops := make([]int, m)
    for i := 0; i < m; i++ {
        fmt.Fscan(in, &ops[i])
    }

    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    for step := 0; step <= m; step++ {
        cur := n + step
        if step > 0 {
            p := ops[step-1]
            v := n + step
            g[p] = append(g[p], v)
            g[v] = append(g[v], p)
        }
        res := countCute(cur, g)
        if step > 0 {
            fmt.Fprint(writer, " ")
        }
        fmt.Fprint(writer, res)
    }
    fmt.Fprintln(writer)
}

func countCute(t int, g [][]int) int {
    res := 0
    for u := 1; u <= t; u++ {
        for v := u + 1; v <= t; v++ {
            path := getPath(u, v, g, t)
            if len(path) == 0 {
                continue
            }
            minIdx, maxIdx := path[0], path[0]
            for _, x := range path[1:] {
                if x < minIdx {
                    minIdx = x
                }
                if x > maxIdx {
                    maxIdx = x
                }
            }
            cond1 := u == minIdx
            cond2 := v == maxIdx
            if cond1 != cond2 {
                res++
            }
        }
    }
    return res
}

func getPath(u, v int, g [][]int, t int) []int {
    parent := make([]int, t+1)
    for i := range parent {
        parent[i] = -1
    }
    q := []int{u}
    parent[u] = 0
    for len(q) > 0 {
        x := q[0]
        q = q[1:]
        if x == v {
            break
        }
        for _, y := range g[x] {
            if y < 1 || y > t {
                continue
            }
            if parent[y] == -1 {
                parent[y] = x
                q = append(q, y)
            }
        }
    }
    if parent[v] == -1 {
        return nil
    }
    var path []int
    cur := v
    for cur != 0 {
        path = append([]int{cur}, path...)
        cur = parent[cur]
    }
    return path
}

