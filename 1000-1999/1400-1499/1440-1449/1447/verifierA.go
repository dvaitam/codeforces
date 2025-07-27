package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }

    binary := os.Args[1]
    rand.Seed(1)
    const t = 100
    ns := make([]int, t)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", t))
    for i := 0; i < t; i++ {
        n := rand.Intn(99) + 2 // 2..100
        ns[i] = n
        sb.WriteString(fmt.Sprintf("%d\n", n))
    }
    cmd := exec.Command(binary)
    cmd.Stdin = strings.NewReader(sb.String())
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Fprintf(os.Stderr, "binary execution error: %v\n", err)
        os.Exit(1)
    }
    reader := bufio.NewReader(bytes.NewReader(output))
    for idx, n := range ns {
        var m int
        if _, err := fmt.Fscan(reader, &m); err != nil {
            fmt.Printf("test %d: failed to read m: %v\n", idx+1, err)
            os.Exit(1)
        }
        if m < 1 || m > 1000 {
            fmt.Printf("test %d: invalid m=%d\n", idx+1, m)
            os.Exit(1)
        }
        ops := make([]int, m)
        for j := 0; j < m; j++ {
            if _, err := fmt.Fscan(reader, &ops[j]); err != nil {
                fmt.Printf("test %d: failed to read op %d: %v\n", idx+1, j+1, err)
                os.Exit(1)
            }
            if ops[j] < 1 || ops[j] > n {
                fmt.Printf("test %d: invalid bag index %d\n", idx+1, ops[j])
                os.Exit(1)
            }
        }
        candies := make([]int, n)
        for i := range candies {
            candies[i] = i + 1
        }
        for opIdx, bag := range ops {
            add := opIdx + 1
            for i := 0; i < n; i++ {
                if i+1 != bag {
                    candies[i] += add
                }
            }
        }
        expected := candies[0]
        for _, c := range candies {
            if c != expected {
                fmt.Printf("test %d: candies not equal\n", idx+1)
                os.Exit(1)
            }
        }
    }
    // ensure no extra output
    var extra string
    if _, err := fmt.Fscan(reader, &extra); err == nil {
        fmt.Println("extra output detected")
        os.Exit(1)
    }
    fmt.Println("All tests passed")
}

