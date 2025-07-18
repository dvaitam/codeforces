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

func countFactor(x, p int) int {
    if x == 0 {
        return 1
    }
    c := 0
    for x%p == 0 {
        c++
        x /= p
    }
    return c
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func solve(mat [][]int) int {
    n := len(mat)
    dp2 := make([][]int, n)
    dp5 := make([][]int, n)
    for i := range dp2 {
        dp2[i] = make([]int, n)
        dp5[i] = make([]int, n)
    }
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            c2 := countFactor(mat[i][j], 2)
            c5 := countFactor(mat[i][j], 5)
            if i == 0 && j == 0 {
                dp2[i][j] = c2
                dp5[i][j] = c5
            } else if i == 0 {
                dp2[i][j] = dp2[i][j-1] + c2
                dp5[i][j] = dp5[i][j-1] + c5
            } else if j == 0 {
                dp2[i][j] = dp2[i-1][j] + c2
                dp5[i][j] = dp5[i-1][j] + c5
            } else {
                dp2[i][j] = min(dp2[i-1][j], dp2[i][j-1]) + c2
                dp5[i][j] = min(dp5[i-1][j], dp5[i][j-1]) + c5
            }
        }
    }
    best := min(dp2[n-1][n-1], dp5[n-1][n-1])
    // check path through zero
    start := make([][]bool, n)
    end := make([][]bool, n)
    for i := range start {
        start[i] = make([]bool, n)
        end[i] = make([]bool, n)
    }
    start[0][0] = true
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if i > 0 && start[i-1][j] {
                start[i][j] = true
            }
            if j > 0 && start[i][j-1] {
                start[i][j] = true
            }
        }
    }
    end[n-1][n-1] = true
    for i := n - 1; i >= 0; i-- {
        for j := n - 1; j >= 0; j-- {
            if i+1 < n && end[i+1][j] {
                end[i][j] = true
            }
            if j+1 < n && end[i][j+1] {
                end[i][j] = true
            }
        }
    }
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if mat[i][j] == 0 && start[i][j] && end[i][j] {
                if best > 1 {
                    best = 1
                }
            }
        }
    }
    return best
}

func generateCase(rng *rand.Rand) (string, int) {
    n := rng.Intn(6) + 2 // 2..7
    mat := make([][]int, n)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", n))
    for i := 0; i < n; i++ {
        mat[i] = make([]int, n)
        for j := 0; j < n; j++ {
            val := rng.Intn(1000)
            mat[i][j] = val
            sb.WriteString(fmt.Sprintf("%d ", val))
        }
        sb.WriteString("\n")
    }
    return sb.String(), solve(mat)
}

func runCase(exe string, input string, expected int) error {
    cmd := exec.Command(exe)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    var got int
    if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
        return fmt.Errorf("bad output: %v", err)
    }
    if got != expected {
        return fmt.Errorf("expected %d got %d", expected, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    exe := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := generateCase(rng)
        if err := runCase(exe, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

