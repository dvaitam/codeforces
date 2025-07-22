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

type testA struct {
    n int
    s string
}

func genTestsA() []testA {
    rng := rand.New(rand.NewSource(42))
    tests := make([]testA, 100)
    for i := range tests {
        n := rng.Intn(49)*2 + 4 // even between 4 and 100
        b := make([]byte, n)
        for j := range b {
            if rng.Intn(2) == 0 {
                b[j] = 'L'
            } else {
                b[j] = 'R'
            }
        }
        tests[i] = testA{n: n, s: string(b)}
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

func checkOutputA(tc testA, out string) error {
    scanner := bufio.NewScanner(strings.NewReader(out))
    scanner.Split(bufio.ScanWords)
    used := make([]bool, tc.n+1)
    pairs := 0
    for pairs < tc.n/2 {
        if !scanner.Scan() {
            return fmt.Errorf("insufficient output")
        }
        aStr := scanner.Text()
        if !scanner.Scan() {
            return fmt.Errorf("insufficient output")
        }
        bStr := scanner.Text()
        a, err := strconv.Atoi(aStr)
        if err != nil {
            return fmt.Errorf("invalid integer %s", aStr)
        }
        b, err := strconv.Atoi(bStr)
        if err != nil {
            return fmt.Errorf("invalid integer %s", bStr)
        }
        if a < 1 || a > tc.n || b < 1 || b > tc.n {
            return fmt.Errorf("index out of range")
        }
        if used[a] || used[b] || a == b {
            return fmt.Errorf("student repeated")
        }
        if abs(a-b) == 1 {
            return fmt.Errorf("adjacent students together")
        }
        if tc.s[a-1] == 'R' && tc.s[b-1] == 'L' {
            return fmt.Errorf("hand conflict")
        }
        used[a] = true
        used[b] = true
        pairs++
    }
    if scanner.Scan() {
        return fmt.Errorf("extra output")
    }
    for i := 1; i <= tc.n; i++ {
        if !used[i] {
            return fmt.Errorf("student %d missing", i)
        }
    }
    return nil
}

func abs(x int) int { if x < 0 { return -x }; return x }

func runCase(bin string, tc testA) error {
    input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
    out, err := run(bin, input)
    if err != nil { return err }
    return checkOutputA(tc, out)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsA()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%v %s\n", i+1, err, tc.n, tc.s)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

