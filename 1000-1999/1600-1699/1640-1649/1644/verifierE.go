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

func expected(n int64, s string) int64 {
    hasR, hasD := false, false
    for _, ch := range s {
        if ch == 'R' { hasR = true }
        if ch == 'D' { hasD = true }
    }
    if !hasR || !hasD {
        return n
    }
    posR, posD := int64(-1), int64(-1)
    for i, ch := range s {
        if ch == 'R' && posR == -1 { posR = int64(i+1) }
        if ch == 'D' && posD == -1 { posD = int64(i+1) }
    }
    var c int64
    if posR < posD { c = posD - 1 } else { c = posR - 1 }
    return n*n - int64(n-1)*c
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesE.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not open testcasesE.txt:", err)
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
    for caseNum := 1; caseNum <= t; caseNum++ {
        if !scanner.Scan() {
            fmt.Fprintf(os.Stderr, "expected %d cases, got %d\n", t, caseNum-1)
            os.Exit(1)
        }
        parts := strings.Fields(scanner.Text())
        if len(parts) != 2 {
            fmt.Fprintf(os.Stderr, "bad line in testcase %d\n", caseNum)
            os.Exit(1)
        }
        nVal, _ := strconv.ParseInt(parts[0], 10, 64)
        s := parts[1]
        exp := expected(nVal, s)
        input := fmt.Sprintf("1\n%d %s\n", nVal, s)
        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(input)
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, out.String())
            os.Exit(1)
        }
        gotStr := strings.TrimSpace(out.String())
        val, err := strconv.ParseInt(gotStr, 10, 64)
        if err != nil || val != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", caseNum, exp, gotStr)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", t)
}

