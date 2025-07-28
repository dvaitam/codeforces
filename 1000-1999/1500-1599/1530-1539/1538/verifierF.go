package main

import (
    "bytes"
    "context"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

type TestCase struct {
    l int64
    r int64
}

func solve(tc TestCase) int64 {
    var ans int64
    for p := int64(1); p <= tc.r; p *= 10 {
        ans += tc.r/p - tc.l/p
    }
    return ans
}

func generateTests() []TestCase {
    rand.Seed(47)
    tests := make([]TestCase, 100)
    for i := range tests {
        l := rand.Int63n(1_000_000_000)
        r := l + rand.Int63n(1_000_000_000-l)+1
        tests[i] = TestCase{l: l, r: r}
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return "", err
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for idx, tc := range tests {
        input := fmt.Sprintf("1\n%d %d\n", tc.l, tc.r)
        want := fmt.Sprintf("%d", solve(tc))
        got, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("all tests passed")
}

