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
    var k int
    if _, err := fmt.Fscan(r, &k); err != nil {
        return ""
    }
    var s string
    fmt.Fscan(r, &s)
    n := (1 << k) - 1
    posToNode := make([]int, n+1)
    pos := 1
    for level := k - 1; level >= 0; level-- {
        for node := 1 << level; node <= (1<<(level+1))-1; node++ {
            posToNode[pos] = node
            pos++
        }
    }
    nodeChar := make([]byte, n+1)
    for p := 1; p <= n; p++ {
        node := posToNode[p]
        nodeChar[node] = s[p-1]
    }
    dp := make([]int, n+1)
    var recalc func(int)
    recalc = func(node int) {
        for {
            ch := nodeChar[node]
            left := node * 2
            if left > n {
                if ch == '?' {
                    dp[node] = 2
                } else {
                    dp[node] = 1
                }
            } else {
                right := left + 1
                if ch == '?' {
                    dp[node] = dp[left] + dp[right]
                } else if ch == '0' {
                    dp[node] = dp[left]
                } else {
                    dp[node] = dp[right]
                }
            }
            if node == 1 {
                break
            }
            node /= 2
        }
    }
    for node := n; node >= 1; node-- {
        recalc(node)
    }
    var q int
    fmt.Fscan(r, &q)
    var out strings.Builder
    for ; q > 0; q-- {
        var p int
        var c string
        fmt.Fscan(r, &p, &c)
        node := posToNode[p]
        nodeChar[node] = c[0]
        recalc(node)
        out.WriteString(fmt.Sprintf("%d\n", dp[1]))
    }
    return out.String()
}

func generateTests() []test {
    rand.Seed(4)
    var tests []test
    // fixed simple test
    inp := "1\n?\n1\n1 0\n"
    tests = append(tests, test{inp, solve(inp)})
    for len(tests) < 100 {
        k := rand.Intn(3) + 1 // up to 4 levels
        n := (1 << k) - 1
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", k))
        for i := 0; i < n; i++ {
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
        q := rand.Intn(5) + 1
        sb.WriteString(fmt.Sprintf("%d\n", q))
        for i := 0; i < q; i++ {
            p := rand.Intn(n) + 1
            r := rand.Intn(3)
            var ch byte
            if r == 0 {
                ch = '0'
            } else if r == 1 {
                ch = '1'
            } else {
                ch = '?'
            }
            sb.WriteString(fmt.Sprintf("%d %c\n", p, ch))
        }
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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

