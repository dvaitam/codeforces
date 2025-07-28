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

func solveB(arr []int) string {
    cnt1, cnt2 := 0, 0
    for _, v := range arr {
        if v == 1 {
            cnt1++
        } else if v == 2 {
            cnt2++
        }
    }
    if cnt1%2 != 0 {
        return "NO"
    }
    if cnt1 == 0 && cnt2%2 == 1 {
        return "NO"
    }
    return "YES"
}

func generateTests() []TestCase {
    r := rand.New(rand.NewSource(42))
    tests := make([]TestCase, 100)
    for i := range tests {
        n := r.Intn(10) + 1
        arr := make([]int, n)
        var sb strings.Builder
        for j := 0; j < n; j++ {
            if j > 0 {
                sb.WriteByte(' ')
            }
            val := 1 + r.Intn(2)
            arr[j] = val
            fmt.Fprintf(&sb, "%d", val)
        }
        expected := solveB(arr)
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
        fmt.Println("usage: go run verifierB.go /path/to/binary")
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

