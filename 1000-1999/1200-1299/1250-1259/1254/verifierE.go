package main

import (
    "bytes"
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

// Embedded testcases (each case: n, n-1 edges, then n values)
const testcases = `6
1 2
2 3
3 4
1 5
4 6
2 0 1 6 0 0
4
1 2
2 3
1 4
0 0 1 0
3
1 2
1 3
2 3 1
3
1 2
1 3
0 2 3
5
1 2
2 3
3 4
1 5
5 4 1 0 4
5
1 2
1 3
1 4
4 5
0 0 0 0 1
6
1 2
2 3
3 4
1 5
2 6
3 3 1 1 2 4
6
1 2
2 3
3 4
3 5
3 6
3 2 2 2 2 2
7
1 2
2 3
1 4
4 5
4 6
4 7
3 0 3 0 3 0 3
2
1 2
1 1
3
1 2
2 3
1 1 1
2
1 2
0 1
6
1 2
2 3
3 4
1 5
5 6
1 0 2 5 1 5
4
1 2
2 3
3 4
2 3 2 3
7
1 2
2 3
1 4
2 5
5 6
6 7
0 3 0 3 0 3 0
5
1 2
2 3
3 4
4 5
1 1 1 2 2
2
1 2
0 0
4
1 2
2 3
3 4
3 1 2 1
2
1 2
0 0
5
1 2
2 3
3 4
4 5
0 0 0 0 0
6
1 2
1 3
1 4
1 5
1 6
1 2 3 3 3 3
6
1 2
1 3
1 4
1 5
1 6
0 3 3 3 3 3
6
1 2
1 3
1 4
1 5
1 6
2 3 3 3 3 3
6
1 2
1 3
1 4
1 5
1 6
3 3 3 3 3 3
6
1 2
1 3
1 4
1 5
1 6
0 2 3 3 3 3
6
1 2
1 3
1 4
1 5
1 6
1 3 3 3 3 3
6
1 2
1 3
1 4
1 5
1 6
3 2 3 3 3 3
7
1 2
2 3
2 4
1 5
5 6
6 7
5 5 5 5 5 5 5
7
1 2
2 3
3 4
3 5
5 6
5 7
1 1 1 1 1 1 1
6
1 2
2 3
3 4
4 5
4 6
3 3 1 3 1 2
6
1 2
2 3
3 4
3 5
5 6
1 1 1 1 1 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 2 2 2 2 2 2 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 2 2 2 2 2 2 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 1 2 2 2 2 2 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 2 1 2 2 2 2 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 2 2 1 2 2 2 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 2 2 2 1 2 2 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 2 2 2 2 1 2 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 2 2 2 2 2 1 2
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
2 2 2 2 2 2 2 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 1 0 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 1 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 0 0 0 1 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 0 0 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 0 0 0 0 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 1 1 1 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 1 1 1 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 1 1 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 1 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 0 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 0 0 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 1 1 1 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 1 1 0 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 1 0 1 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 1 0 1 1 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 1 0 1 1 1 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 1 0 1 1 1 1 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 0 1 1 1 1 1 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 1 1 1 1 1 1 1
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
1 0 0 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 1 0 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 1 0 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 1 0 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 0 1 0 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 0 0 1 0 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 0 0 0 1 0
8
1 2
2 3
3 4
4 5
5 6
6 7
7 8
0 0 0 0 0 0 0 1
9
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
1 2 1 2 1 2 1 2 1
9
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
2 1 2 1 2 1 2 1 2
9
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
2 2 1 2 1 2 1 2 1
9
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
1 2 1 2 1 2 1 2 0
9
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
2 2 1 2 1 2 1 1 1
9
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
1 1 1 2 1 1 1 1 1
10
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10
1 2 1 2 1 2 1 2 1 2
10
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10
2 1 2 1 2 1 2 1 2 1
10
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10
2 2 1 2 1 2 1 2 1 2
10
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10
1 2 1 2 1 2 1 2 1 0
10
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10
2 2 1 2 1 2 1 1 1 1
10
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10
1 1 1 2 1 1 1 1 1 1
2
1 2
0 0
1
1
0
1
1
1
1
1
1
1
1
1
1
1
`

// referenceSolve embeds the simple solver from 1254A.go
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
        s := parts[2]
        cases = append(cases, testCase{n: n, k: kVal, s: s})
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
