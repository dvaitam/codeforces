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

const mod int64 = 1000000007

func powmod(base, exp int64) int64 {
    res := int64(1)
    base %= mod
    for exp > 0 {
        if exp&1 == 1 {
            res = res * base % mod
        }
        base = base * base % mod
        exp >>= 1
    }
    return res
}

func solveG(n, k int) string {
    mu := make([]int, k+1)
    lp := make([]int, k+1)
    primes := make([]int, 0)
    mu[1] = 1
    for i := 2; i <= k; i++ {
        if lp[i] == 0 {
            lp[i] = i
            primes = append(primes, i)
            mu[i] = -1
        }
        for _, p := range primes {
            if p > lp[i] || p*i > k {
                break
            }
            lp[p*i] = p
            if i%p == 0 {
                mu[p*i] = 0
                break
            } else {
                mu[p*i] = -mu[i]
            }
        }
    }
    powArr := make([]int64, k+1)
    for i := 1; i <= k; i++ {
        powArr[i] = powmod(int64(i), int64(n))
    }
    b := make([]int64, k+1)
    for d := 1; d <= k; d++ {
        if mu[d] == 0 {
            continue
        }
        md := int64(mu[d])
        for m := d; m <= k; m += d {
            b[m] += md * powArr[m/d]
        }
    }
    ans := int64(0)
    for i := 1; i <= k; i++ {
        bi := b[i] % mod
        if bi < 0 {
            bi += mod
        }
        ans = (ans + int64(int(bi)^i)) % mod
    }
    return fmt.Sprintf("%d", ans)
}

func generateG(rng *rand.Rand) (string, string) {
    n := rng.Intn(5) + 1
    k := rng.Intn(20) + 1
    input := fmt.Sprintf("%d %d\n", n, k)
    return input, solveG(n, k)
}

func runCase(bin, in, exp string) error {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    if got != exp {
        return fmt.Errorf("expected %s got %s", exp, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    for i := 0; i < 100; i++ {
        in, exp := generateG(rng)
        if err := runCase(bin, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

