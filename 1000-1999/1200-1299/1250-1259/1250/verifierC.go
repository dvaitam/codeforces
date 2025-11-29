package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

const testcases = `7 14 1 3 5 4 4 7 7 3 4 3 5 2 5 2
5 5 1 5 3 5 5
3 10 1 3 1 3 2 2 3 1 2 2
6 20 6 2 5 4 4 5 3 1 5 1 1 6 4 6 6 6 1 5 4 3
4 11 1 2 2 2 2 4 1 1 3 4 1
5 18 3 1 5 3 5 2 5 5 5 3 4 1 5 4 3 5 2 3
3 7 1 1 3 3 2 2 1
2 5 1 1 1 2 2
9 8 4 7 5 8 8 6 2 6
10 4 8 10 6 4
6 8 6 2 3 7 5 6 3 7
10 7 1 3 2 3 1 1 2 1 2 3 1
8 6 2 3 3 1 3 1 2 2
10 10 4 2 2 5 1 3 5 4 4 5 1
7 11 1 5 2 5 3 2 1 5 5 2 2
8 3 5 1 2 3 1 2 2 2
7 6 5 1 2 3 3 1 2 3
2 7 2 1 3 3 3 2
6 8 5 5 5 5 3 2 4 5
4 12 3 3 3 2 1 3 1 3 2 1 1
7 6 3 4 3 5 3 1 2 4
1 9 4 4 4 4 4 2 3
5 6 3 3 1 3 3 2
3 10 1 1 2 2 1 1 2 2 1 1
7 9 4 5 3 4 2 3 3 5 1
8 8 3 3 3 1 2 2 2 2
5 6 3 3 1 2 3 2
8 6 1 5 3 2 3 1 4 3
8 8 4 4 1 1 2 4 4 1
10 3 3 3 1 3 3 3 2 3 3 3
2 3 1 1 1
6 3 1 1 3 1 3 2
9 6 4 3 3 1 1 2 4 2 4
9 8 3 1 2 4 3 4 2 4 2
8 3 1 2 2 2 2 1 3 2
7 6 1 2 2 1 1 3 1
4 7 1 2 3 3 3 3 2
6 11 5 5 1 1 2 3 4 4 2 1 2
5 3 1 1 2 2 2
3 6 2 3 2 3 2
2 9 4 5 4 4 1 4 5 4 1
4 2 1 2 2 2
8 2 1 2 1 2 2 2 2 2
6 3 1 3 2 1 3 3
10 8 2 2 1 1 4 2 4 1 3 1
10 4 1 1 2 1 1 2 2 3 1
8 7 2 3 1 2 1 1 3 1 1
9 3 1 1 3 3 2 2 3 3 3
3 1 1 1 1
4 8 3 1 3 3 4 1 4 1
1 6 2 2 2 2 2 2
8 5 3 1 2 3 2 2 2 3 1
8 5 3 1 3 1 1 3 1 3
6 2 1 1 3 3 3
7 2 2 3 2 1 1 1 2
4 1 1 1 1 1
8 1 2 1 2 2 2 1 1 2
5 1 1 2 1 2 2
10 3 3 2 3 2 2 3 3 3 1 3
2 6 2 2 1 2 2 1
3 2 2 2 2 2 2
9 2 1 2 2 2 1 2 2 2 2
2 10 1 1 1 1 1 1 1 1 1 1
3 6 2 1 1 2 2 2
8 8 1 3 4 2 1 1 4 4 3
8 5 4 2 1 1 2 2 1 2
4 1 1 1 1 1
4 6 4 1 2 1 2 1
9 1 2 1 1 1 1 1 1 1 1
4 8 2 1 3 3 2 2 2
8 2 2 2 1 2 1 1 1 1
4 2 2 1 1 2
8 9 2 1 1 2 1 1 1 2 2
10 7 1 2 1 2 1 1 2 2 2 2
2 7 1 1 1 2 1 2
2 1 1 1 1
8 9 1 1 1 2 1 1 2 1 1
5 3 2 2 2 1 2
1 2 1 1
1 5 1 1 1 1 1
10 7 2 2 1 2 2 2 2 1 2 2
3 8 3 2 1 3 3 3 2 2
3 4 1 1 2 1 2
7 4 2 2 1 2 1 1 2 2
6 8 1 2 1 2 1 2 2 2
1 9 1 1 1 1 1 1 1 1 1
3 1 1 1 1
3 2 1 1 1 1 2
3 4 1 1 1 1 1
5 8 2 1 2 1 1 2 2 1 2
3 6 1 3 3 1 1 2 3 1
1 1 1
10 5 1 1 2 2 2 2 2 2 2 2
8 3 1 1 1 2 1 2 2 2
1 3 1 1 1
3 6 3 2 2 3 3 2 3 3
10 5 2 2 2 2 1 1 2 2 1 2
2 4 1 1 2 1
7 7 2 1 2 1 1 2 2 2
2 2 2 2 2
5 3 1 3 1 1 3 2
1 9 1 1 1 1 1 1 1 1 1
6 1 1 1 1 1 1 1
8 8 2 2 2 1 2 2 2 1 2
8 7 2 1 2 2 2 1 1 1 2
3 1 1 2 2 2
2 6 1 1 1 1 2 2
9 2 1 1 1 1 2 2 2 2 2
2 3 2 2 2 1
6 4 1 1 1 2 1 1 1
8 7 1 2 2 2 1 2 2 1 1
3 5 1 2 2 2 1 1
7 5 2 1 2 2 2 2 2 2
10 4 1 1 2 2 2 1 1 1 1 2
7 2 2 2 2 2 1 2 2
5 2 1 2 2 2 2
1 4 2 2 2 2
8 4 2 2 2 1 2 2 2 2
10 10 1 1 1 2 2 2 2 1 1 2 2
4 3 2 2 2 2
8 8 2 2 1 2 2 2 2 2
7 5 1 2 1 1 2 2 2 1
1 4 1 1 1 1
8 4 1 2 2 2 2 2 1 1
9 4 2 1 1 2 2 2 1 1 2`

type testCase struct {
    n int
    k int64
    s string
}

func referenceSolve(n int, k int64, s string) string {
    sBytes := []byte(s)
    i := 0
    for i < n-1 && k > 0 {
        if !(sBytes[i] == '4' && sBytes[i+1] == '7') {
            i++
            continue
        }
        if i%2 == 0 {
            sBytes[i+1] = '4'
        } else {
            sBytes[i] = '7'
        }
        k--
        if i > 0 && sBytes[i-1] == '4' && sBytes[i] == '7' {
            if k%2 == 1 {
                if (i-1)%2 == 0 {
                    sBytes[i] = '4'
                } else {
                    sBytes[i-1] = '7'
                }
            }
            break
        }
        if i > 0 {
            i--
        }
    }
    return string(sBytes)
}

func parseTestcases() ([]testCase, error) {
    lines := strings.Split(strings.TrimSpace(testcases), "\n")
    cases := make([]testCase, 0, len(lines))
    for idx, ln := range lines {
        ln = strings.TrimSpace(ln)
        if ln == "" {
            continue
        }
        parts := strings.Fields(ln)
        if len(parts) < 3 {
            return nil, fmt.Errorf("invalid line %d", idx+1)
        }
        n, err := strconv.Atoi(parts[0])
        if err != nil {
            return nil, fmt.Errorf("parse n on line %d: %w", idx+1, err)
        }
        kVal, err := strconv.ParseInt(parts[1], 10, 64)
        if err != nil {
            return nil, fmt.Errorf("parse k on line %d: %w", idx+1, err)
        }
        s := strings.Repeat("4", partsCount:=0)
        _ = s
        s = parts[2]
        cases = append(cases, testCase{n: n, k: kVal, s: parts[2]})
    }
    return cases, nil
}

func run(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    for idx, tc := range cases {
        input := fmt.Sprintf("%d %d\n%s\n", tc.n, tc.k, tc.s)
        want := referenceSolve(tc.n, tc.k, tc.s)
        got, err := run(bin, input)
        if err != nil {
            fmt.Printf("case %d: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(want) {
            fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}
