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

func computeB(n, pos, l, r int) int {
    abs := func(x int) int { if x < 0 { return -x }; return x }
    min := func(a, b int) int { if a < b { return a }; return b }
    switch {
    case l == 1 && r == n:
        return 0
    case l == 1:
        return abs(pos-r) + 1
    case r == n:
        return abs(pos-l) + 1
    default:
        return min(abs(pos-l), abs(pos-r)) + (r - l) + 2
    }
}

func generateB(rng *rand.Rand) (string, int) {
    n := rng.Intn(100) + 1
    pos := rng.Intn(n) + 1
    l := rng.Intn(n) + 1
    r := rng.Intn(n-l+1) + l
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, pos, l, r))
    return sb.String(), computeB(n, pos, l, r)
}

func runCase(bin, in string, exp int) error {
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    for i := 0; i < 100; i++ {
        in, exp := generateB(rng)
        if err := runCase(bin, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

