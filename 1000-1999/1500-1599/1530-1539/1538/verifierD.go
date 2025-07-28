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
    a int64
    b int64
    k int
}

func primeCount(x int64) int {
    cnt := 0
    for x%2 == 0 {
        cnt++
        x /= 2
    }
    for i := int64(3); i*i <= x; i += 2 {
        for x%i == 0 {
            cnt++
            x /= i
        }
    }
    if x > 1 {
        cnt++
    }
    return cnt
}

func solve(tc TestCase) string {
    cntA := primeCount(tc.a)
    cntB := primeCount(tc.b)
    maxOps := cntA + cntB
    if tc.k == 1 {
        if tc.a != tc.b && (tc.a%tc.b == 0 || tc.b%tc.a == 0) {
            return "Yes"
        }
        return "No"
    }
    if tc.k >= 2 && tc.k <= maxOps {
        return "Yes"
    }
    return "No"
}

func generateTests() []TestCase {
    rand.Seed(45)
    tests := make([]TestCase, 100)
    for i := range tests {
        a := rand.Int63n(1000000) + 1
        b := rand.Int63n(1000000) + 1
        k := rand.Intn(6) + 1
        tests[i] = TestCase{a: a, b: b, k: k}
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
        fmt.Println("usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for idx, tc := range tests {
        input := fmt.Sprintf("1\n%d %d %d\n", tc.a, tc.b, tc.k)
        want := solve(tc)
        got, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("all tests passed")
}

