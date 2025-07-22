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

type testF struct {
    n int
    a int
    b int
    h []int
}

func genTestsF() []testF {
    rng := rand.New(rand.NewSource(47))
    tests := make([]testF, 100)
    for i := range tests {
        n := rng.Intn(8) + 1
        h := make([]int, n)
        total := 0
        for j := range h {
            h[j] = rng.Intn(10) + 1
            total += h[j]
        }
        a := rng.Intn(total + 10)
        b := rng.Intn(total + 10)
        tests[i] = testF{n: n, a: a, b: b, h: h}
    }
    return tests
}

func min(a, b int) int { if a < b { return a }; return b }

func solveF(tc testF) int {
    n := tc.n
    a := tc.a
    b := tc.b
    h := tc.h
    total := 0
    for _, v := range h { total += v }
    if total > a+b { return -1 }
    const INF = int(1e9)
    prev0 := make([]int, a+1)
    prev1 := make([]int, a+1)
    for i := range prev0 { prev0[i] = INF; prev1[i] = INF }
    if h[0] <= a { prev0[h[0]] = 0 }
    prev1[0] = 0
    for i := 1; i < n; i++ {
        hi := h[i]
        hip := h[i-1]
        next0 := make([]int, a+1)
        next1 := make([]int, a+1)
        for j := range next0 { next0[j] = INF; next1[j] = INF }
        cost := min(hip, hi)
        for r := 0; r <= a; r++ {
            v0 := prev0[r]
            if v0 < INF {
                nr := r + hi
                if nr <= a && v0 < next0[nr] {
                    next0[nr] = v0
                }
                if v0+cost < next1[r] { next1[r] = v0 + cost }
            }
            v1 := prev1[r]
            if v1 < INF {
                nr := r + hi
                if nr <= a && v1+cost < next0[nr] { next0[nr] = v1 + cost }
                if v1 < next1[r] { next1[r] = v1 }
            }
        }
        prev0, prev1 = next0, next1
    }
    ans := INF
    start := 0
    if total > b { start = total - b }
    for r := start; r <= a; r++ {
        if prev0[r] < ans { ans = prev0[r] }
        if prev1[r] < ans { ans = prev1[r] }
    }
    if ans == INF { return -1 }
    return ans
}

func run(bin string, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return out.String(), nil
}

func runCase(bin string, tc testF) error {
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d\n%d %d\n", tc.n, tc.a, tc.b)
    for i, v := range tc.h {
        if i > 0 { sb.WriteByte(' ') }
        sb.WriteString(strconv.Itoa(v))
    }
    sb.WriteByte('\n')
    out, err := run(bin, sb.String())
    if err != nil { return err }
    expected := solveF(tc)
    valStr := strings.TrimSpace(out)
    val, err := strconv.Atoi(valStr)
    if err != nil { return fmt.Errorf("bad output") }
    if val != expected { return fmt.Errorf("expected %d got %d", expected, val) }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsF()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

