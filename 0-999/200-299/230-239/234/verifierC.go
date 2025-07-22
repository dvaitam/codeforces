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

type testC struct {
    n int
    t []int
}

func genTestsC() []testC {
    rng := rand.New(rand.NewSource(44))
    tests := make([]testC, 100)
    for i := range tests {
        n := rng.Intn(20) + 2
        arr := make([]int, n)
        for j := range arr {
            v := rng.Intn(21) - 10
            if v >= 0 {
                v++
            }
            arr[j] = v
        }
        tests[i] = testC{n: n, t: arr}
    }
    return tests
}

func solveC(tc testC) int {
    n := tc.n
    t := tc.t
    prefix := make([]int, n+1)
    for i := 1; i <= n; i++ {
        prefix[i] = prefix[i-1]
        if t[i-1] >= 0 {
            prefix[i]++
        }
    }
    suffix := make([]int, n+2)
    for i := n; i >= 1; i-- {
        suffix[i] = suffix[i+1]
        if t[i-1] <= 0 {
            suffix[i]++
        }
    }
    ans := n
    for k := 1; k <= n-1; k++ {
        cost := prefix[k] + suffix[k+1]
        if cost < ans {
            ans = cost
        }
    }
    return ans
}

func run(bin string, input string) (string, error) {
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
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return out.String(), nil
}

func runCase(bin string, tc testC) error {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", tc.n))
    for i, v := range tc.t {
        if i > 0 { sb.WriteByte(' ') }
        sb.WriteString(strconv.Itoa(v))
    }
    sb.WriteByte('\n')
    out, err := run(bin, sb.String())
    if err != nil { return err }
    expected := solveC(tc)
    valStr := strings.TrimSpace(out)
    val, err := strconv.Atoi(valStr)
    if err != nil { return fmt.Errorf("invalid output") }
    if val != expected { return fmt.Errorf("expected %d got %d", expected, val) }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsC()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

