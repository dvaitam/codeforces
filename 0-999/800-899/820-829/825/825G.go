package main

import (
    "bufio"
    "fmt"
    "os"
)

const maxLog = 20

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, q int
    fmt.Fscan(in, &n, &q)

    adj := make([][]int, n+1)
    for i := 0; i < n-1; i++ {
        var x, y int
        fmt.Fscan(in, &x, &y)
        adj[x] = append(adj[x], y)
        adj[y] = append(adj[y], x)
    }

    // Build parent, depth, and binary lifting tables
    parent := make([][]int, maxLog)
    best := make([][]int, maxLog)
    for i := 0; i < maxLog; i++ {
        parent[i] = make([]int, n+1)
        best[i] = make([]int, n+1)
    }
    depth := make([]int, n+1)

    // iterative DFS to avoid recursion depth issues
    stack := []int{1}
    parent[0][1] = 0
    best[0][1] = 1
    for len(stack) > 0 {
        v := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        for _, to := range adj[v] {
            if to == parent[0][v] {
                continue
            }
            parent[0][to] = v
            depth[to] = depth[v] + 1
            if v < to {
                best[0][to] = v
            } else {
                best[0][to] = to
            }
            stack = append(stack, to)
        }
    }

    for k := 1; k < maxLog; k++ {
        for v := 1; v <= n; v++ {
            p := parent[k-1][v]
            parent[k][v] = parent[k-1][p]
            if p == 0 {
                best[k][v] = best[k-1][v]
            } else {
                b := best[k-1][v]
                if best[k-1][p] < b {
                    b = best[k-1][p]
                }
                best[k][v] = b
            }
        }
    }

    // store black vertices
    black := make([]bool, n+1)
    blacks := make([]int, 0)

    last := 0
    for i := 0; i < q; i++ {
        var t, z int
        fmt.Fscan(in, &t, &z)
        x := (z + last) % n + 1
        if t == 1 {
            if !black[x] {
                black[x] = true
                blacks = append(blacks, x)
            }
        } else {
            ans := n + 1
            for _, b := range blacks {
                v1, v2 := x, b
                cur := n + 1
                if depth[v1] < depth[v2] {
                    v1, v2 = v2, v1
                }
                diff := depth[v1] - depth[v2]
                for k := maxLog - 1; k >= 0; k-- {
                    if diff&(1<<uint(k)) != 0 {
                        if best[k][v1] < cur {
                            cur = best[k][v1]
                        }
                        v1 = parent[k][v1]
                    }
                }
                if v1 == v2 {
                    if v1 < cur {
                        cur = v1
                    }
                } else {
                    for k := maxLog - 1; k >= 0; k-- {
                        if parent[k][v1] != parent[k][v2] {
                            if best[k][v1] < cur {
                                cur = best[k][v1]
                            }
                            if best[k][v2] < cur {
                                cur = best[k][v2]
                            }
                            v1 = parent[k][v1]
                            v2 = parent[k][v2]
                        }
                    }
                    if best[0][v1] < cur {
                        cur = best[0][v1]
                    }
                    if best[0][v2] < cur {
                        cur = best[0][v2]
                    }
                    lca := parent[0][v1]
                    if lca < cur {
                        cur = lca
                    }
                }
                if cur < ans {
                    ans = cur
                }
            }
            fmt.Fprintln(out, ans)
            last = ans
        }
    }
}

