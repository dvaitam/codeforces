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

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func solve(input string) string {
    r := strings.NewReader(strings.TrimSpace(input))
    var t int
    fmt.Fscan(r, &t)
    var out strings.Builder
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(r, &n)
        arr := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(r, &arr[i])
        }
        evens := make([]int, 0)
        odds := make([]int, 0)
        for _, v := range arr {
            if v%2 == 0 {
                evens = append(evens, v)
            } else {
                odds = append(odds, v)
            }
        }
        b := append(evens, odds...)
        count := 0
        for i := 0; i < n; i++ {
            for j := i + 1; j < n; j++ {
                if gcd(b[i], 2*b[j]) > 1 {
                    count++
                }
            }
        }
        out.WriteString(fmt.Sprintf("%d\n", count))
    }
    return out.String()
}

func generateTests() []test {
    rand.Seed(2)
    var tests []test
    fixed := [][]int{
        {6, 3, 5, 3},
        {1, 1},
        {2, 4},
        {5, 7, 9, 11},
    }
    for _, arr := range fixed {
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
        for i, v := range arr {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprintf("%d", v))
        }
        sb.WriteString("\n")
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    for len(tests) < 100 {
        n := rand.Intn(8) + 2
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i := 0; i < n; i++ {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprintf("%d", rand.Intn(100)+1))
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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

