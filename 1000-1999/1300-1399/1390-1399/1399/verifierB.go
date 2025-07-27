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

type testCaseB struct {
    a []int64
    b []int64
}

func parseTestcasesB(path string) ([]testCaseB, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var cases []testCaseB
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        fields := strings.Fields(line)
        if len(fields) < 1 {
            return nil, fmt.Errorf("bad line: %s", line)
        }
        n, err := strconv.Atoi(fields[0])
        if err != nil {
            return nil, err
        }
        if len(fields)-1 != 2*n {
            return nil, fmt.Errorf("expected %d numbers", 2*n)
        }
        a := make([]int64, n)
        b := make([]int64, n)
        for i := 0; i < n; i++ {
            v, err := strconv.ParseInt(fields[1+i], 10, 64)
            if err != nil {
                return nil, err
            }
            a[i] = v
        }
        for i := 0; i < n; i++ {
            v, err := strconv.ParseInt(fields[1+n+i], 10, 64)
            if err != nil {
                return nil, err
            }
            b[i] = v
        }
        cases = append(cases, testCaseB{a, b})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cases, nil
}

func solveB(a, b []int64) int64 {
    minA := a[0]
    minB := b[0]
    for i := 1; i < len(a); i++ {
        if a[i] < minA {
            minA = a[i]
        }
        if b[i] < minB {
            minB = b[i]
        }
    }
    var moves int64
    for i := 0; i < len(a); i++ {
        da := a[i] - minA
        db := b[i] - minB
        if da > db {
            moves += da
        } else {
            moves += db
        }
    }
    return moves
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
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcasesB("testcasesB.txt")
    if err != nil {
        fmt.Println("failed to parse testcases:", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        n := len(tc.a)
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i, v := range tc.a {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(strconv.FormatInt(v, 10))
        }
        sb.WriteByte('\n')
        for i, v := range tc.b {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(strconv.FormatInt(v, 10))
        }
        sb.WriteByte('\n')
        expected := strconv.FormatInt(solveB(tc.a, tc.b), 10)
        got, err := run(bin, sb.String())
        if err != nil {
            fmt.Printf("case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != expected {
            fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

