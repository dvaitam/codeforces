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
    r := strings.NewReader(strings.TrimSpace(input))
    var t int
    fmt.Fscan(r, &t)
    var out strings.Builder
    for ; t > 0; t-- {
        var s string
        fmt.Fscan(r, &s)
        dp0, dp1 := 0, 0
        var total int64
        for i := 0; i < len(s); i++ {
            ch := s[i]
            expected0, expected1 := byte('0'), byte('1')
            if i%2 == 1 {
                expected0, expected1 = '1', '0'
            }
            if ch == '?' || ch == expected0 {
                dp0++
            } else {
                dp0 = 0
            }
            if ch == '?' || ch == expected1 {
                dp1++
            } else {
                dp1 = 0
            }
            if dp0 > dp1 {
                total += int64(dp0)
            } else {
                total += int64(dp1)
            }
        }
        out.WriteString(fmt.Sprintf("%d\n", total))
    }
    return out.String()
}

func generateTests() []test {
    rand.Seed(3)
    var tests []test
    fixed := []string{"0??10", "0", "???", "1010", "?1??1"}
    for _, s := range fixed {
        inp := fmt.Sprintf("1\n%s\n", s)
        tests = append(tests, test{inp, solve(inp)})
    }
    for len(tests) < 100 {
        l := rand.Intn(15) + 1
        var sb strings.Builder
        sb.WriteString("1\n")
        for i := 0; i < l; i++ {
            r := rand.Intn(3)
            if r == 0 {
                sb.WriteByte('0')
            } else if r == 1 {
                sb.WriteByte('1')
            } else {
                sb.WriteByte('?')
            }
        }
        sb.WriteString("\n")
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", i+1, t.input, t.expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

