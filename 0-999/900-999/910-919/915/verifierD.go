package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

type edge struct{ u, v int }

func checkD(n int, adj [][]int, skipU, skipV int) bool {
    indeg := make([]int, n+1)
    for u := 1; u <= n; u++ {
        for _, v := range adj[u] {
            if u == skipU && v == skipV {
                continue
            }
            indeg[v]++
        }
    }
    q := make([]int, 0, n)
    for i := 1; i <= n; i++ {
        if indeg[i] == 0 {
            q = append(q, i)
        }
    }
    count := 0
    for head := 0; head < len(q); head++ {
        v := q[head]
        count++
        for _, to := range adj[v] {
            if v == skipU && to == skipV {
                continue
            }
            indeg[to]--
            if indeg[to] == 0 {
                q = append(q, to)
            }
        }
    }
    return count == n
}

func findCycle(n int, adj [][]int) ([]edge, bool) {
    color := make([]int, n+1)
    stack := make([]int, 0, n)
    pos := make([]int, n+1)
    var cycle []edge
    var found bool
    var dfs func(int)
    dfs = func(v int) {
        color[v] = 1
        pos[v] = len(stack)
        stack = append(stack, v)
        for _, to := range adj[v] {
            if found {
                return
            }
            if color[to] == 0 {
                dfs(to)
            } else if color[to] == 1 {
                found = true
                idx := pos[to]
                nodes := append([]int{}, stack[idx:]...)
                nodes = append(nodes, to)
                for i := 0; i < len(nodes)-1; i++ {
                    cycle = append(cycle, edge{nodes[i], nodes[i+1]})
                }
                return
            }
        }
        stack = stack[:len(stack)-1]
        color[v] = 2
    }
    for i := 1; i <= n && !found; i++ {
        if color[i] == 0 {
            dfs(i)
        }
    }
    return cycle, found
}

func solveD(n int, edges []edge) string {
    adj := make([][]int, n+1)
    for _, e := range edges {
        adj[e.u] = append(adj[e.u], e.v)
    }
    if checkD(n, adj, -1, -1) {
        return "YES"
    }
    cyc, ok := findCycle(n, adj)
    if !ok {
        return "NO"
    }
    for _, e := range cyc {
        if checkD(n, adj, e.u, e.v) {
            return "YES"
        }
    }
    return "NO"
}

func generateD(rng *rand.Rand) (string, string) {
    n := rng.Intn(7) + 2
    maxM := n * (n - 1)
    m := rng.Intn(min(maxM, 10)) + 1
    edges := make([]edge, 0, m)
    seen := make(map[[2]int]bool)
    for len(edges) < m {
        u := rng.Intn(n) + 1
        v := rng.Intn(n) + 1
        if u == v {
            continue
        }
        key := [2]int{u, v}
        if seen[key] {
            continue
        }
        seen[key] = true
        edges = append(edges, edge{u, v})
    }
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
    for _, e := range edges {
        sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
    }
    return sb.String(), solveD(n, edges)
}

func min(a, b int) int { if a < b { return a }; return b }

func runCase(bin, input, exp string) error {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    if got != exp {
        return fmt.Errorf("expected %s got %s", exp, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    for i := 0; i < 100; i++ {
        in, exp := generateD(rng)
        if err := runCase(bin, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

