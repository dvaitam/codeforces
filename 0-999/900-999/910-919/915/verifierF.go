package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

// DSU structure used in original solution
type DSU struct {
    parent []int
    size   []int
}

func NewDSU(n int) *DSU {
    d := &DSU{parent: make([]int, n), size: make([]int, n)}
    for i := 0; i < n; i++ {
        d.parent[i] = i
        d.size[i] = 1
    }
    return d
}

func (d *DSU) find(x int) int {
    for d.parent[x] != x {
        d.parent[x] = d.parent[d.parent[x]]
        x = d.parent[x]
    }
    return x
}

func (d *DSU) union(a, b int) int64 {
    ra := d.find(a)
    rb := d.find(b)
    if ra == rb {
        return 0
    }
    if d.size[ra] < d.size[rb] {
        ra, rb = rb, ra
    }
    pairs := int64(d.size[ra]) * int64(d.size[rb])
    d.parent[rb] = ra
    d.size[ra] += d.size[rb]
    return pairs
}

func solveF(n int, vals []int, edges [][2]int) string {
    g := make([][]int, n)
    for _, e := range edges {
        x, y := e[0], e[1]
        g[x] = append(g[x], y)
        g[y] = append(g[y], x)
    }
    type node struct{ val, idx int }
    order := make([]node, n)
    for i := 0; i < n; i++ {
        order[i] = node{val: vals[i], idx: i}
    }
    sort.Slice(order, func(i, j int) bool { return order[i].val < order[j].val })
    dsu := NewDSU(n)
    active := make([]bool, n)
    var maxSum int64
    for _, p := range order {
        v := p.idx
        active[v] = true
        for _, to := range g[v] {
            if active[to] {
                maxSum += int64(p.val) * dsu.union(v, to)
            }
        }
    }
    sort.Slice(order, func(i, j int) bool { return order[i].val > order[j].val })
    dsu = NewDSU(n)
    for i := range active { active[i] = false }
    var minSum int64
    for _, p := range order {
        v := p.idx
        active[v] = true
        for _, to := range g[v] {
            if active[to] {
                minSum += int64(p.val) * dsu.union(v, to)
            }
        }
    }
    return fmt.Sprintf("%d", maxSum-minSum)
}

func generateTree(rng *rand.Rand, n int) [][2]int {
    edges := make([][2]int, 0, n-1)
    for i := 1; i < n; i++ {
        j := rng.Intn(i)
        edges = append(edges, [2]int{i, j})
    }
    return edges
}

func generateF(rng *rand.Rand) (string, string) {
    n := rng.Intn(15) + 1
    vals := make([]int, n)
    for i := 0; i < n; i++ {
        vals[i] = rng.Intn(20) + 1
    }
    edges := generateTree(rng, n)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", n))
    for i := 0; i < n; i++ {
        if i > 0 { sb.WriteByte(' ') }
        sb.WriteString(fmt.Sprintf("%d", vals[i]))
    }
    sb.WriteByte('\n')
    for _, e := range edges {
        sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
    }
    return sb.String(), solveF(n, vals, edges)
}

func runCase(bin, in, exp string) error {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    for i := 0; i < 100; i++ {
        in, exp := generateF(rng)
        if err := runCase(bin, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

