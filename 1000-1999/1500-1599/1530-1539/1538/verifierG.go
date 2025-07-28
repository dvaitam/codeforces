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
    x int64
    y int64
    a int64
    b int64
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

func feasible(k, x, y, a, b int64) bool {
    if a < b {
        a, b = b, a
        x, y = y, x
    }
    if a == b {
        return k <= min(x, y)/a
    }
    diff := a - b
    if k*(a+b) > x+y {
        return false
    }
    lower := (k*a - y + diff - 1) / diff
    if lower < 0 {
        lower = 0
    }
    upper := (x - k*b) / diff
    if upper > k {
        upper = k
    }
    if upper < 0 {
        return false
    }
    return lower <= upper
}

func maxSets(x, y, a, b int64) int64 {
    if x < y {
        x, y = y, x
    }
    if a < b {
        a, b = b, a
    }
    if a == b {
        return min(x, y) / a
    }
    hi := (x + y) / (a + b)
    lo := int64(0)
    var res int64
    for lo <= hi {
        mid := (lo + hi) / 2
        if feasible(mid, x, y, a, b) {
            res = mid
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return res
}

func solve(tc TestCase) int64 {
    return maxSets(tc.x, tc.y, tc.a, tc.b)
}

func generateTests() []TestCase {
    rand.Seed(48)
    tests := make([]TestCase, 100)
    for i := range tests {
        x := rand.Int63n(100000) + 1
        y := rand.Int63n(100000) + 1
        a := rand.Int63n(1000) + 1
        b := rand.Int63n(1000) + 1
        tests[i] = TestCase{x: x, y: y, a: a, b: b}
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
        fmt.Println("usage: go run verifierG.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for idx, tc := range tests {
        input := fmt.Sprintf("1\n%d %d %d %d\n", tc.x, tc.y, tc.a, tc.b)
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

