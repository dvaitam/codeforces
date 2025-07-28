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
    minIdx, maxIdx := 0, 0
    for i := 1; i < n; i++ {
        if arr[i] < arr[minIdx] {
            minIdx = i
        }
        if arr[i] > arr[maxIdx] {
            maxIdx = i
        }
    }
    minPos := minIdx + 1
    maxPos := maxIdx + 1
    if minPos > maxPos {
        minPos, maxPos = maxPos, minPos
    }
    fromLeft := maxPos
    fromRight := n - minPos + 1
    mixed1 := minPos + (n - maxPos + 1)
    mixed2 := maxPos + (n - minPos + 1)
    ans := fromLeft
    if fromRight < ans {
        ans = fromRight
    }
    if mixed1 < ans {
        ans = mixed1
    }
    if mixed2 < ans {
        ans = mixed2
    }
    return ans
}

func generateTests() []TestCase {
    rand.Seed(42)
    tests := make([]TestCase, 100)
    for i := range tests {
        n := rand.Intn(99) + 2 // 2..100
        perm := rand.Perm(n)
        arr := make([]int, n)
        for j, v := range perm {
            arr[j] = v + 1
        }
        tests[i] = TestCase{n: n, arr: arr}
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
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

