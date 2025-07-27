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

func expected(n, m, r, c int) int {
    d1 := (r - 1) + (c - 1)
    d2 := (r - 1) + (m - c)
    d3 := (n - r) + (c - 1)
    d4 := (n - r) + (m - c)
    ans := d1
    if d2 > ans {
        ans = d2
    }
    if d3 > ans {
        ans = d3
    }
    if d4 > ans {
        ans = d4
    }
    return ans
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
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
        if len(parts) != 4 {
            fmt.Printf("test %d: invalid line\n", idx)
            os.Exit(1)
        }
        n, _ := strconv.Atoi(parts[0])
        m, _ := strconv.Atoi(parts[1])
        r, _ := strconv.Atoi(parts[2])
        c, _ := strconv.Atoi(parts[3])
        expect := expected(n, m, r, c)
        input := fmt.Sprintf("1\n%d %d %d %d\n", n, m, r, c)
        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(input)
        var out bytes.Buffer
        var stderr bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &stderr
        err := cmd.Run()
        if err != nil {
            fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
            os.Exit(1)
        }
        gotStr := strings.TrimSpace(out.String())
        got, err := strconv.Atoi(gotStr)
        if err != nil {
            fmt.Printf("test %d: cannot parse output %q\n", idx, gotStr)
            os.Exit(1)
        }
        if got != expect {
            fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}
