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

func apply(n int, states []string, keys []int) []string {
    wordIdx := make(map[string]int)
    words := []string{}
    init := make([]bool, n)
    for i, s := range states {
        if _, ok := wordIdx[s]; !ok {
            wordIdx[s] = len(words)
            words = append(words, s)
        }
        init[i] = wordIdx[s] == 1
    }
    parity := make([]bool, n+1)
    for _, k := range keys {
        if k >= 1 && k <= n {
            parity[k] = !parity[k]
        }
    }
    for i := 1; i <= n; i++ {
        if !parity[i] {
            continue
        }
        for j := i; j <= n; j += i {
            init[j-1] = !init[j-1]
        }
    }
    res := make([]string, n)
    for i, b := range init {
        if b {
            if len(words) > 1 {
                res[i] = words[1]
            } else {
                res[i] = words[0]
            }
        } else {
            res[i] = words[0]
        }
    }
    return res
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(5))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(6) + 1
        words := []string{"on", "off"}
        states := make([]string, n)
        for j := 0; j < n; j++ {
            states[j] = words[r.Intn(2)]
        }
        k := r.Intn(n) + 1
        keys := make([]int, k)
        for j := 0; j < k; j++ {
            keys[j] = r.Intn(n) + 1
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for j, s := range states {
            sb.WriteString(s)
            if j+1 < n {
                sb.WriteByte(' ')
            }
        }
        sb.WriteString("\n")
        sb.WriteString(fmt.Sprintf("%d\n", k))
        for j, x := range keys {
            sb.WriteString(fmt.Sprintf("%d", x))
            if j+1 < k {
                sb.WriteByte(' ')
            }
        }
        sb.WriteString("\n")
        final := apply(n, states, keys)
        ans := strings.Join(final, " ")
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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

