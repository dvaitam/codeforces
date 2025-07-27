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

type testCaseD struct {
    s string
}

func parseTestcasesD(path string) ([]testCaseD, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var cases []testCaseD
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        parts := strings.Fields(line)
        if len(parts) != 2 {
            return nil, fmt.Errorf("bad line: %s", line)
        }
        n, err := strconv.Atoi(parts[0])
        if err != nil {
            return nil, err
        }
        if n != len(parts[1]) {
            return nil, fmt.Errorf("length mismatch")
        }
        cases = append(cases, testCaseD{parts[1]})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cases, nil
}

func solveD(s string) (int, []int) {
    var stack0, stack1 []int
    res := make([]int, len(s))
    cnt := 0
    for i := 0; i < len(s); i++ {
        bit := s[i] - '0'
        var id int
        if bit == 0 {
            if len(stack1) > 0 {
                id = stack1[len(stack1)-1]
                stack1 = stack1[:len(stack1)-1]
            } else {
                cnt++
                id = cnt
            }
            stack0 = append(stack0, id)
        } else {
            if len(stack0) > 0 {
                id = stack0[len(stack0)-1]
                stack0 = stack0[:len(stack0)-1]
            } else {
                cnt++
                id = cnt
            }
            stack1 = append(stack1, id)
        }
        res[i] = id
    }
    return cnt, res
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
        fmt.Println("usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcasesD("testcasesD.txt")
    if err != nil {
        fmt.Println("failed to parse testcases:", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d %s\n", len(tc.s), tc.s))
        k, seq := solveD(tc.s)
        expected := fmt.Sprintf("%d\n", k)
        for i, v := range seq {
            if i > 0 {
                expected += " "
            }
            expected += strconv.Itoa(v)
        }
        got, err := run(bin, sb.String())
        if err != nil {
            fmt.Printf("case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(expected) {
            fmt.Printf("case %d failed:\nexpected:%s\ngot:%s\n", idx+1, expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

