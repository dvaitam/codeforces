package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
)

type TestCase struct {
    input    string
    expected string
}

func solveG(n, m int, edges [][2]int) []int {
    g := make([][]int, n)
    for _, e := range edges {
        u, v := e[0], e[1]
        u--
        v--
        g[u] = append(g[u], v)
    }
    dist := make([]int, n)
    for i := range dist {
        dist[i] = -1
    }
    q := make([]int, 0, n)
    q = append(q, 0)
    dist[0] = 0
    for head := 0; head < len(q); head++ {
        u := q[head]
        for _, v := range g[u] {
            if dist[v] == -1 {
                dist[v] = dist[u] + 1
                q = append(q, v)
            }
        }
    }
    ord := make([]int, n)
    for i := range ord {
        ord[i] = i
    }
    sort.Slice(ord, func(i, j int) bool { return dist[ord[i]] > dist[ord[j]] })
    dp := make([]int, n)
    copy(dp, dist)
    for _, u := range ord {
        for _, v := range g[u] {
            if dist[u] < dist[v] {
                if dp[u] > dp[v] {
                    dp[u] = dp[v]
                }
            } else {
                if dp[u] > dist[v] {
                    dp[u] = dist[v]
                }
            }
        }
    }
    return dp
}

func generateTests() []TestCase {
    r := rand.New(rand.NewSource(42))
    tests := make([]TestCase, 100)
    for i := range tests {
        n := r.Intn(6) + 1
        m := r.Intn(n*(n-1)/2 + 1)
        edges := make([][2]int, m)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("1\n%d %d\n", n, m))
        for j := 0; j < m; j++ {
            u := r.Intn(n) + 1
            v := r.Intn(n) + 1
            edges[j] = [2]int{u, v}
            sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
        }
        ans := solveG(n, m, edges)
        var out strings.Builder
        for j, v := range ans {
            if j > 0 {
                out.WriteByte(' ')
            }
            out.WriteString(fmt.Sprintf("%d", v))
        }
        tests[i] = TestCase{input: sb.String(), expected: strings.TrimSpace(out.String())}
    }
    return tests
}

func run(binary string, tc TestCase) (string, error) {
    cmd := exec.Command(binary)
    cmd.Stdin = strings.NewReader(tc.input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierG.go /path/to/binary")
        os.Exit(1)
    }
    binary := os.Args[1]
    tests := generateTests()
    for i, tc := range tests {
        got, err := run(binary, tc)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if got != tc.expected {
            fmt.Printf("Test %d failed: expected %q got %q\nInput:\n%s", i+1, tc.expected, got, tc.input)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

