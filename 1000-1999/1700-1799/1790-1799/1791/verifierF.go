package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "strconv"
)

type query struct {
    typ int
    l   int
    r   int
}

const testcaseData = `100
5 3
12 17 1 15 8
1 2 2
2 4
1 4 4
5 2
1 7 14 9 6
2 2
1 2 5
2 2
1 1
1 1 1
1 2 2
2 5
7 6
1 2 2
1 2 2
1 1 2
1 2 2
1 2 2
3 3
10 16 11
1 2 3
1 1 2
1 3 3
4 1
18 14 12 13
1 4 4
2 5
7 4
1 2 2
2 2
2 1
2 2
1 2 2
2 3
17 20
2 1
2 2
1 2 2
3 2
3 5 10
2 1
1 1 3
5 4
2 8 20 12 9
2 4
1 1 1
2 3
1 2 3
4 1
6 14 12 5
1 4 4
2 4
20 6
2 2
2 2
2 2
2 2
2 1
13 18
1 2 2
2 1
16 9
2 1
3 5
2 10 12
2 2
2 3
2 2
1 3 3
2 3
3 3
9 15 10
2 2
2 3
2 3
4 3
6 15 12 11
1 2 2
2 4
2 1
4 2
20 19 17 14
2 3
1 2 2
5 4
20 6 8 6 2
2 2
1 1 2
1 3 3
2 2
5 1
14 15 12 13 20
1 5 5
2 3
1 12
2 2
2 1
2 1
5 5
10 4 10 18 17
2 5
2 3
1 4 5
2 4
1 2 5
5 4
7 5 20 3 12
1 4 4
2 5
1 3 5
2 4
4 2
16 10 16 13
2 2
2 3
4 3
14 1 11 10
2 3
1 4 4
1 4 4
3 1
5 13 1
2 3
5 3
8 16 2 8 16
2 2
2 3
2 5
4 5
20 4 1 5
2 3
2 3
1 4 4
2 2
1 1 3
5 4
4 18 7 1 14
2 5
2 4
2 4
1 3 4
1 3
1
2 1
2 1
1 1 1
5 4
11 18 5 14 19
1 1 2
2 1
1 1 3
2 1
4 4
2 4 4 8
1 2 3
2 2
1 2 3
1 1 2
1 5
4
2 1
2 1
1 1 1
1 1 1
2 1
4 3
12 2 5 10
1 3 3
2 1
1 3 4
5 2
8 12 15 13 6
1 1 3
1 5 5
1 1
17
2 1
4 4
16 11 4 17
1 4 4
1 4 4
1 1 4
1 2 4
5 4
12 8 10 11 20
2 4
2 2
2 3
2 3
5 1
19 19 7 6 13
1 1 1
1 2
7
1 1 1
2 1
1 4
16
2 1
2 1
2 1
1 1 1
1 3
13
2 1
2 1
2 1
3 2
2 16 2
2 2
2 1
1 4
16
1 1 1
2 1
1 1 1
1 1 1
1 5
13
2 1
1 1 1
2 1
2 1
1 1 1
1 5
14
2 1
2 1
2 1
1 1 1
2 1
1 2
16
1 1 1
1 1 1
4 1
2 9 14 15
2 1
1 2
6
2 1
1 1 1
1 1
18
1 1 1
5 1
15 18 12 16 14
2 5
2 1
10 15
1 1 1
3 4
3 9 13
2 1
1 3 3
2 1
1 2 2
3 1
16 8 2
1 2 3
1 1
14
2 1
5 1
10 5 19 1 15
1 1 1
4 3
4 13 18 5
2 1
2 2
1 3 4
3 3
2 3 18
1 1 1
1 1 2
2 2
5 1
8 17 14 5 12
1 3 4
1 2
19
1 1 1
1 1 1
4 4
4 11 11 15
1 4 4
1 4 4
1 4 4
1 1 2
4 1
7 1 2 19
1 4 4
4 4
7 7 9 7
2 3
1 3 4
1 1 4
2 3
2 2
8 16
2 1
1 1 2
1 3
7
2 1
2 1
1 1 1
1 5
20
1 1 1
2 1
2 1
2 1
2 1
5 2
14 12 19 1 17
1 3 4
1 1 1
3 5
19 14 10
1 2 2
1 1 1
2 2
1 2 2
1 1 2
5 4
13 1 13 13 6
1 5 5
2 3
1 5 5
1 1 2
3 3
9 20 3
1 1 3
2 1
1 3 3
5 5
18 14 14 8 14
2 5
1 3 4
2 1
2 2
2 1
2 1
10 19
1 2 2
3 5
4 11 3
1 2 2
1 2 2
2 3
2 1
2 1
1 2
6
1 1 1
1 1 1
5 4
11 6 18 19 4
2 1
1 2 3
2 3
1 4 4
2 4
18 15
2 1
1 2 2
1 2 2
2 1
4 3
1 4 16 12
1 3 4
1 3 4
1 1 3
4 3
3 17 14 2
1 4 4
2 1
1 3 3
5 4
14 6 18 13 8
1 4 5
1 1 1
1 3 3
1 3 5
4 4
17 3 2 7
1 4 4
2 1
2 4
1 1 2
4 5
14 19 20 19
1 2 2
2 4
2 4
1 2 3
1 3 3
4 1
12 7 13 14
1 3 3
4 3
10 1 4 6
2 1
1 4 4
2 2
2 1
16 11
2 1
4 2
3 12 20 1
1 2 3
1 3 3
3 3
7 5 10
1 2 2
1 3 3
2 2
3 4
18 19 15
1 3 3
2 1
1 3 3
1 1 2
1 1
1
1 1 1
5 5
8 10 2 17 8
1 4 5
1 5 5
2 4
1 1 3
2 5
5 4
13 10 6 19 3
2 1
1 5 5
1 5 5
2 4
2 4
12 3
2 1
1 2 2
1 1 1
2 1
3 3
9 9 16
1 3 3
2 2
1 3 3
4 2
8 9 19 15
1 2 4
1 3 3
1 4
2
2 1
1 1 1
2 1
1 1 1
5 1
2 3 15 10 17
1 3 3
4 3
9 15 10 10
1 3 4
2 4
1 2 3
3 5
7 15 3
1 2 2
2 2
2 1
1 1 2
1 3 3
5 3
1 5 14 9 19
2 5
1 5 5
1 4 5
3 3
20 4 7
2 2
2 3
2 2
1 1
4
2 1
`

func sumDigits(x int) int {
    res := 0
    for x > 0 {
        res += x % 10
        x /= 10
    }
    return res
}

func computeExpected() (string, error) {
    fields := strings.Fields(testcaseData)
    if len(fields) == 0 {
        return "", fmt.Errorf("no testcases")
    }
    pos := 0
    t, err := strconv.Atoi(fields[pos])
    if err != nil {
        return "", fmt.Errorf("bad test count: %w", err)
    }
    pos++
    var out strings.Builder
    for caseNum := 0; caseNum < t; caseNum++ {
        if pos+1 >= len(fields) {
            return "", fmt.Errorf("case %d: missing n/q", caseNum+1)
        }
        n, err := strconv.Atoi(fields[pos])
        if err != nil {
            return "", fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
        }
        pos++
        q, err := strconv.Atoi(fields[pos])
        if err != nil {
            return "", fmt.Errorf("case %d: bad q: %w", caseNum+1, err)
        }
        pos++
        a := make([]int, n+2)
        for i := 1; i <= n; i++ {
            if pos >= len(fields) {
                return "", fmt.Errorf("case %d: missing array value", caseNum+1)
            }
            v, err := strconv.Atoi(fields[pos])
            if err != nil {
                return "", fmt.Errorf("case %d: bad array value: %w", caseNum+1, err)
            }
            a[i] = v
            pos++
        }
        parent := make([]int, n+2)
        for i := 0; i <= n+1; i++ {
            parent[i] = i
        }
        var find func(int) int
        find = func(x int) int {
            if parent[x] != x {
                parent[x] = find(parent[x])
            }
            return parent[x]
        }
        for qi := 0; qi < q; qi++ {
            if pos >= len(fields) {
                return "", fmt.Errorf("case %d: missing query type", caseNum+1)
            }
            typ, err := strconv.Atoi(fields[pos])
            if err != nil {
                return "", fmt.Errorf("case %d: bad query type: %w", caseNum+1, err)
            }
            pos++
            if pos >= len(fields) {
                return "", fmt.Errorf("case %d: missing l", caseNum+1)
            }
            l, err := strconv.Atoi(fields[pos])
            if err != nil {
                return "", fmt.Errorf("case %d: bad l: %w", caseNum+1, err)
            }
            pos++
            if typ == 1 {
                if pos >= len(fields) {
                    return "", fmt.Errorf("case %d: missing r", caseNum+1)
                }
                r, err := strconv.Atoi(fields[pos])
                if err != nil {
                    return "", fmt.Errorf("case %d: bad r: %w", caseNum+1, err)
                }
                pos++
                for p := find(l); p <= r; p = find(p+1) {
                    a[p] = sumDigits(a[p])
                    if a[p] < 10 {
                        parent[p] = find(p + 1)
                    }
                }
            } else {
                out.WriteString(strconv.Itoa(a[l]))
                out.WriteByte('\n')
            }
        }
    }
    return strings.TrimSpace(out.String()), nil
}

func run(bin string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(testcaseData)
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
        fmt.Println("usage: verifierF /path/to/binary")
        os.Exit(1)
    }
    expected, err := computeExpected()
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to compute expected outputs: %v\n", err)
        os.Exit(1)
    }
    got, err := run(os.Args[1])
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
        os.Exit(1)
    }
    if expected != got {
        fmt.Fprintf(os.Stderr, "output mismatch\nExpected:\n%s\nGot:\n%s\n", expected, got)
        os.Exit(1)
    }
    fmt.Println("All tests passed")
}
