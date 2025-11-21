package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
)

type fastScanner struct {
    r *bufio.Reader
}

func newFastScanner() *fastScanner {
    return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
    sign, val := 1, 0
    c, err := fs.r.ReadByte()
    for (c < '0' || c > '9') && c != '-' {
        c, err = fs.r.ReadByte()
        if err != nil {
            return 0
        }
    }
    if c == '-' {
        sign = -1
        c, err = fs.r.ReadByte()
        if err != nil {
            return 0
        }
    }
    for c >= '0' && c <= '9' {
        val = val*10 + int(c-'0')
        c, err = fs.r.ReadByte()
        if err != nil {
            break
        }
    }
    return sign * val
}

type edgeInput struct {
    u, v int
    w   int64
}

type edge struct {
    to  int
    rev int
    cap int64
}

type dinic struct {
    n     int
    graph [][]edge
    level []int
    prog  []int
}

func newDinic(n int) *dinic {
    g := make([][]edge, n)
    level := make([]int, n)
    prog := make([]int, n)
    return &dinic{n: n, graph: g, level: level, prog: prog}
}

func (d *dinic) addEdge(u, v int, cap int64) {
    d.graph[u] = append(d.graph[u], edge{to: v, rev: len(d.graph[v]), cap: cap})
    d.graph[v] = append(d.graph[v], edge{to: u, rev: len(d.graph[u]) - 1, cap: 0})
}

func (d *dinic) bfs(s, t int) bool {
    for i := 0; i < d.n; i++ {
        d.level[i] = -1
    }
    queue := []int{s}
    d.level[s] = 0
    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        for _, e := range d.graph[v] {
            if e.cap > 0 && d.level[e.to] < 0 {
                d.level[e.to] = d.level[v] + 1
                queue = append(queue, e.to)
            }
        }
    }
    return d.level[t] >= 0
}

func minInt64(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

func (d *dinic) dfs(v, t int, f int64) int64 {
    if v == t {
        return f
    }
    for ; d.prog[v] < len(d.graph[v]); d.prog[v]++ {
        i := d.prog[v]
        e := d.graph[v][i]
        if e.cap > 0 && d.level[e.to] == d.level[v]+1 {
            pushed := d.dfs(e.to, t, minInt64(f, e.cap))
            if pushed > 0 {
                d.graph[v][i].cap -= pushed
                rev := d.graph[v][i].rev
                d.graph[e.to][rev].cap += pushed
                return pushed
            }
        }
    }
    return 0
}

func (d *dinic) maxFlow(s, t int) int64 {
    flow := int64(0)
    for d.bfs(s, t) {
        for i := range d.prog {
            d.prog[i] = 0
        }
        for {
            pushed := d.dfs(s, t, math.MaxInt64)
            if pushed == 0 {
                break
            }
            flow += pushed
        }
    }
    return flow
}

func gcd64(a, b int64) int64 {
    if a < 0 {
        a = -a
    }
    if b < 0 {
        b = -b
    }
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func solveParity(adj [][]int) []int {
    n := len(adj)
    rows := make([]uint64, n)
    for i := 0; i < n; i++ {
        deg := 0
        for j := 0; j < n; j++ {
            if adj[i][j]&1 != 0 {
                deg ^= 1
            }
        }
        var row uint64
        for j := 0; j < n; j++ {
            coeff := 0
            if j == i {
                coeff = deg
            } else {
                coeff = adj[i][j] & 1
            }
            if coeff != 0 {
                row |= 1 << uint(j)
            }
        }
        if deg != 0 {
            row |= 1 << uint(n)
        }
        rows[i] = row
    }
    where := make([]int, n)
    for i := range where {
        where[i] = -1
    }
    row := 0
    for col := 0; col < n && row < n; col++ {
        sel := -1
        for i := row; i < n; i++ {
            if (rows[i]>>uint(col))&1 != 0 {
                sel = i
                break
            }
        }
        if sel == -1 {
            continue
        }
        rows[row], rows[sel] = rows[sel], rows[row]
        where[col] = row
        for i := 0; i < n; i++ {
            if i != row && ((rows[i]>>uint(col))&1) != 0 {
                rows[i] ^= rows[row]
            }
        }
        row++
    }
    sol := make([]int, n)
    for col := 0; col < n; col++ {
        if where[col] != -1 {
            sol[col] = int((rows[where[col]] >> uint(n)) & 1)
        } else {
            sol[col] = 0
        }
    }
    return sol
}

func main() {
    fs := newFastScanner()
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    t := fs.nextInt()
    for ; t > 0; t-- {
        n := fs.nextInt()
        m := fs.nextInt()
        edges := make([]edgeInput, m)
        adj := make([][]int, n)
        for i := range adj {
            adj[i] = make([]int, n)
        }
        for i := 0; i < m; i++ {
            u := fs.nextInt() - 1
            v := fs.nextInt() - 1
            w := fs.nextInt()
            edges[i] = edgeInput{u: u, v: v, w: int64(w)}
            parity := w & 1
            adj[u][v] = parity
            adj[v][u] = parity
        }
        g := int64(0)
        stop := false
        if n > 1 {
            for i := 0; i < n && !stop; i++ {
                for j := i + 1; j < n; j++ {
                    flowNet := newDinic(n)
                    for _, e := range edges {
                        flowNet.addEdge(e.u, e.v, e.w)
                        flowNet.addEdge(e.v, e.u, e.w)
                    }
                    flow := flowNet.maxFlow(i, j)
                    g = gcd64(g, flow)
                    if g == 1 {
                        stop = true
                        break
                    }
                }
            }
        }
        if n <= 1 || g == 0 || g > 1 {
            fmt.Fprintln(out, 1)
            fmt.Fprintln(out, n)
            if n > 0 {
                for i := 1; i <= n; i++ {
                    if i == n {
                        fmt.Fprintln(out, i)
                    } else {
                        fmt.Fprint(out, i, " ")
                    }
                }
            } else {
                fmt.Fprintln(out)
            }
            continue
        }
        colors := solveParity(adj)
        groups := make([][]int, 2)
        for i := 0; i < n; i++ {
            c := colors[i] & 1
            groups[c] = append(groups[c], i+1)
        }
        fmt.Fprintln(out, 2)
        fmt.Fprintln(out, len(groups[0]))
        if len(groups[0]) > 0 {
            for idx, v := range groups[0] {
                if idx+1 == len(groups[0]) {
                    fmt.Fprintln(out, v)
                } else {
                    fmt.Fprint(out, v, " ")
                }
            }
        } else {
            fmt.Fprintln(out)
        }
        fmt.Fprintln(out, len(groups[1]))
        if len(groups[1]) > 0 {
            for idx, v := range groups[1] {
                if idx+1 == len(groups[1]) {
                    fmt.Fprintln(out, v)
                } else {
                    fmt.Fprint(out, v, " ")
                }
            }
        } else {
            fmt.Fprintln(out)
        }
    }
}
