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

type testH struct {
    a []int
    b []int
}

func genTestsH() []testH {
    rng := rand.New(rand.NewSource(49))
    tests := make([]testH, 100)
    for i := range tests {
        n := rng.Intn(5) + 1
        m := rng.Intn(5) + 1
        A := make([]int, n)
        B := make([]int, m)
        for j := 0; j < n; j++ { A[j] = rng.Intn(2) }
        for j := 0; j < m; j++ { B[j] = rng.Intn(2) }
        tests[i] = testH{a: A, b: B}
    }
    return tests
}

func solveOps(a, b []int) int {
    n := len(a)
    m := len(b)
    total := n + m
    ia, ib, im := n-1, m-1, total-1
    c := 0
    ops := 0
    for ia >= 0 || ib >= 0 {
        for ia >= 0 && a[ia] == c { ia--; im-- }
        for ib >= 0 && b[ib] == c { ib--; im-- }
        ops++
        if c == 0 { c = 1 } else { c = 0 }
    }
    return ops - 1
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

func applyOps(order []int, a, b []int, flips []int) ([]int, error) {
    n := len(a)
    m := len(b)
    vals := make([]int, n+m)
    for i, idx := range order {
        if idx < 1 || idx > n+m { return nil, fmt.Errorf("index out of range") }
        if idx <= n {
            vals[i] = a[idx-1]
        } else {
            vals[i] = b[idx-n-1]
        }
    }
    // check relative order
    posA := 0
    posB := 0
    for _, idx := range order {
        if idx <= n {
            posA++
            if idx != posA {
                return nil, fmt.Errorf("order not preserved")
            }
        } else {
            posB++
            if idx != n+posB {
                return nil, fmt.Errorf("order not preserved")
            }
        }
    }
    for _, k := range flips {
        if k < 1 || k > len(vals) { return nil, fmt.Errorf("bad flip size") }
        for i := 0; i < k/2; i++ {
            vals[i], vals[k-1-i] = vals[k-1-i]^1, vals[i]^1
        }
        if k%2 == 1 { vals[k/2] ^= 1 }
    }
    return vals, nil
}

func runCase(bin string, tc testH) error {
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d\n", len(tc.a))
    for i, v := range tc.a { if i>0 {sb.WriteByte(' ')}; sb.WriteString(strconv.Itoa(v)) }
    sb.WriteByte('\n')
    fmt.Fprintf(&sb, "%d\n", len(tc.b))
    for i, v := range tc.b { if i>0 {sb.WriteByte(' ')}; sb.WriteString(strconv.Itoa(v)) }
    sb.WriteByte('\n')
    out, err := run(bin, sb.String())
    if err != nil { return err }
    scanner := bufio.NewScanner(strings.NewReader(out))
    scanner.Split(bufio.ScanWords)
    order := make([]int, len(tc.a)+len(tc.b))
    for i := range order {
        if !scanner.Scan() { return fmt.Errorf("missing order element") }
        v, err := strconv.Atoi(scanner.Text())
        if err != nil { return fmt.Errorf("bad order element") }
        order[i] = v
    }
    if !scanner.Scan() { return fmt.Errorf("missing op count") }
    x, err := strconv.Atoi(scanner.Text())
    if err != nil { return fmt.Errorf("bad op count") }
    flips := make([]int, x)
    for i := 0; i < x; i++ {
        if !scanner.Scan() { return fmt.Errorf("missing flip value") }
        v, err := strconv.Atoi(scanner.Text())
        if err != nil { return fmt.Errorf("bad flip value") }
        flips[i] = v
    }
    if scanner.Scan() { return fmt.Errorf("extra output") }
    result, err := applyOps(order, tc.a, tc.b, flips)
    if err != nil { return err }
    for _, v := range result {
        if v != 0 { return fmt.Errorf("final deck not all face down") }
    }
    if x != solveOps(tc.a, tc.b) { return fmt.Errorf("operations not minimal") }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsH()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

