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

const mod int64 = 998244353

func expected(n, m, k, q int, xs, ys []int) int64 {
    rowUsed := make([]bool, n+1)
    colUsed := make([]bool, m+1)
    cntRow, cntCol := 0, 0
    ans := int64(1)
    for i := q - 1; i >= 0; i-- {
        if cntRow == n || cntCol == m {
            break
        }
        x, y := xs[i], ys[i]
        if !rowUsed[x] || !colUsed[y] {
            ans = ans * int64(k) % mod
            if !rowUsed[x] {
                rowUsed[x] = true
                cntRow++
            }
            if !colUsed[y] {
                colUsed[y] = true
                cntCol++
            }
        }
    }
    return ans
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesD.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not open testcasesD.txt:", err)
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
        if len(parts) < 4 {
            fmt.Fprintf(os.Stderr, "bad line in testcase %d\n", caseNum)
            os.Exit(1)
        }
        n, _ := strconv.Atoi(parts[0])
        m, _ := strconv.Atoi(parts[1])
        k, _ := strconv.Atoi(parts[2])
        q, _ := strconv.Atoi(parts[3])
        if len(parts) != 4+2*q {
            fmt.Fprintf(os.Stderr, "case %d: expected %d numbers got %d\n", caseNum, 4+2*q, len(parts))
            os.Exit(1)
        }
        xs := make([]int, q)
        ys := make([]int, q)
        for i := 0; i < q; i++ {
            xs[i], _ = strconv.Atoi(parts[4+2*i])
            ys[i], _ = strconv.Atoi(parts[4+2*i+1])
        }
        exp := expected(n, m, k, q, xs, ys)
        var input strings.Builder
        input.WriteString("1\n")
        input.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, k, q))
        for i := 0; i < q; i++ {
            input.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
        }
        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(input.String())
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, out.String())
            os.Exit(1)
        }
        got := strings.TrimSpace(out.String())
        val, err := strconv.ParseInt(got, 10, 64)
        if err != nil || val%mod != exp%mod {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", caseNum, exp%mod, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", t)
}

