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

func isGood(s string) bool {
    zeros := strings.Count(s, "0")
    ones := len(s) - zeros
    return zeros != ones
}

func expected(n int, s string) (int, []string) {
    zeros := strings.Count(s, "0")
    ones := n - zeros
    if n%2 == 1 || zeros != ones {
        return 1, []string{s}
    }
    return 2, []string{s[:1], s[1:]}
}

func run(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesA.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    idx := 0
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        idx++
        parts := strings.Fields(line)
        if len(parts) != 2 {
            fmt.Fprintf(os.Stderr, "invalid test case format on line %d\n", idx)
            os.Exit(1)
        }
        n, err := strconv.Atoi(parts[0])
        if err != nil {
            fmt.Fprintf(os.Stderr, "invalid n on line %d\n", idx)
            os.Exit(1)
        }
        s := parts[1]
        input := fmt.Sprintf("%d\n%s\n", n, s)
        out, err := run(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
            os.Exit(1)
        }
        fields := strings.Fields(out)
        if len(fields) < 2 {
            fmt.Fprintf(os.Stderr, "test %d: output too short\n", idx)
            os.Exit(1)
        }
        k, err := strconv.Atoi(fields[0])
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: invalid k\n", idx)
            os.Exit(1)
        }
        if len(fields) != 1+k {
            fmt.Fprintf(os.Stderr, "test %d: expected %d substrings got %d\n", idx, k, len(fields)-1)
            os.Exit(1)
        }
        substrs := fields[1 : 1+k]
        if strings.Join(substrs, "") != s {
            fmt.Fprintf(os.Stderr, "test %d: concatenation mismatch\n", idx)
            os.Exit(1)
        }
        for _, sub := range substrs {
            if !isGood(sub) {
                fmt.Fprintf(os.Stderr, "test %d: substring %s not good\n", idx, sub)
                os.Exit(1)
            }
        }
        expK, _ := expected(n, s)
        if k != expK {
            fmt.Fprintf(os.Stderr, "test %d: expected k=%d got %d\n", idx, expK, k)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}

