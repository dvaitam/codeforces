package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type TestCase struct {
    input string
    ans   string
}

func poly(a []int) []int {
    coeffs := []int{1}
    for _, ai := range a {
        m := len(coeffs)
        nxt := make([]int, m+1)
        for j := 0; j < m; j++ {
            nxt[j] += ai * coeffs[j]
            nxt[j+1] += coeffs[j]
        }
        coeffs = nxt
    }
    return coeffs
}

func formatPoly(coeffs []int) string {
    var sb strings.Builder
    first := true
    deg := len(coeffs) - 1
    for k := deg; k >= 0; k-- {
        c := coeffs[k]
        if c == 0 {
            continue
        }
        if first {
            if c < 0 {
                sb.WriteByte('-')
                c = -c
            }
            first = false
        } else {
            if c < 0 {
                sb.WriteString(" - ")
                c = -c
            } else {
                sb.WriteString(" + ")
            }
        }
        if k == 0 {
            sb.WriteString(strconv.Itoa(c))
        } else {
            if c != 1 {
                sb.WriteString(strconv.Itoa(c))
                sb.WriteByte('*')
            }
            sb.WriteByte('X')
            if k != 1 {
                sb.WriteByte('^')
                sb.WriteString(strconv.Itoa(k))
            }
        }
    }
    return "p(x) = " + sb.String()
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(6))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(4) + 1
        a := make([]int, n)
        for j := 0; j < n; j++ {
            v := r.Intn(7) - 3
            a[j] = v
        }
        coeffs := poly(a)
        ans := formatPoly(coeffs)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for j, v := range a {
            if j > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprintf("%d", v))
        }
        sb.WriteString("\n")
        cases[i] = TestCase{input: sb.String(), ans: ans}
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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

