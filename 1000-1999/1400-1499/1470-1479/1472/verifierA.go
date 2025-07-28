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

func solveA(w, h, n int64) string {
    count := int64(1)
    for w%2 == 0 {
        count <<= 1
        w >>= 1
    }
    for h%2 == 0 {
        count <<= 1
        h >>= 1
    }
    if count >= n {
        return "YES"
    }
    return "NO"
}

func generateTests() []TestCase {
    r := rand.New(rand.NewSource(42))
    tests := make([]TestCase, 100)
    for i := range tests {
        w := r.Int63n(1000) + 1
        h := r.Int63n(1000) + 1
        n := r.Int63n(1000) + 1
        expected := solveA(w, h, n)
        input := fmt.Sprintf("1\n%d %d %d\n", w, h, n)
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
    cmd.Env = append(os.Environ(), "LC_ALL=C")
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    binary := os.Args[1]
    tests := generateTests()
    for i, tc := range tests {
        got, err := run(binary, tc)
        if err != nil {
            os.Exit(1)
        }
        if got != tc.expected {
            fmt.Printf("Test %d failed: expected %q got %q\nInput:\n%s", i+1, tc.expected, got, tc.input)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

