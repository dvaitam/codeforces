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

func expected(n, x int, a []int) []int {
    pref := make([]int, n+1)
    for i := 1; i <= n; i++ {
        pref[i] = pref[i-1] + a[i-1]
    }
    maxSum := make([]int, n+1)
    maxSum[0] = 0
    for l := 1; l <= n; l++ {
        best := -1 << 60
        for i := 0; i+l <= n; i++ {
            s := pref[i+l] - pref[i]
            if s > best {
                best = s
            }
        }
        maxSum[l] = best
    }
    res := make([]int, n+1)
    for k := 0; k <= n; k++ {
        ans := 0
        for l := 0; l <= n; l++ {
            add := k
            if l < k {
                add = l
            }
            val := maxSum[l] + add*x
            if val > ans {
                ans = val
            }
        }
        res[k] = ans
    }
    return res
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesC.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not open testcasesC.txt:", err)
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
        if len(parts) < 2 {
            fmt.Fprintf(os.Stderr, "bad line in testcase %d\n", caseNum)
            os.Exit(1)
        }
        n, _ := strconv.Atoi(parts[0])
        x, _ := strconv.Atoi(parts[1])
        if len(parts)-2 != n {
            fmt.Fprintf(os.Stderr, "testcase %d: expected %d numbers got %d\n", caseNum, n, len(parts)-2)
            os.Exit(1)
        }
        arr := make([]int, n)
        for i := 0; i < n; i++ {
            arr[i], _ = strconv.Atoi(parts[2+i])
        }
        exp := expected(n, x, arr)
        var input strings.Builder
        input.WriteString("1\n")
        input.WriteString(fmt.Sprintf("%d %d\n", n, x))
        for i := 0; i < n; i++ {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(fmt.Sprintf("%d", arr[i]))
        }
        input.WriteByte('\n')
        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(input.String())
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, out.String())
            os.Exit(1)
        }
        gotParts := strings.Fields(strings.TrimSpace(out.String()))
        if len(gotParts) != n+1 {
            fmt.Fprintf(os.Stderr, "case %d: expected %d numbers got %d\n", caseNum, n+1, len(gotParts))
            os.Exit(1)
        }
        for i := 0; i <= n; i++ {
            val, err := strconv.Atoi(gotParts[i])
            if err != nil || val != exp[i] {
                fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", caseNum, exp, gotParts)
                os.Exit(1)
            }
        }
    }
    fmt.Printf("All %d tests passed\n", t)
}

