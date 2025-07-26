package main

import (
    "bytes"
    "fmt"
    "math/big"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type TestCase struct {
    input string
    ans   string
}

func solve(a, b string) string {
    A := new(big.Int)
    B := new(big.Int)
    A.SetString(a, 10)
    B.SetString(b, 10)
    return new(big.Int).Add(A, B).String()
}

func randomInt(r *rand.Rand, digits int) string {
    b := make([]byte, digits)
    for i := 0; i < digits; i++ {
        d := byte(r.Intn(10))
        if i == 0 && d == 0 {
            d = byte(r.Intn(9) + 1)
        }
        b[i] = '0' + d
    }
    return string(b)
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(3))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        d1 := r.Intn(20) + 1
        d2 := r.Intn(20) + 1
        a := randomInt(r, d1)
        b := randomInt(r, d2)
        input := fmt.Sprintf("%s\n%s\n", a, b)
        ans := solve(a, b)
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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

