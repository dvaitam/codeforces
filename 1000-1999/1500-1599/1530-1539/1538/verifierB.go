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
    n   int
    arr []int
}

func solve(n int, arr []int) int {
    sum := 0
    for _, v := range arr {
        sum += v
    }
    if sum%n != 0 {
        return -1
    }
    avg := sum / n
    ans := 0
    for _, v := range arr {
        if v > avg {
            ans++
        }
    }
    return ans
}

func generateTests() []TestCase {
    rand.Seed(43)
    tests := make([]TestCase, 100)
    for i := range tests {
        n := rand.Intn(20) + 1 // 1..20
        arr := make([]int, n)
        for j := range arr {
            arr[j] = rand.Intn(21) // 0..20
        }
        tests[i] = TestCase{n: n, arr: arr}
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
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for idx, tc := range tests {
        input := fmt.Sprintf("1\n%d\n", tc.n)
        for i, v := range tc.arr {
            if i+1 == tc.n {
                input += fmt.Sprintf("%d\n", v)
            } else {
                input += fmt.Sprintf("%d ", v)
            }
        }
        want := fmt.Sprintf("%d", solve(tc.n, tc.arr))
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

