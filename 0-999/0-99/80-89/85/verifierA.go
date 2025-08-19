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

func generateCase(rng *rand.Rand) string {
    n := rng.Intn(100) + 1
    return fmt.Sprintf("%d\n", n)
}

func runProg(exe, input string) (string, error) {
    cmd := exec.Command(exe)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func validate(n int, out string) bool {
    out = strings.TrimSpace(out)
    if out == "-1" {
        // For 4 x n, a valid tiling with the cut property exists for all n>=1 per oracle.
        return false
    }
    lines := strings.Split(out, "\n")
    // Remove any trailing empty lines
    var rows []string
    for _, ln := range lines {
        ln = strings.TrimRight(ln, "\r")
        if ln != "" {
            rows = append(rows, ln)
        }
    }
    if len(rows) != 4 {
        return false
    }
    grid := make([][]byte, 4)
    for i := 0; i < 4; i++ {
        if len(rows[i]) != n {
            return false
        }
        grid[i] = []byte(rows[i])
        for j := 0; j < n; j++ {
            if grid[i][j] < 'a' || grid[i][j] > 'z' {
                return false
            }
        }
    }
    // Each cell must be part of exactly one 1x2 or 2x1 domino of same color
    dirs := [][2]int{{1,0},{-1,0},{0,1},{0,-1}}
    for r := 0; r < 4; r++ {
        for c := 0; c < n; c++ {
            ch := grid[r][c]
            same := 0
            for _, d := range dirs {
                nr, nc := r+d[0], c+d[1]
                if nr >= 0 && nr < 4 && nc >= 0 && nc < n && grid[nr][nc] == ch {
                    same++
                }
            }
            if same != 1 {
                return false
            }
        }
    }
    // Every vertical cut must be crossed by at least one horizontal domino
    for c := 0; c < n-1; c++ {
        crossed := false
        for r := 0; r < 4; r++ {
            if grid[r][c] == grid[r][c+1] {
                crossed = true
                break
            }
        }
        if !crossed {
            return false
        }
    }
    return true
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        input := generateCase(rng)
        got, err := runProg(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
            os.Exit(1)
        }
        // parse n from input
        var n int
        fmt.Sscanf(strings.TrimSpace(input), "%d", &n)
        if !validate(n, got) {
            fmt.Fprintf(os.Stderr, "case %d mismatch\ninput:%s\noutput:\n%s\n", i+1, input, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
