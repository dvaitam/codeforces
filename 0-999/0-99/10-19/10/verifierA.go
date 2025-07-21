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

type interval struct {
    l, r int
}

func compute(n, p1, p2, p3, t1, t2 int, ivs []interval) int64 {
    pc := int64(0)
    pc += int64(ivs[0].r-ivs[0].l) * int64(p1)
    prev := ivs[0].r
    for i := 1; i < n; i++ {
        gap := ivs[i].l - prev
        if gap <= t1 {
            pc += int64(gap) * int64(p1)
        } else {
            pc += int64(t1) * int64(p1)
            gap -= t1
            if gap <= t2 {
                pc += int64(gap) * int64(p2)
            } else {
                pc += int64(t2) * int64(p2)
                gap -= t2
                pc += int64(gap) * int64(p3)
            }
        }
        pc += int64(ivs[i].r-ivs[i].l) * int64(p1)
        prev = ivs[i].r
    }
    return pc
}

func generateCase(rng *rand.Rand) (string, int64) {
    n := rng.Intn(5) + 1
    p1 := rng.Intn(101)
    p2 := rng.Intn(101)
    p3 := rng.Intn(101)
    t1 := rng.Intn(60) + 1
    t2 := rng.Intn(60) + 1
    ivs := make([]interval, n)
    cur := rng.Intn(10)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d\n", n, p1, p2, p3, t1, t2))
    for i := 0; i < n; i++ {
        start := cur + rng.Intn(10)
        length := rng.Intn(10) + 1
        end := start + length
        if end > 1440 {
            end = 1440
        }
        ivs[i] = interval{start, end}
        sb.WriteString(fmt.Sprintf("%d %d\n", start, end))
        cur = end + rng.Intn(10)
    }
    ans := compute(n, p1, p2, p3, t1, t2, ivs)
    return sb.String(), ans
}

func runCase(exe, input string, expected int64) error {
    cmd := exec.Command(exe)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    outStr := strings.TrimSpace(out.String())
    var got int64
    fmt.Sscan(outStr, &got)
    if got != expected {
        return fmt.Errorf("expected %d got %s", expected, outStr)
    }
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    exe := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        input, expected := generateCase(rng)
        if err := runCase(exe, input, expected); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

