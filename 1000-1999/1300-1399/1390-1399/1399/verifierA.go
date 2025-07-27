package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

type testCaseA struct {
    arr []int
}

func parseTestcasesA(path string) ([]testCaseA, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var cases []testCaseA
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
        cases = append(cases, testCaseA{arr})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cases, nil
}

func solveA(arr []int) string {
    sort.Ints(arr)
    for i := 1; i < len(arr); i++ {
        if arr[i]-arr[i-1] > 1 {
            return "NO"
        }
    }
    return "YES"
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcasesA("testcasesA.txt")
    if err != nil {
        fmt.Println("failed to parse testcases:", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
        for i, v := range tc.arr {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(strconv.Itoa(v))
        }
        sb.WriteByte('\n')
        expected := solveA(tc.arr)
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

