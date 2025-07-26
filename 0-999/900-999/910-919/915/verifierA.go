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

type caseA struct {
    n, k int
    buckets []int
}

func computeA(c caseA) int {
    best := int(1<<31 - 1)
    for _, a := range c.buckets {
        if c.k%a == 0 {
            hours := c.k / a
            if hours < best {
                best = hours
            }
        }
    }
    return best
}

func generateA(rng *rand.Rand) (string, int) {
    n := rng.Intn(10) + 1
    k := rng.Intn(100) + 1
    buckets := make([]int, n)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
    divisor := rng.Intn(100) + 1
    for divisor == 0 || k%divisor != 0 {
        divisor = rng.Intn(100) + 1
    }
    pos := rng.Intn(n)
    for i := 0; i < n; i++ {
        if i == pos {
            buckets[i] = divisor
        } else {
            buckets[i] = rng.Intn(100) + 1
        }
        if i > 0 {
            sb.WriteByte(' ')
        }
        sb.WriteString(fmt.Sprintf("%d", buckets[i]))
    }
    sb.WriteByte('\n')
    return sb.String(), computeA(caseA{n: n, k: k, buckets: buckets})
}

func runCase(bin string, in string, exp int) error {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    var got int
    if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
        return fmt.Errorf("failed to parse output: %v", err)
    }
    if got != exp {
        return fmt.Errorf("expected %d got %d", exp, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    for i := 0; i < 100; i++ {
        in, exp := generateA(rng)
        if err := runCase(bin, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

