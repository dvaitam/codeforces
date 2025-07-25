package main

import (
    "bufio"
    "bytes"
    "context"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

type Edge struct { to int; c int }

type State struct { v int; used int }

func solveE(in io.Reader) string {
    reader := bufio.NewReader(in)
    var n, m int
    if _, err := fmt.Fscan(reader, &n, &m); err != nil {
        return ""
    }
    g := make([][]Edge, n+1)
    for i := 0; i < m; i++ {
        var x, y, z int
        fmt.Fscan(reader, &x, &y, &z)
        g[x] = append(g[x], Edge{y, z})
        g[y] = append(g[y], Edge{x, z})
    }
    var a, b int
    fmt.Fscan(reader, &a, &b)

    visited := make([][2]bool, n+1)
    q := []State{{a, 0}}
    visited[a][0] = true
    for head := 0; head < len(q); head++ {
        cur := q[head]
        if cur.v == b && cur.used == 1 {
            return "YES"
        }
        for _, e := range g[cur.v] {
            next := cur.used | e.c
            if !visited[e.to][next] {
                visited[e.to][next] = true
                q = append(q, State{e.to, next})
            }
        }
    }
    return "NO"
}

func genTests() []string {
    rng := rand.New(rand.NewSource(5))
    tests := make([]string, 100)
    for i := 0; i < 100; i++ {
        n := rng.Intn(6) + 2
        maxEdges := n * (n - 1) / 2
        m := rng.Intn(maxEdges) + 1
        edges := make([][3]int, m)
        used := map[[2]int]bool{}
        for j := 0; j < m; j++ {
            x := rng.Intn(n) + 1
            y := rng.Intn(n) + 1
            for x == y || used[[2]int{min(x,y), max(x,y)}] {
                x = rng.Intn(n) + 1
                y = rng.Intn(n) + 1
            }
            used[[2]int{min(x,y), max(x,y)}] = true
            z := rng.Intn(2)
            edges[j] = [3]int{x, y, z}
        }
        a := rng.Intn(n) + 1
        b := rng.Intn(n) + 1
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
        for _, e := range edges {
            sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
        }
        sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
        tests[i] = sb.String()
    }
    return tests
}

func min(a,b int) int { if a<b { return a }; return b }
func max(a,b int) int { if a>b { return a }; return b }

func runBinary(bin, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run verifierE.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        expected := solveE(strings.NewReader(tc))
        actual, err := runBinary(bin, tc)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
            return
        }
        if actual != strings.TrimSpace(expected) {
            fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, tc, expected, actual)
            return
        }
    }
    fmt.Println("All tests passed!")
}

