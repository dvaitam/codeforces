package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

const INF = 1000000000

type Edge struct { v, c int }

var (
    n, llimit, rlimit int
    adj                [][]Edge
    vis                []bool
    s, f               []int
    to                 []int
    w                  []int
    px, py             []int
    ff, q              []int
    xx, yy, mid, ans   int
    tot, maxL          int
)

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func root(u, fa, size int) int {
    ret := -1
    s[u] = 1
    f[u] = 0
    for _, e := range adj[u] {
        v := e.v
        if v == fa || vis[v] {
            continue
        }
        x := root(v, u, size)
        s[u] += s[v]
        if f[u] < s[v] {
            f[u] = s[v]
        }
        if ret < 0 || f[x] < f[ret] {
            ret = x
        }
    }
    if f[u] < size-s[u] {
        f[u] = size - s[u]
    }
    if ret < 0 || f[u] < f[ret] {
        ret = u
    }
    return ret
}

func dfs(u, fa, dep, sum int) {
    if sum > f[dep] {
        f[dep] = sum
        py[dep] = u
    }
    for _, e := range adj[u] {
        v := e.v
        if v == fa || vis[v] {
            continue
        }
        if e.c >= mid {
            dfs(v, u, dep+1, sum+1)
        } else {
            dfs(v, u, dep+1, sum-1)
        }
    }
}

func check(u int) bool {
    px[0] = u
    size0 := s[to[tot-1]]
    for i := 1; i <= size0; i++ {
        ff[i] = -INF
    }
    k := 0
    for i := 0; i < tot; i++ {
        v := to[i]
        // clear f for depths
        for j := 1; j <= s[v]; j++ {
            f[j] = -INF
        }
        if w[v] >= mid {
            dfs(v, u, 1, 1)
        } else {
            dfs(v, u, 1, -1)
        }
        qh, qt := 0, -1
        var j int
        for j = 1; j <= s[v]; j++ {
            if f[j] == -INF {
                break
            }
            for k >= 0 && j+k >= llimit {
                for qh <= qt && ff[q[qt]] < ff[k] {
                    qt--
                }
                qt++
                q[qt] = k
                k--
            }
            for qh <= qt && q[qh]+j > rlimit {
                qh++
            }
            if qh <= qt && ff[q[qh]]+f[j] >= 0 {
                xx = px[q[qh]]
                yy = py[j]
                ans = mid
                return true
            }
        }
        for m := 1; m < j; m++ {
            if f[m] > ff[m] {
                ff[m] = f[m]
                px[m] = py[m]
            }
        }
        k = j - 1
    }
    return false
}

func gao(u, size int) {
    tot = 0
    for _, e := range adj[u] {
        v := e.v
        if vis[v] {
            continue
        }
        if s[v] > s[u] {
            s[v] = size - s[u]
        }
        to[tot] = v
        w[v] = e.c
        tot++
    }
    if tot == 0 {
        return
    }
    sort.Slice(to[:tot], func(i, j int) bool {
        return s[to[i]] < s[to[j]]
    })
    l, r := ans+1, maxL
    mid = l
    if !check(u) {
        return
    }
    for l <= r {
        mid = (l + r) >> 1
        if check(u) {
            l = mid + 1
        } else {
            r = mid - 1
        }
    }
}

func solve(u, size int) {
    x := root(u, 0, size)
    vis[x] = true
    gao(x, size)
    for _, e := range adj[x] {
        v := e.v
        if vis[v] {
            continue
        }
        solve(v, s[v])
    }
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    fmt.Fscan(reader, &n, &llimit, &rlimit)
    adj = make([][]Edge, n+1)
    vis = make([]bool, n+1)
    s = make([]int, n+1)
    f = make([]int, n+1)
    to = make([]int, n+1)
    w = make([]int, n+1)
    px = make([]int, n+1)
    py = make([]int, n+1)
    ff = make([]int, n+1)
    q = make([]int, n+1)
    maxL = 0
    for i := 1; i < n; i++ {
        var a, b, c int
        fmt.Fscan(reader, &a, &b, &c)
        adj[a] = append(adj[a], Edge{b, c})
        adj[b] = append(adj[b], Edge{a, c})
        if c > maxL {
            maxL = c
        }
    }
    ans = -1
    solve(1, n)
    fmt.Fprintln(writer, xx, yy)
}
