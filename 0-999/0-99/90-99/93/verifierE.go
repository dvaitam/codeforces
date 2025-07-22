package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

func compute(n uint64, arr []uint64) int64 {
    sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
    if len(arr) > 0 && arr[0] == 1 {
        return 0
    }
    var res int64
    var dfs func(idx int, prod uint64, depth int)
    dfs = func(idx int, prod uint64, depth int) {
        for i := idx; i < len(arr); i++ {
            ai := arr[i]
            if prod > n/ai {
                continue
            }
            newProd := prod * ai
            cnt := int64(n / newProd)
            if depth%2 == 0 {
                res += cnt
            } else {
                res -= cnt
            }
            dfs(i+1, newProd, depth+1)
        }
    }
    dfs(0, 1, 0)
    damage := int64(n) - res
    return damage
}

type testCase struct {
    n  uint64
    k  int
    a  []uint64
}

func runCase(bin string, tc testCase) error {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
    for i := 0; i < tc.k; i++ {
        if i > 0 {
            sb.WriteByte(' ')
        }
        sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
    }
    sb.WriteByte('\n')
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(sb.String())
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    var val int64
    if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &val); err != nil {
        return fmt.Errorf("invalid output: %v", err)
    }
    expected := compute(tc.n, append([]uint64(nil), tc.a...))
    if val != expected {
        return fmt.Errorf("expected %d got %d", expected, val)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    var cases []testCase
    cases = append(cases, testCase{n: 1, k: 1, a: []uint64{2}})
    cases = append(cases, testCase{n: 10, k: 2, a: []uint64{2,3}})
    cases = append(cases, testCase{n: 1000, k: 3, a: []uint64{2,5,7}})
    for i := 0; i < 100; i++ {
        n := rng.Uint64()%1_000_000 + 1
        k := rng.Intn(5) + 1
        arr := make([]uint64, k)
        used := make(map[uint64]bool)
        for j := 0; j < k; j++ {
            for {
                v := rng.Uint64()%1000 + 2
                if !used[v] {
                    used[v] = true
                    arr[j] = v
                    break
                }
            }
        }
        cases = append(cases, testCase{n: n, k: k, a: arr})
    }

    for i, tc := range cases {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

