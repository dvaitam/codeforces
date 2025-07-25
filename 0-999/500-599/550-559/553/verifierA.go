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

func parseTestcases(path string) ([]testCaseA, error) {
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
    cases, err := parseTestcases("testcasesA.txt")
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

