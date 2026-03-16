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

type testCaseA struct {
    k int
    c []int
}

const testcasesRaw = `100
7
7 1 3 5 1 1 1
5
12 5 1 1 1
5
1 5 1 1 1
5
1 3 2 1 1
6
3 5 1 1 1 1
9
1 5 1 1 1 1 1 1 1
4
2 4 5 3
4
3 5 2 1
2
2 2
9
2 3 5 1 1 1 1 1 1
8
2 1 1 1 1 1 1 1
1
20
5
3 2 11 1 3
3
1 2 1
9
5 1 1 1 1 1 1 1 1
2
10 2
2
11 6
4
1 5 1 4
4
3 6 4 2
1
4
3
1 2 7
1
4
4
4 1 1 1
2
1 2
4
6 1 1 1
4
1 2 1 1
5
1 1 3 1 1
6
2 1 5 2 1 1
10
2 1 1 1 1 1 1 1 1 1
3
9 2 2
2
11 5
3
1 1 1
10
5 3 2 2 1 1 1 1 1 1
1
18
5
2 4 1 1 1
6
10 1 1 1 1 1
2
1 1
6
1 1 3 1 1 1
10
1 4 4 1 1 1 1 1 1 1
8
8 1 1 2 1 1 1 1
1
14
4
3 1 1 1
7
7 3 1 1 1 3 1
1
17
10
1 1 1 2 1 1 1 1 1 1
2
1 1
8
13 1 1 1 1 1 1 1
5
1 1 1 1 1
9
3 1 4 1 1 1 1 1 1
7
5 1 2 2 3 2 1
6
6 2 3 5 1 1
1
9
3
5 1 1
6
5 1 2 1 2 1
4
2 1 1 1
5
6 5 4 1 1
2
16 3
8
3 1 4 1 1 1 1 1
5
11 1 1 1 1
2
1 3
4
1 7 1 2
2
9 5
9
4 2 1 1 1 1 1 1 1
10
2 1 2 1 1 1 1 1 1 1
8
3 2 1 2 3 1 1 1
4
7 6 2 3
2
5 12
10
6 2 1 1 1 1 1 1 1 1
6
12 1 1 1 1 1
5
4 2 1 1 1
2
2 16
2
3 6
1
10
1
15
6
3 3 8 2 2 1
9
1 2 2 6 2 1 1 1 1
8
7 5 3 1 1 1 1 1
5
6 1 4 1 1
5
7 7 1 1 2
3
5 3 2
1
2
8
2 4 3 1 1 1 1 1
1
20
8
4 1 1 4 1 1 1 1
10
3 6 1 3 1 1 1 1 1 1
1
20
8
6 1 5 1 1 1 1 1
2
7 11
1
13
8
2 4 2 1 1 1 1 1
5
1 1 1 1 1
9
1 2 4 1 1 1 1 1 1
2
13 2
3
1 1 1
7
1 1 1 1 1 1 1
2
12 1
1
12
6
1 1 2 1 1 2
4
1 1 1 1
5
12 1 1 1 1`

func parseTestcases() ([]testCaseA, error) {
    in := bufio.NewReader(strings.NewReader(testcasesRaw))
    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return nil, err
    }
    cases := make([]testCaseA, T)
    for i := 0; i < T; i++ {
        var k int
        if _, err := fmt.Fscan(in, &k); err != nil {
            return nil, err
        }
        c := make([]int, k)
        for j := 0; j < k; j++ {
            fmt.Fscan(in, &c[j])
        }
        cases[i] = testCaseA{k: k, c: c}
    }
    return cases, nil
}

const mod int64 = 1000000007

func modPow(a, e int64) int64 {
    res := int64(1)
    for e > 0 {
        if e&1 == 1 {
            res = res * a % mod
        }
        a = a * a % mod
        e >>= 1
    }
    return res
}

func modInv(a int64) int64 { return modPow(a, mod-2) }

func solveCase(tc testCaseA) string {
    sum := 0
    for _, v := range tc.c {
        sum += v
    }
    maxN := sum
    fact := make([]int64, maxN+1)
    invFact := make([]int64, maxN+1)
    fact[0] = 1
    for i := 1; i <= maxN; i++ {
        fact[i] = fact[i-1] * int64(i) % mod
    }
    invFact[maxN] = modInv(fact[maxN])
    for i := maxN; i > 0; i-- {
        invFact[i-1] = invFact[i] * int64(i) % mod
    }
    res := int64(1)
    pref := tc.c[0]
    for i := 1; i < tc.k; i++ {
        n := pref + tc.c[i] - 1
        choose := fact[n] * invFact[tc.c[i]-1] % mod * invFact[pref] % mod
        res = res * choose % mod
        pref += tc.c[i]
    }
    return strconv.FormatInt(res, 10)
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", tc.k))
        for j, v := range tc.c {
            if j > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(strconv.Itoa(v))
        }
        sb.WriteByte('\n')
        expected := solveCase(tc)
        got, err := run(bin, sb.String())
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

