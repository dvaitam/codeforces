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

const mod = 1000000009

func expected(n int64) string {
    n %= mod
    return fmt.Sprintf("%d", (n*n)%mod)
}

func genCase(rng *rand.Rand) (string, string) {
    n := int64(rng.Intn(1000000) + 2)
    input := fmt.Sprintf("%d\n", n)
    return input, expected(n)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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

