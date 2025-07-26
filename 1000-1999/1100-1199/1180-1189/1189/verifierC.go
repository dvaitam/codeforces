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

func candies(seq []int) int {
    cnt := 0
    for len(seq) > 1 {
        next := make([]int, len(seq)/2)
        for i := 0; i < len(seq); i += 2 {
            sum := seq[i] + seq[i+1]
            if sum >= 10 {
                cnt++
            }
            next[i/2] = sum % 10
        }
        seq = next
    }
    return cnt
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesC.txt")
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
        if len(parts) < 3 {
            fmt.Fprintf(os.Stderr, "invalid test case on line %d\n", idx)
            os.Exit(1)
        }
        pos := 0
        n, err := strconv.Atoi(parts[pos])
        if err != nil || n <= 0 {
            fmt.Fprintf(os.Stderr, "invalid n on line %d\n", idx)
            os.Exit(1)
        }
        pos++
        if len(parts) < pos+n+1 {
            fmt.Fprintf(os.Stderr, "invalid line length on line %d\n", idx)
            os.Exit(1)
        }
        digits := make([]int, n)
        for i := 0; i < n; i++ {
            v, err := strconv.Atoi(parts[pos+i])
            if err != nil {
                fmt.Fprintf(os.Stderr, "invalid digit on line %d\n", idx)
                os.Exit(1)
            }
            digits[i] = v
        }
        pos += n
        q, err := strconv.Atoi(parts[pos])
        if err != nil {
            fmt.Fprintf(os.Stderr, "invalid q on line %d\n", idx)
            os.Exit(1)
        }
        pos++
        if len(parts) != pos+2*q {
            fmt.Fprintf(os.Stderr, "invalid query count on line %d\n", idx)
            os.Exit(1)
        }
        queries := make([][2]int, q)
        for i := 0; i < q; i++ {
            l, _ := strconv.Atoi(parts[pos])
            r, _ := strconv.Atoi(parts[pos+1])
            queries[i] = [2]int{l, r}
            pos += 2
        }
        var sb strings.Builder
        fmt.Fprintf(&sb, "%d\n", n)
        for i := 0; i < n; i++ {
            if i > 0 {
                sb.WriteByte(' ')
            }
            fmt.Fprintf(&sb, "%d", digits[i])
        }
        sb.WriteString("\n")
        fmt.Fprintf(&sb, "%d\n", q)
        expected := make([]int, q)
        for i, qr := range queries {
            fmt.Fprintf(&sb, "%d %d\n", qr[0], qr[1])
            seg := make([]int, qr[1]-qr[0]+1)
            copy(seg, digits[qr[0]-1:qr[1]])
            expected[i] = candies(seg)
        }
        out, err := run(bin, sb.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
            os.Exit(1)
        }
        outLines := strings.Fields(out)
        if len(outLines) != q {
            fmt.Fprintf(os.Stderr, "test %d: expected %d lines got %d\n", idx, q, len(outLines))
            os.Exit(1)
        }
        for i := 0; i < q; i++ {
            val, err := strconv.Atoi(outLines[i])
            if err != nil || val != expected[i] {
                fmt.Fprintf(os.Stderr, "test %d query %d failed: expected %d got %s\n", idx, i+1, expected[i], outLines[i])
                os.Exit(1)
            }
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}

