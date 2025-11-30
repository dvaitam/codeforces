package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

const (
    mod   int64 = 998244353
    maxN        = 500000 + 5
)

// Base64-encoded contents of testcasesD.txt.
const testcasesD64 = "MiAxCjIgMAo2IDEgMCAwIDEgMQozIDAgMAozIDAgMAo2IDEgMCAxIDEgMAozIDAgMQo0IDAgMCAxCjMgMCAxCjUgMSAwIDEgMAozIDAgMQozIDAgMAoyIDEKNiAwIDAgMCAwIDAKMyAxIDAKNSAxIDEgMCAwCjQgMCAxIDAKMyAxIDEKMiAxCjUgMCAxIDEgMAo1IDAgMSAxIDAKNCAxIDEgMAo2IDEgMCAxIDAgMQo1IDAgMCAwIDEKMiAxCjUgMSAwIDEgMQozIDAgMQo1IDEgMCAxIDAKNCAxIDAgMQo1IDAgMCAxIDAKNiAxIDEgMSAwIDEKNiAxIDAgMSAxIDEKMyAxIDAKMyAwIDEKMyAwIDEKNiAxIDEgMCAwIDEKNiAwIDAgMCAwIDAKMyAwIDAKNSAxIDEgMSAwCjUgMSAwIDEgMQo2IDAgMSAxIDAgMQo0IDEgMSAxCjUgMSAwIDAgMAoyIDAKMyAwIDEKNiAwIDEgMCAwIDAKMiAxCjUgMSAwIDEgMQo2IDEgMSAwIDAgMAo2IDAgMCAwIDEgMAo1IDEgMSAwIDAKNSAwIDEgMCAxCjIgMQo0IDAgMSAwCjUgMCAxIDEgMQo1IDEgMCAwIDEKMyAxIDAKNSAwIDAgMSAxCjYgMSAwIDEgMSAwCjIgMQozIDEgMAozIDAgMAo2IDAgMCAwIDEgMQoyIDEKNiAxIDEgMCAwIDAKMyAxIDEKNCAwIDAgMQo0IDEgMCAwCjQgMSAwIDEKNSAwIDAgMSAwCjUgMSAxIDEgMAoyIDEKNCAxIDAgMQo0IDEgMSAxCjUgMSAwIDEgMQoyIDEKMyAxIDAKMyAxIDEKMyAwIDEKNSAxIDEgMCAxCjYgMSAxIDAgMSAwCjIgMQo0IDEgMCAxCjIgMAo1IDAgMSAwIDEKMyAwIDAKNCAwIDAgMAoyIDEKNiAwIDAgMCAwIDAKNSAwIDAgMSAwCjYgMCAwIDAgMSAxCjUgMCAxIDAgMQoyIDEKMiAxCjQgMSAxIDAKNiAwIDEgMCAwIDEKMiAwCjYgMCAxIDAgMSAwCjUgMCAwIDEgMAo2IDEgMSAwIDEgMQo="

var (
    inv  []int64
    fact []int64
)

func precompute() {
    if len(inv) > 0 {
        return
    }
    inv = make([]int64, maxN)
    fact = make([]int64, maxN)
    inv[1] = 1
    for i := 2; i < maxN; i++ {
        inv[i] = mod - (mod/int64(i))*inv[int(mod%int64(i))]%mod
    }
    fact[0] = 1
    for i := 1; i < maxN; i++ {
        fact[i] = (fact[i-1] * int64(i)) % mod
    }
}

func modAdd(a, b int64) int64 {
    x := a + b
    if x >= mod {
        x -= mod
    }
    return x
}

func modMul(a, b int64) int64 { return (a * b) % mod }

type testCase struct {
    n int
    a []int
}

// solve implements the logic from 1806D.go.
func solve(tc testCase) string {
    precompute()
    n := tc.n
    a := tc.a
    pref := int64(1)
    sum := int64(0)
    var out bytes.Buffer
    for k := 1; k <= n-1; k++ {
        if a[k] == 1 {
            pref = modMul(pref, modMul(int64(k-1), inv[k]))
        }
        if a[k] == 0 {
            add := modMul(inv[k], pref)
            sum = modAdd(sum, add)
        }
        ans := modMul(fact[k], sum)
        if k > 1 {
            out.WriteByte(' ')
        }
        out.WriteString(strconv.FormatInt(ans, 10))
    }
    return out.String()
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesD64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    cases := make([]testCase, 0)
    for {
        if !sc.Scan() {
            break
        }
        nVal, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("parse n: %v", err)
        }
        a := make([]int, nVal)
        for i := 1; i <= nVal-1; i++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("truncated testcase with n=%d", nVal)
            }
            v, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("parse a[%d]: %v", i, err)
            }
            a[i] = v
        }
        cases = append(cases, testCase{n: nVal, a: a})
    }
    if len(cases) == 0 {
        return nil, fmt.Errorf("no testcases found")
    }
    return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "bin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    bin, cleanup, err := buildIfGo(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer cleanup()

    for idx, tc := range cases {
        var input strings.Builder
        input.WriteString("1\n")
        input.WriteString(fmt.Sprintf("%d\n", tc.n))
        for i := 1; i <= tc.n-1; i++ {
            if i > 1 {
                input.WriteByte(' ')
            }
            input.WriteString(strconv.Itoa(tc.a[i]))
        }
        input.WriteByte('\n')

        want := solve(tc)
        got, err := runCandidate(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

