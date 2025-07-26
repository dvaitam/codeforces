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
    input string
    ans   string
}

func solve(n, k, n1 int) string {
    if n1 >= n || k >= 4 {
        return "YES"
    }
    return "NO"
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(1))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(20) + 1
        k := r.Intn(6)
        n1 := r.Intn(20) + 1
        input := fmt.Sprintf("%d %d %d\n", n, k, n1)
        ans := solve(n, k, n1)
        cases[i] = TestCase{input: input, ans: ans}
    }
    return cases
}

func run(bin string, in string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        got, err := run(bin, tc.input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != tc.ans {
            fmt.Fprintf(os.Stderr, "test %d failed: input %q expected %q got %q\n", i+1, tc.input, tc.ans, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

