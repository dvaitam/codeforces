package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type edge struct{ to int; w int64 }

type test struct{ input, expected string }

func solve(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var n int
    if _, err := fmt.Fscan(in, &n); err != nil { return "" }
    adj := make([][]edge, n+1)
    var total int64
    for i:=0;i<n-1;i++ {
        var u,v int; var w int64
        fmt.Fscan(in,&u,&v,&w)
        adj[u] = append(adj[u], edge{v,w})
        adj[v] = append(adj[v], edge{u,w})
        total += w
    }
    dist := make([]int64,n+1)
    vis := make([]bool,n+1)
    q := []int{1}; vis[1]=true
    var maxd int64
    for idx:=0; idx<len(q); idx++ {
        cur := q[idx]
        for _, e := range adj[cur] {
            if !vis[e.to] {
                vis[e.to]=true
                dist[e.to]=dist[cur]+e.w
                if dist[e.to]>maxd { maxd = dist[e.to] }
                q = append(q,e.to)
            }
        }
    }
    ans := total*2 - maxd
    return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
    rand.Seed(45)
    var tests []test
    fixed := []string{ "1\n", "2\n1 2 5\n" }
    for _, f := range fixed {
        tests = append(tests, test{f, solve(f)})
    }
    for len(tests) < 100 {
        n := rand.Intn(15)+1
        if n==1 { tests = append(tests, test{"1\n", "0"}); continue }
        edges := make([][3]int64, n-1)
        for i:=0;i<n-1;i++ { edges[i][0]=int64(i+1); edges[i][1]=int64(rand.Intn(i+1)+1); edges[i][2]=int64(rand.Intn(1000)) }
        var sb strings.Builder
        sb.WriteString(strconv.Itoa(n)+"\n")
        for _, e := range edges {
            sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
        }
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for i, t := range tests {
        got, err := runBinary(bin, t.input)
        if err != nil {
            fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if got != strings.TrimSpace(t.expected) {
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

