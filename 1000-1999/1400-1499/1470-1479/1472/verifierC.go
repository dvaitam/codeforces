package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type TestCase struct {
    input    string
    expected string
}

func solveC(a []int64) string {
    n := len(a)
    dp := make([]int64, n)
    var ans int64
    for i := n - 1; i >= 0; i-- {
        dp[i] = a[i]
        next := i + int(a[i])
        if next < n {
            dp[i] += dp[next]
        }
        if dp[i] > ans {
            ans = dp[i]
        }
    }
    return fmt.Sprintf("%d", ans)
}

func generateTests() []TestCase {
    r := rand.New(rand.NewSource(42))
    tests := make([]TestCase, 100)
    for i := range tests {
        n := r.Intn(10) + 1
        arr := make([]int64, n)
        var sb strings.Builder
        for j := 0; j < n; j++ {
            arr[j] = int64(r.Intn(10) + 1)
            if j > 0 {
                sb.WriteByte(' ')
            }
            fmt.Fprintf(&sb, "%d", arr[j])
        }
        expected := solveC(arr)
        input := fmt.Sprintf("1\n%d\n%s\n", n, sb.String())
        tests[i] = TestCase{input: input, expected: expected}
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
        fmt.Println("usage: go run verifierC.go /path/to/binary")
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

