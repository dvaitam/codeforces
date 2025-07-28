package main

import (
    "bytes"
    "context"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

type TestCase struct {
    n   int
    l   int64
    r   int64
    arr []int64
}

func countPairs(arr []int64, limit int64) int64 {
    var cnt int64
    j := len(arr) - 1
    for i := 0; i < len(arr); i++ {
        for j > i && arr[i]+arr[j] > limit {
            j--
        }
        if j <= i {
            break
        }
        cnt += int64(j - i)
    }
    return cnt
}

func solve(tc TestCase) int64 {
    b := append([]int64(nil), tc.arr...)
    sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
    return countPairs(b, tc.r) - countPairs(b, tc.l-1)
}

func generateTests() []TestCase {
    rand.Seed(44)
    tests := make([]TestCase, 100)
    for i := range tests {
        n := rand.Intn(20) + 1
        l := rand.Int63n(50) + 1
        r := l + rand.Int63n(50)
        arr := make([]int64, n)
        for j := range arr {
            arr[j] = rand.Int63n(100) + 1
        }
        tests[i] = TestCase{n: n, l: l, r: r, arr: arr}
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
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for idx, tc := range tests {
        input := fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.l, tc.r)
        for i, v := range tc.arr {
            if i+1 == tc.n {
                input += fmt.Sprintf("%d\n", v)
            } else {
                input += fmt.Sprintf("%d ", v)
            }
        }
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

