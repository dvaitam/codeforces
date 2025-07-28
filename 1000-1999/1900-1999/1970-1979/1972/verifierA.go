package main

import (
    "bytes"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

// expectedOps computes the answer for problem A
func expectedOps(a, b []int) int {
    j := 0
    ops := 0
    n := len(a)
    for i := 0; i < n; i++ {
        if j < n && a[j] <= b[i] {
            j++
        } else {
            ops++
        }
    }
    return ops
}

func runTest(binary string, n int, a, b []int) error {
    var input strings.Builder
    input.WriteString("1\n")
    input.WriteString(fmt.Sprintf("%d\n", n))
    for i, v := range a {
        if i > 0 {
            input.WriteByte(' ')
        }
        input.WriteString(fmt.Sprintf("%d", v))
    }
    input.WriteByte('\n')
    for i, v := range b {
        if i > 0 {
            input.WriteByte(' ')
        }
        input.WriteString(fmt.Sprintf("%d", v))
    }
    input.WriteByte('\n')

    cmd := exec.Command(binary)
    cmd.Stdin = strings.NewReader(input.String())
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = io.Discard

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("execution error: %v", err)
    }
    result := strings.TrimSpace(out.String())
    expected := fmt.Sprintf("%d", expectedOps(a, b))
    if result != expected {
        return fmt.Errorf("expected %s, got %s", expected, result)
    }
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    binary := os.Args[1]
    rand.Seed(time.Now().UnixNano())
    const tests = 100
    for t := 0; t < tests; t++ {
        n := rand.Intn(100) + 1
        a := make([]int, n)
        b := make([]int, n)
        cur := rand.Intn(1000) + 1
        for i := 0; i < n; i++ {
            cur += rand.Intn(50) + 1
            if cur > 1_000_000_000 {
                cur = 1_000_000_000
            }
            a[i] = cur
        }
        cur = a[0] + rand.Intn(50)
        if cur > 1_000_000_000 {
            cur = 1_000_000_000
        }
        for i := 0; i < n; i++ {
            if cur < a[i] {
                cur = a[i]
            }
            cur += rand.Intn(50)
            if cur > 1_000_000_000 {
                cur = 1_000_000_000
            }
            b[i] = cur
        }
        if err := runTest(binary, n, a, b); err != nil {
            fmt.Printf("Test %d failed: %v\n", t+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", tests)
}

