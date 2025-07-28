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

func solveF(n, m int, pairs [][2]int) string {
    colMask := make(map[int]int)
    for _, p := range pairs {
        r, c := p[0], p[1]
        if r == 1 {
            colMask[c] |= 1
        } else {
            colMask[c] |= 2
        }
    }
    cols := make([]int, 0, len(colMask))
    for c := range colMask {
        cols = append(cols, c)
    }
    sort.Ints(cols)
    pendingRow := 0
    pendingCol := 0
    ok := true
    for _, c := range cols {
        mask := colMask[c]
        if mask == 3 {
            if pendingRow != 0 {
                ok = false
                break
            }
            continue
        }
        row := 1
        if mask == 2 {
            row = 2
        }
        if pendingRow == 0 {
            pendingRow = row
            pendingCol = c
        } else {
            diff := c - pendingCol
            same := 0
            if row == pendingRow {
                same = 1
            }
            if diff%2 != same {
                ok = false
                break
            }
            pendingRow = 0
        }
    }
    if pendingRow != 0 {
        ok = false
    }
    if ok {
        return "YES"
    }
    return "NO"
}

func generateTests() []TestCase {
    r := rand.New(rand.NewSource(42))
    tests := make([]TestCase, 100)
    for i := range tests {
        n := r.Intn(5) + 1
        m := r.Intn(10)
        pairs := make([][2]int, m)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("1\n%d %d\n", n, m))
        for j := 0; j < m; j++ {
            r1 := r.Intn(2) + 1
            c := r.Intn(10) + 1
            pairs[j] = [2]int{r1, c}
            sb.WriteString(fmt.Sprintf("%d %d\n", r1, c))
        }
        expected := solveF(n, m, pairs)
        tests[i] = TestCase{input: sb.String(), expected: expected}
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
        fmt.Println("usage: go run verifierF.go /path/to/binary")
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

