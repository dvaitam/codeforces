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

func runCandidate(bin string, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func expected(n int, perm []int, queries []int) string {
    pos := make([]int, n+1)
    for i, v := range perm {
        pos[v] = i + 1
    }
    var vasya, petya int64
    for _, q := range queries {
        p := pos[q]
        vasya += int64(p)
        petya += int64(n - p + 1)
    }
    return fmt.Sprintf("%d %d", vasya, petya)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    f, err := os.Open("testcasesB.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    idx := 0
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        idx++
        fields := strings.Fields(line)
        if len(fields) < 3 {
            fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
            os.Exit(1)
        }
        n, _ := strconv.Atoi(fields[0])
        if len(fields) < 1+n+1 {
            fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
            os.Exit(1)
        }
        perm := make([]int, n)
        for i := 0; i < n; i++ {
            perm[i], _ = strconv.Atoi(fields[1+i])
        }
        mIndex := 1 + n
        m, _ := strconv.Atoi(fields[mIndex])
        if len(fields) != 1+n+1+m {
            fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
            os.Exit(1)
        }
        queries := make([]int, m)
        for i := 0; i < m; i++ {
            queries[i], _ = strconv.Atoi(fields[mIndex+1+i])
        }
        input := fmt.Sprintf("%d\n", n)
        for i, v := range perm {
            if i > 0 {
                input += " "
            }
            input += fmt.Sprintf("%d", v)
        }
        input += "\n"
        input += fmt.Sprintf("%d\n", m)
        for i, v := range queries {
            if i > 0 {
                input += " "
            }
            input += fmt.Sprintf("%d", v)
        }
        input += "\n"
        want := expected(n, perm, queries)
        got, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}

