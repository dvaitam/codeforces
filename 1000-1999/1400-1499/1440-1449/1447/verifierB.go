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

func expected(grid [][]int64) int64 {
    var sumAbs int64
    countNeg := 0
    countZero := 0
    minAbs := int64(1<<62)
    for _, row := range grid {
        for _, v := range row {
            if v < 0 {
                countNeg++
            }
            if v == 0 {
                countZero++
            }
            if v < 0 {
                v = -v
            }
            if v < minAbs {
                minAbs = v
            }
            sumAbs += v
        }
    }
    if countNeg%2 != 0 && countZero == 0 {
        sumAbs -= 2 * minAbs
    }
    return sumAbs
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    binary := os.Args[1]
    rand.Seed(1)
    const t = 100
    type test struct {
        n, m int
        grid [][]int64
    }
    tests := make([]test, t)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", t))
    for i := 0; i < t; i++ {
        n := rand.Intn(9) + 2   //2..10
        m := rand.Intn(9) + 2
        g := make([][]int64, n)
        sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
        for r := 0; r < n; r++ {
            g[r] = make([]int64, m)
            for c := 0; c < m; c++ {
                val := int64(rand.Intn(201) - 100)
                g[r][c] = val
                sb.WriteString(fmt.Sprintf("%d ", val))
            }
            sb.WriteString("\n")
        }
        tests[i] = test{n: n, m: m, grid: g}
    }
    cmd := exec.Command(binary)
    cmd.Stdin = strings.NewReader(sb.String())
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Fprintf(os.Stderr, "binary execution error: %v\n", err)
        os.Exit(1)
    }
    reader := bufio.NewReader(bytes.NewReader(output))
    for i, tst := range tests {
        var got int64
        if _, err := fmt.Fscan(reader, &got); err != nil {
            fmt.Printf("test %d: failed to read output: %v\n", i+1, err)
            os.Exit(1)
        }
        exp := expected(tst.grid)
        if got != exp {
            fmt.Printf("test %d: expected %d, got %d\n", i+1, exp, got)
            os.Exit(1)
        }
    }
    var extra string
    if _, err := fmt.Fscan(reader, &extra); err == nil {
        fmt.Println("extra output detected")
        os.Exit(1)
    }
    fmt.Println("All tests passed")
}

