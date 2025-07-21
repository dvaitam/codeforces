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

func run(bin, input string) (string, error) {
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
    return strings.TrimSpace(out.String()), nil
}

func expected(n, m, x1, y1, x2, y2 int64) string {
    dx := x1 - x2
    if dx < 0 {
        dx = -dx
    }
    dy := y1 - y2
    if dy < 0 {
        dy = -dy
    }
    W := n - dx
    H := m - dy
    overlapW := n - 2*dx
    if overlapW < 0 {
        overlapW = 0
    }
    overlapH := m - 2*dy
    if overlapH < 0 {
        overlapH = 0
    }
    melted := 2*W*H - overlapW*overlapH
    total := n * m
    unmelted := total - melted
    return fmt.Sprintf("%d", unmelted)
}

func genCase(rng *rand.Rand) (string, string) {
    n := int64(rng.Intn(100) + 2)
    m := int64(rng.Intn(100) + 2)
    x1 := int64(rng.Intn(int(n))) + 1
    y1 := int64(rng.Intn(int(m))) + 1
    for {
        x2 := int64(rng.Intn(int(n))) + 1
        y2 := int64(rng.Intn(int(m))) + 1
        if x2 != x1 || y2 != y1 {
            exp := expected(n, m, x1, y1, x2, y2)
            input := fmt.Sprintf("1\n%d %d %d %d %d %d\n", n, m, x1, y1, x2, y2)
            return input, exp
        }
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := genCase(rng)
        out, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        if out != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

