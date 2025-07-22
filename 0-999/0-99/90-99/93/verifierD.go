package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

const MOD int64 = 1000000007
const inv2 int64 = 500000004

func compute(L, R int64) int64 {
    allowedAdj := func(i, j int) bool {
        if i == j {
            return false
        }
        if (i == 0 && j == 1) || (i == 1 && j == 0) {
            return false
        }
        if (i == 2 && j == 3) || (i == 3 && j == 2) {
            return false
        }
        return true
    }
    forbiddenTriple := map[[3]int]bool{
        {3, 0, 2}: true,
        {2, 0, 3}: true,
    }
    var pair [][2]int
    idx := make(map[[2]int]int)
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            if allowedAdj(i, j) {
                id := len(pair)
                pair = append(pair, [2]int{i, j})
                idx[[2]int{i, j}] = id
            }
        }
    }
    m := len(pair)
    T := make([][]int64, m)
    for i := range T {
        T[i] = make([]int64, m)
    }
    for a := 0; a < m; a++ {
        i1, i2 := pair[a][0], pair[a][1]
        for k := 0; k < 4; k++ {
            if !allowedAdj(i2, k) {
                continue
            }
            if forbiddenTriple[[3]int{i1, i2, k}] {
                continue
            }
            b := idx[[2]int{i2, k}]
            T[a][b] = 1
        }
    }
    size := m + 1
    M := make([][]int64, size)
    for i := range M {
        M[i] = make([]int64, size)
    }
    for i := 0; i < m; i++ {
        for j := 0; j < m; j++ {
            M[i][j] = T[i][j]
        }
    }
    for j := 0; j < m; j++ {
        M[m][j] = 1
    }
    M[m][m] = 1
    v2 := make([]int64, size)
    for i := 0; i < m; i++ {
        v2[i] = 1
    }
    v2[m] = int64(m)
    mul := func(A, B [][]int64) [][]int64 {
        n := len(A)
        C := make([][]int64, n)
        for i := range C {
            C[i] = make([]int64, n)
        }
        for i := 0; i < n; i++ {
            for k := 0; k < n; k++ {
                if A[i][k] == 0 {
                    continue
                }
                aik := A[i][k]
                for j := 0; j < n; j++ {
                    C[i][j] = (C[i][j] + aik*B[k][j]) % MOD
                }
            }
        }
        return C
    }
    var pow func(mat [][]int64, e int64) [][]int64
    pow = func(mat [][]int64, e int64) [][]int64 {
        n := len(mat)
        res := make([][]int64, n)
        for i := 0; i < n; i++ {
            res[i] = make([]int64, n)
            res[i][i] = 1
        }
        base := mat
        for e > 0 {
            if e&1 == 1 {
                res = mul(res, base)
            }
            base = mul(base, base)
            e >>= 1
        }
        return res
    }
    sumA := func(b int64) int64 {
        if b < 2 {
            return 0
        }
        Mexp := pow(M, b-2)
        var s int64
        for j := 0; j < size; j++ {
            s = (s + v2[j]*Mexp[j][m]) % MOD
        }
        return s
    }
    var S1 int64
    if L <= 1 && R >= 1 {
        S1 = (S1 + 4) % MOD
    }
    l2 := L
    if l2 < 2 {
        l2 = 2
    }
    if R >= 2 && l2 <= R {
        S1 = (S1 + (sumA(R)-sumA(l2-1)+MOD)%MOD) % MOD
    }
    k1 := (L + 2) / 2
    if L%2 == 0 {
        k1 = L/2 + 1
    }
    k2 := (R + 1) / 2
    var S2 int64
    if k1 <= k2 {
        if k1 <= 1 && k2 >= 1 {
            S2 = (S2 + 4) % MOD
        }
        kk1 := k1
        if kk1 < 2 {
            kk1 = 2
        }
        if k2 >= 2 && kk1 <= k2 {
            S2 = (S2 + (sumA(k2)-sumA(kk1-1)+MOD)%MOD) % MOD
        }
    }
    ans := (S1 + S2) % MOD * inv2 % MOD
    return ans
}

type testCase struct {
    L, R int64
}

func runCase(bin string, tc testCase) error {
    input := fmt.Sprintf("%d %d\n", tc.L, tc.R)
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
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
    expected := compute(tc.L, tc.R)
    if val != expected {
        return fmt.Errorf("expected %d got %d", expected, val)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    var cases []testCase
    cases = append(cases, testCase{L: 1, R: 1})
    cases = append(cases, testCase{L: 1, R: 3})
    cases = append(cases, testCase{L: 123, R: 12345})
    for i := 0; i < 100; i++ {
        L := rng.Int63n(1_000_000_000) + 1
        R := L + rng.Int63n(1_000_000_000-L+1)
        cases = append(cases, testCase{L: L, R: R})
    }

    for i, tc := range cases {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %d %d\n", i+1, err, tc.L, tc.R)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

