package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type testCaseB struct {
    n int
    k uint64
}

func parseTestcases(path string) ([]testCaseB, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    in := bufio.NewReader(f)
    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return nil, err
    }
    cases := make([]testCaseB, T)
    for i := 0; i < T; i++ {
        var n int
        var k uint64
        if _, err := fmt.Fscan(in, &n, &k); err != nil {
            return nil, err
        }
        cases[i] = testCaseB{n: n, k: k}
    }
    return cases, nil
}

func solveCase(tc testCaseB) string {
    n := tc.n
    k := tc.k
    fib := make([]uint64, n+2)
    fib[0] = 1
    fib[1] = 1
    for i := 2; i <= n; i++ {
        fib[i] = fib[i-1] + fib[i-2]
    }
    if k > fib[n] {
        k = fib[n]
    }
    res := make([]int, 0, n)
    i := 1
    for i <= n {
        count := fib[n-i]
        if k <= count {
            res = append(res, i)
            i++
        } else {
            k -= count
            res = append(res, i+1, i)
            i += 2
        }
    }
    var sb strings.Builder
    for idx, v := range res {
        if idx > 0 {
            sb.WriteByte(' ')
        }
        sb.WriteString(strconv.Itoa(v))
    }
    return sb.String()
}

func run(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcases("testcasesB.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
        expected := solveCase(tc)
        got, err := run(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != expected {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

