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

func solve(n int, s, t string) string {
    ls := len(s)
    lt := len(t)
    i := 0
    max := ls
    if lt < max {
        max = lt
    }
    for i < max && s[i] == t[i] {
        i++
    }
    ops := (ls - i) + (lt - i)
    if ops <= n {
        return "YES"
    }
    return "NO"
}

func genString(r *rand.Rand, l int) string {
    if l <= 0 {
        l = 1
    }
    b := make([]byte, l)
    for i := range b {
        b[i] = byte('a' + r.Intn(3))
    }
    return string(b)
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(4))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(20) + 1
        s := genString(r, r.Intn(10)+1)
        t := genString(r, r.Intn(10)+1)
        input := fmt.Sprintf("%d\n%s\n%s\n", n, s, t)
        ans := solve(n, s, t)
        cases[i] = TestCase{input: input, ans: ans}
    }
    return cases
}

func run(bin, in string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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

