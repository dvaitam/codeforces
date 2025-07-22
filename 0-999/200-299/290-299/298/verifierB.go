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

func solve(t int, sx, sy, ex, ey int64, wind string) string {
    dx := ex - sx
    dy := ey - sy
    for i := 0; i < t; i++ {
        switch wind[i] {
        case 'E':
            if dx > 0 {
                dx--
            }
        case 'W':
            if dx < 0 {
                dx++
            }
        case 'N':
            if dy > 0 {
                dy--
            }
        case 'S':
            if dy < 0 {
                dy++
            }
        }
        if dx == 0 && dy == 0 {
            return fmt.Sprintf("%d\n", i+1)
        }
    }
    return "-1\n"
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
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return out.String(), nil
}

func generateCase(rng *rand.Rand) (string, string) {
    t := rng.Intn(40) + 1
    sx := rng.Intn(21) - 10
    sy := rng.Intn(21) - 10
    ex := rng.Intn(21) - 10
    ey := rng.Intn(21) - 10
    if sx == ex && sy == ey {
        ex++
    }
    dirs := []byte{'E', 'S', 'W', 'N'}
    var sb strings.Builder
    for i := 0; i < t; i++ {
        sb.WriteByte(dirs[rng.Intn(4)])
    }
    wind := sb.String()
    input := fmt.Sprintf("%d %d %d %d %d\n%s\n", t, sx, sy, ex, ey, wind)
    expect := solve(t, int64(sx), int64(sy), int64(ex), int64(ey), wind)
    return input, expect
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := generateCase(rng)
        out, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        if strings.TrimSpace(out) != strings.TrimSpace(exp) {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(out), in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

