package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type test struct {
    input    string
    expected string
}

func solve(input string) string {
    lines := strings.Split(strings.TrimSpace(input), "\n")
    if len(lines) < 2 {
        return ""
    }
    a := strings.TrimSpace(lines[0])
    b := strings.TrimSpace(lines[1])
    n := len(a)
    res := make([]byte, n)
    for i := 0; i < n; i++ {
        if a[i] != b[i] {
            res[i] = '1'
        } else {
            res[i] = '0'
        }
    }
    return string(res)
}

func generateTests() []test {
    rand.Seed(42)
    var tests []test
    fixed := []struct{ a, b string }{
        {"0", "0"},
        {"1", "0"},
        {"1010", "0101"},
        {"1111", "0000"},
        {"101010", "101010"},
    }
    for _, f := range fixed {
        inp := fmt.Sprintf("%s\n%s\n", f.a, f.b)
        tests = append(tests, test{inp, solve(inp)})
    }
    for len(tests) < 100 {
        n := rand.Intn(100) + 1
        a := make([]byte, n)
        b := make([]byte, n)
        for i := 0; i < n; i++ {
            a[i] = byte('0' + rand.Intn(2))
            b[i] = byte('0' + rand.Intn(2))
        }
        inp := fmt.Sprintf("%s\n%s\n", string(a), string(b))
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for i, t := range tests {
        got, err := runBinary(bin, t.input)
        if err != nil {
            fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if got != strings.TrimSpace(t.expected) {
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

