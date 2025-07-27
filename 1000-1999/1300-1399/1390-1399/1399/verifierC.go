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

type testCaseC struct {
    w []int
}

func parseTestcasesC(path string) ([]testCaseC, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var cases []testCaseC
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
        if len(fields)-1 != n {
            return nil, fmt.Errorf("expected %d numbers", n)
        }
        arr := make([]int, n)
        for i := 0; i < n; i++ {
            v, err := strconv.Atoi(fields[i+1])
            if err != nil {
                return nil, err
            }
            arr[i] = v
        }
        cases = append(cases, testCaseC{arr})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cases, nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func solveC(w []int) int {
    n := len(w)
    freq := make([]int, n+1)
    for _, v := range w {
        if v >= 1 && v <= n {
            freq[v]++
        }
    }
    best := 0
    for s := 2; s <= 2*n; s++ {
        cnt := 0
        for i := 1; i < s-i && i <= n; i++ {
            j := s - i
            if j > n {
                continue
            }
            cnt += min(freq[i], freq[j])
        }
        if s%2 == 0 {
            i := s / 2
            if i >= 1 && i <= n {
                cnt += freq[i] / 2
            }
        }
        if cnt > best {
            best = cnt
        }
    }
    return best
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
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcasesC("testcasesC.txt")
    if err != nil {
        fmt.Println("failed to parse testcases:", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        n := len(tc.w)
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i, v := range tc.w {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(strconv.Itoa(v))
        }
        sb.WriteByte('\n')
        expected := strconv.Itoa(solveC(tc.w))
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

