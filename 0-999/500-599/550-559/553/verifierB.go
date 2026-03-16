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

const testcasesRaw = `100
3 3
2 2
2 2
8 31
7 7
2 2
1 1
7 20
1 1
5 4
10 14
6 1
1 1
9 1
7 7
7 1
9 15
8 32
9 15
6 4
4 4
5 1
7 18
2 1
5 2
6 12
9 28
9 54
4 3
5 8
9 26
10 5
8 16
7 14
3 2
9 45
6 2
8 33
2 1
9 54
7 12
8 2
8 3
5 7
3 1
9 15
1 1
9 36
4 4
9 23
10 46
8 18
9 39
1 1
9 52
3 3
9 14
7 2
8 24
10 71
4 5
7 16
6 7
6 1
9 35
10 79
6 8
10 4
4 2
9 38
3 1
9 52
5 1
2 1
1 1
1 1
4 3
2 1
6 5
2 1
3 2
9 11
5 5
8 21
8 31
2 1
5 7
6 7
4 3
2 2
9 14
10 56
1 1
1 1
3 1
3 2
9 44
7 18
4 5
8 15`

func parseTestcases() ([]testCaseB, error) {
    in := bufio.NewReader(strings.NewReader(testcasesRaw))
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
    cases, err := parseTestcases()
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

