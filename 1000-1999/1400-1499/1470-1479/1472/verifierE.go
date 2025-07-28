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

func solveE(h, w []int) []int {
    n := len(h)
    for i := 0; i < n; i++ {
        if h[i] > w[i] {
            h[i], w[i] = w[i], h[i]
        }
    }
    idx := make([]int, n)
    for i := range idx { idx[i] = i }
    sort.Slice(idx, func(i, j int) bool { return h[idx[i]] < h[idx[j]] })
    tmp := -1
    ans := make([]int, n)
    for i := range ans { ans[i] = -1 }
    for i := 0; i < n; {
        j := i
        hh := h[idx[i]]
        for j < n && h[idx[j]] == hh { j++ }
        for k := i; k < j; k++ {
            id := idx[k]
            if tmp != -1 && w[tmp] < w[id] {
                ans[id] = tmp
            }
        }
        for k := i; k < j; k++ {
            id := idx[k]
            if tmp == -1 || w[tmp] > w[id] {
                tmp = id
            }
        }
        i = j
    }
    res := make([]int, n)
    for i := range ans {
        if ans[i] >= 0 {
            res[i] = ans[i] + 1
        } else {
            res[i] = -1
        }
    }
    return res
}

func generateTests() []TestCase {
    r := rand.New(rand.NewSource(42))
    tests := make([]TestCase, 100)
    for i := range tests {
        n := r.Intn(10) + 1
        h := make([]int, n)
        w := make([]int, n)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("1\n%d\n", n))
        for j := 0; j < n; j++ {
            h[j] = r.Intn(20) + 1
            w[j] = r.Intn(20) + 1
            sb.WriteString(fmt.Sprintf("%d %d\n", h[j], w[j]))
        }
        ans := solveE(append([]int(nil), h...), append([]int(nil), w...))
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
        fmt.Println("usage: go run verifierE.go /path/to/binary")
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

