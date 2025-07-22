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

type testG struct {
    n int
}

func genTestsG() []testG {
    rng := rand.New(rand.NewSource(48))
    tests := make([]testG, 100)
    for i := range tests {
        tests[i] = testG{n: rng.Intn(18) + 2}
    }
    return tests
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

func checkSchedule(n int, m int, practices [][]int) error {
    r := 0
    for (1<<r) < n { r++ }
    if m != r { return fmt.Errorf("m should be %d", r) }
    for _, p := range practices {
        if len(p) == 0 || len(p) >= n { return fmt.Errorf("invalid team size") }
        seen := make(map[int]bool)
        for _, v := range p {
            if v < 1 || v > n { return fmt.Errorf("index out of range") }
            if seen[v] { return fmt.Errorf("duplicate index in practice") }
            seen[v] = true
        }
    }
    sep := make([][]bool, n)
    for i := range sep { sep[i] = make([]bool, n) }
    for _, p := range practices {
        mark := make([]bool, n+1)
        for _, v := range p { mark[v] = true }
        for i := 1; i <= n; i++ {
            for j := i + 1; j <= n; j++ {
                if mark[i] != mark[j] {
                    sep[i-1][j-1] = true
                    sep[j-1][i-1] = true
                }
            }
        }
    }
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            if !sep[i][j] { return fmt.Errorf("pair %d %d never separated", i+1, j+1) }
        }
    }
    return nil
}

func runCase(bin string, tc testG) error {
    input := fmt.Sprintf("%d\n", tc.n)
    out, err := run(bin, input)
    if err != nil { return err }
    scanner := bufio.NewScanner(strings.NewReader(out))
    scanner.Split(bufio.ScanWords)
    if !scanner.Scan() { return fmt.Errorf("missing m") }
    m, err := strconv.Atoi(scanner.Text())
    if err != nil { return fmt.Errorf("invalid m") }
    practices := make([][]int, m)
    for i := 0; i < m; i++ {
        if !scanner.Scan() { return fmt.Errorf("missing team size") }
        f, err := strconv.Atoi(scanner.Text())
        if err != nil { return fmt.Errorf("bad team size") }
        p := make([]int, f)
        for j := 0; j < f; j++ {
            if !scanner.Scan() { return fmt.Errorf("missing player index") }
            v, err := strconv.Atoi(scanner.Text())
            if err != nil { return fmt.Errorf("bad index") }
            p[j] = v
        }
        practices[i] = p
    }
    if scanner.Scan() { return fmt.Errorf("extra output") }
    return checkSchedule(tc.n, m, practices)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsG()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

