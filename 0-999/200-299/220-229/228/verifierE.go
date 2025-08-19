package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type edge struct{ u, v, c int }

func runBinary(path, input string) (string, error) {
    cmd := exec.Command(path)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func genTests() []string {
    r := rand.New(rand.NewSource(42))
    tests := []string{}
    for len(tests) < 100 {
        n := r.Intn(3) + 2 // 2..4
        maxM := n * (n - 1) / 2
        m := r.Intn(maxM + 1)
        var sb strings.Builder
        fmt.Fprintf(&sb, "%d %d\n", n, m)
        edges := make(map[[2]int]bool)
        for i := 0; i < m; i++ {
            var a, b int
            for {
                a = r.Intn(n) + 1
                b = r.Intn(n) + 1
                if a != b && !edges[[2]int{a, b}] && !edges[[2]int{b, a}] {
                    break
                }
            }
            edges[[2]int{a, b}] = true
            c := r.Intn(2)
            fmt.Fprintf(&sb, "%d %d %d\n", a, b, c)
        }
        tests = append(tests, sb.String())
    }
    return tests
}

// Check satisfiability of constraints and validate a candidate selection
// Edges: want (sel[u] xor sel[v]) == 1 - c
func validateSelection(input, output string) error {
    lines := strings.Split(strings.TrimSpace(input), "\n")
    var n, m int
    fmt.Sscanf(lines[0], "%d %d", &n, &m)
    es := make([]edge, m)
    for i := 0; i < m; i++ {
        var a, b, c int
        fmt.Sscanf(lines[1+i], "%d %d %d", &a, &b, &c)
        es[i] = edge{u: a, v: b, c: c}
    }
    outLines := []string{}
    for _, ln := range strings.Split(strings.TrimSpace(output), "\n") {
        t := strings.TrimSpace(ln)
        if t != "" {
            outLines = append(outLines, t)
        }
    }
    if len(outLines) == 0 {
        return fmt.Errorf("empty output")
    }
    if strings.EqualFold(outLines[0], "impossible") {
        if satisfiable(n, es) {
            return fmt.Errorf("claimed impossible but constraints satisfiable")
        }
        return nil
    }
    // parse count and list
    k, err := strconv.Atoi(outLines[0])
    if err != nil || k < 0 || k > n {
        return fmt.Errorf("invalid count")
    }
    sel := make([]int, n+1)
    if k == 0 {
        // no cities
    } else {
        if len(outLines) < 2 {
            return fmt.Errorf("missing city list")
        }
        fields := strings.Fields(outLines[1])
        // allow duplicates (multiple days on same city)
        for _, f := range fields {
            v, e := strconv.Atoi(f)
            if e != nil || v < 1 || v > n {
                return fmt.Errorf("invalid city index")
            }
            sel[v] ^= 1
        }
        // ensure parity of toggles equals reported k modulo 2? Not necessary.
    }
    // Check constraints
    for _, e := range es {
        if (sel[e.u]^sel[e.v]) != (1 - e.c) {
            return fmt.Errorf("selection does not satisfy edge %d-%d-%d", e.u, e.v, e.c)
        }
    }
    return nil
}

// DSU with parity to test satisfiability
func satisfiable(n int, es []edge) bool {
    par := make([]int, n+1)
    d := make([]int, n+1)
    for i := 1; i <= n; i++ { par[i] = i }
    var fp func(int) int
    fp = func(x int) int {
        if par[x] == x { return x }
        p := par[x]
        r := fp(p)
        d[x] ^= d[p]
        par[x] = r
        return r
    }
    for _, e := range es {
        u, v, c := e.u, e.v, e.c
        // want sel[u]^sel[v] == 1-c
        want := 1 - c
        ru := fp(u); rv := fp(v)
        if ru == rv {
            if (d[u]^d[v]) != want { return false }
            continue
        }
        par[ru] = rv
        d[ru] = want ^ d[u] ^ d[v]
    }
    return true
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, in := range tests {
        out, err := runBinary(bin, in)
        if err != nil {
            fmt.Printf("test %d runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if err := validateSelection(in, out); err != nil {
            fmt.Printf("test %d failed:\ninput:\n%serror: %v\noutput:\n%s\n", i+1, in, err, out)
            os.Exit(1)
        }
    }
    fmt.Printf("ok %d tests\n", len(tests))
}
