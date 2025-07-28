package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func verifyOutput(n int, out string) error {
    lines := strings.Split(strings.TrimSpace(out), "\n")
    if len(lines) != n {
        return fmt.Errorf("expected %d permutations got %d", n, len(lines))
    }
    seen := make(map[string]bool)
    for i, line := range lines {
        fields := strings.Fields(line)
        if len(fields) != n {
            return fmt.Errorf("line %d: expected %d numbers got %d", i+1, n, len(fields))
        }
        perm := make([]int, n)
        used := make([]bool, n+1)
        for j, f := range fields {
            v, err := strconv.Atoi(f)
            if err != nil {
                return fmt.Errorf("line %d: invalid number", i+1)
            }
            if v < 1 || v > n || used[v] {
                return fmt.Errorf("line %d: invalid permutation", i+1)
            }
            used[v] = true
            perm[j] = v
            if j >= 2 && perm[j-2]+perm[j-1] == v {
                return fmt.Errorf("line %d: not anti-Fibonacci", i+1)
            }
        }
        key := strings.Join(fields, " ")
        if seen[key] {
            return fmt.Errorf("duplicate permutation on line %d", i+1)
        }
        seen[key] = true
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesB.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not open testcasesB.txt:", err)
        os.Exit(1)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    if !scanner.Scan() {
        fmt.Fprintln(os.Stderr, "empty testcases file")
        os.Exit(1)
    }
    var t int
    fmt.Sscan(scanner.Text(), &t)
    for i := 1; i <= t; i++ {
        if !scanner.Scan() {
            fmt.Fprintf(os.Stderr, "expected %d cases, got %d\n", t, i-1)
            os.Exit(1)
        }
        var n int
        fmt.Sscan(scanner.Text(), &n)
        in := fmt.Sprintf("1\n%d\n", n)
        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(in)
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i, err, out.String())
            os.Exit(1)
        }
        if err := verifyOutput(n, out.String()); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", t)
}

