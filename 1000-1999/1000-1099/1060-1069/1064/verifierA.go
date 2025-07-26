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

func runCandidate(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func solve(a, b, c int) int {
    if a > b {
        a, b = b, a
    }
    if b > c {
        b, c = c, b
    }
    if a > b {
        a, b = b, a
    }
    if a+b > c {
        return 0
    }
    return c - (a + b) + 1
}

func generateCase(rng *rand.Rand) (int, int, int) {
    return rng.Intn(100) + 1, rng.Intn(100) + 1, rng.Intn(100) + 1
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    tests := []struct{ a, b, c int }{
        {1, 1, 1},
        {1, 1, 2},
        {1, 2, 3},
        {2, 3, 5},
        {2, 3, 4},
        {100, 100, 100},
        {1, 100, 1},
        {50, 50, 100},
        {1, 1, 100},
        {3, 4, 5},
    }
    for i := 0; i < 100; i++ {
        a, b, c := generateCase(rng)
        tests = append(tests, struct{ a, b, c int }{a, b, c})
    }

    for idx, tc := range tests {
        input := fmt.Sprintf("%d %d %d\n", tc.a, tc.b, tc.c)
        expected := fmt.Sprintf("%d", solve(tc.a, tc.b, tc.c))
        out, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
            os.Exit(1)
        }
        if strings.TrimSpace(out) != expected {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, expected, out, input)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

