package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func expected(s string) string {
    keys := make(map[rune]bool)
    for _, ch := range s {
        if ch >= 'a' && ch <= 'z' {
            keys[ch] = true
        } else {
            if !keys[ch+32] { // convert to lowercase
                return "NO"
            }
        }
    }
    return "YES"
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesA.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not open testcasesA.txt:", err)
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
        s := strings.TrimSpace(scanner.Text())
        in := fmt.Sprintf("1\n%s\n", s)
        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(in)
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i, err, out.String())
            os.Exit(1)
        }
        got := strings.TrimSpace(out.String())
        exp := expected(s)
        if got != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i, exp, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", t)
}

