package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"
)

func solve(reader io.Reader) (string, error) {
    in := bufio.NewReader(reader)
    var n, m int
    if _, err := fmt.Fscan(in, &n, &m); err != nil {
        return "", err
    }
    a := make([][]int, n)
    for i := 0; i < n; i++ {
        a[i] = make([]int, m)
        for j := 0; j < m; j++ {
            if _, err := fmt.Fscan(in, &a[i][j]); err != nil {
                return "", err
            }
        }
    }
    r := make([]int, n)
    c := make([]int, m)
    rowSum := make([]int, n)
    colSum := make([]int, m)
    for i := 0; i < n; i++ {
        sum := 0
        for j := 0; j < m; j++ {
            sum += a[i][j]
        }
        rowSum[i] = sum
    }
    for j := 0; j < m; j++ {
        sum := 0
        for i := 0; i < n; i++ {
            sum += a[i][j]
        }
        colSum[j] = sum
    }
    changed := true
    for changed {
        changed = false
        // flip negative-sum rows
        for i := 0; i < n; i++ {
            if rowSum[i] < 0 {
                changed = true
                r[i] ^= 1
                rowSum[i] = -rowSum[i]
                for j := 0; j < m; j++ {
                    oldParity := (r[i]^1 + c[j]) & 1
                    var oldVal int
                    if oldParity == 1 {
                        oldVal = -a[i][j]
                    } else {
                        oldVal = a[i][j]
                    }
                    colSum[j] += -2 * oldVal
                }
            }
        }
        // flip negative-sum columns
        for j := 0; j < m; j++ {
            if colSum[j] < 0 {
                changed = true
                c[j] ^= 1
                colSum[j] = -colSum[j]
                for i := 0; i < n; i++ {
                    oldParity := (r[i] + (c[j]^1)) & 1
                    var oldVal int
                    if oldParity == 1 {
                        oldVal = -a[i][j]
                    } else {
                        oldVal = a[i][j]
                    }
                    rowSum[i] += -2 * oldVal
                }
            }
        }
    }
    // collect result
    var rows []int
    for i := 0; i < n; i++ {
        if r[i] == 1 {
            rows = append(rows, i+1)
        }
    }
    var cols []int
    for j := 0; j < m; j++ {
        if c[j] == 1 {
            cols = append(cols, j+1)
        }
    }
    var sb strings.Builder
    writer := bufio.NewWriter(&sb)
    // output rows
    fmt.Fprint(writer, len(rows))
    if len(rows) > 0 {
        fmt.Fprint(writer, " ")
        for i, v := range rows {
            if i > 0 {
                fmt.Fprint(writer, " ")
            }
            fmt.Fprint(writer, v)
        }
    }
    fmt.Fprintln(writer)
    // output columns
    fmt.Fprint(writer, len(cols))
    if len(cols) > 0 {
        fmt.Fprint(writer, " ")
        for i, v := range cols {
            if i > 0 {
                fmt.Fprint(writer, " ")
            }
            fmt.Fprint(writer, v)
        }
    }
    fmt.Fprintln(writer)

    if err := writer.Flush(); err != nil {
        return "", err
    }
    return sb.String(), nil
}

const testcasesRaw = `4 4
-5 -1 3 2
1 -1 2 0
4 -2 3 -3
-1 -3 -4 4

3 5
4 -3 -1 -4 -4
5 0 2 3 -4
0 1 0 4 5

2 5
2 2 3 -1 -5
3 -5 -4 1 5

1 5
2 0 -2 0 -4

2 5
-2 -2 -3 3 2
-4 -4 0 3 2

1 3
3 -1 -4

5 3
3 -2 4
3 4 -1
2 -4 4
1 0 4
-2 -1 -3

2 2
-5 4
5 -1

4 1
-4
5
-3
-3

1 1
3

4 5
-1 3 -2 -2 5
4 1 4 -1 2
2 5 5 0 -4
0 4 -4 2 4

3 2
-2 -5
-1 -4
-2 0

2 3
1 -5 -4
-3 -2 -5

5 5
4 5 -4 -5 -4
5 -2 4 4 -4
1 -4 0 -4 -5
4 -5 -2 -3 -4
2 -2 -5 5 -5

5 4
4 -4 -1 -4
-2 -4 5 -1
0 1 -3 -5
3 2 -5 4
-4 1 -2 -1

3 4
4 -3 5 -2
-5 5 -3 -3
0 3 -1 -4

5 4
5 -3 -5 2
5 1 4 3
-1 5 0 1
5 -1 -3 3
-5 2 -4 0

1 5
-1 -3 -2 2 0

5 3
5 0 4
5 4 -3
-1 1 1
5 -4 -5
4 -2 0

2 2
-2 5
2 1

5 4
-5 1 4 1
5 -5 -3 2
-4 -1 -3 2
3 2 3 4
-5 -5 2 0

3 4
-5 1 -2 3
5 -4 -3 -5
1 5 1 0

1 2
-5 -5

5 5
-4 -2 -4 4 5
-2 -1 -1 -3 -4
2 1 5 -4 -5
-1 2 -4 -1 -3
5 3 5 5 0

1 2
-1 -5

1 1
-2

3 5
0 0 4 -5 4
5 2 5 2 5
1 0 3 -3 -2

4 5
-1 -5 -3 -3 -1
0 0 0 -4 0
4 -5 -5 -1 -3
-3 4 -1 0 1

5 2
-1 -4
2 -2
-5 -1
-3 3
-4 -1

4 3
-1 1 -4
-4 3 2
2 0 0
-4 2 -4

4 4
-5 -1 0 5
-3 -3 5 4
1 5 -4 -4
-4 -2 -2 -5

4 1
-4
1
3
3

3 4
2 4 5 -2
1 -4 0 -2
-1 4 -3 1

2 3
-4 -4 -5
3 2 5

2 1
2
1

3 2
5 -5
-2 4
-3 -4

2 4
1 0 3 -3
-4 4 2 -3

5 4
5 5 1 3
2 5 0 2
2 5 5 -2
3 4 -2 -5
0 0 0 -5

5 2
-1 4
-3 1
4 -1
2 -4
-4 3

1 1
-2

2 1
-1
-5

4 3
-3 -3 5
2 0 3
1 3 3
-5 4 -4

5 5
-4 1 -2 -1 3
4 1 2 1 4
4 -2 -5 5 -5
-3 -1 3 4 -1
0 -4 2 -1 -1

4 4
1 -5 -3 5
-3 -2 -1 0
-5 -5 2 1
-3 2 4 -4

2 3
1 -5 4
2 1 2

1 1
2

2 1
-5
4

5 2
5 0
-4 3
5 0
-2 1
2 -4

1 5
2 4 5 0 5

1 5
-1 -3 1 -1 5

1 5
-2 -5 1 2 0

2 4
0 5 -4 -5
-5 2 -1 -5

5 5
4 -2 -2 -4 5
3 3 1 3 -1
-4 -3 1 4 1
-4 -4 1 -4 -4
1 -3 -5 2 1

4 1
2
0
-1
-4

3 1
-4
0
-5

3 3
-3 -5 -2
0 -4 4
-3 -2 -5

2 1
-5
-1

3 1
4
-2
-3

2 4
-4 2 0 -1
-3 -5 -2 0

3 4
-1 -1 3 5
0 -3 4 -4
-4 3 4 -1

2 4
-3 -3 -2 0
3 -2 -2 -3

3 3
1 5 -5
-3 4 -5
1 -4 -4

2 4
-1 3 1 -3
4 1 -1 5

3 1
-2
2
5

3 5
-5 1 1 -5 1
0 2 -2 0 -1
2 -4 -3 -4 -1

1 5
4 -3 2 1 -3

4 4
-3 -2 2 0
3 -3 0 2
5 5 -4 2
-2 -1 -5 2

5 4
-5 -2 -1 -4
5 -1 3 4
-3 1 2 -4
5 2 -2 3
1 -1 5 -5

1 3
5 -5 -5

3 4
3 4 1 2
-4 -1 0 -1
5 -2 4 -4

1 1
-1

3 5
0 -4 3 -2 -3
-4 1 -1 -1 3
-3 4 3 5 -2

5 1
1
5
3
1
-1

3 4
0 4 5 -3
-3 -4 -4 1
1 4 2 -3

5 3
0 5 2
1 -2 2
2 3 0
2 5 -5
2 -1 -3

4 1
4
-2
-5
0

4 4
-5 3 -4 5
-4 5 5 1
-5 0 -5 -4
4 -5 -1 5

3 2
-3 4
-1 -2
-4 1

4 3
1 -3 0
1 5 5
1 -3 2
-3 3 0

2 2
-3 2
0 1

4 4
1 -2 -2 2
-2 4 -5 1
-5 -2 5 -4
-3 0 -5 5

2 2
4 -1
4 -4

5 3
0 1 2
-5 5 3
5 5 3
1 4 2
2 -1 2

2 3
-1 -5 -5
-5 -3 0

1 3
5 -5 -3

1 4
5 -2 4 1

5 2
2 -2
0 4
-4 4
-4 0
0 3

4 3
-1 -5 3
-5 -2 0
-4 -2 3
0 -2 -2

3 3
-1 3 1
-1 2 0
-2 -5 -1

5 1
-5
2
2
2
-5

4 4
2 2 -4 -4
-4 -2 -4 -3
1 -2 2 4
-4 1 3 1

1 2
-2 2

2 2
-1 0
0 1

1 5
-1 4 3 -2 -1

4 5
4 2 3 5 -1
-1 -2 -5 -4 4
-4 -3 1 -2 -2
-1 5 -5 3 3

4 1
-4
1
5
-1

1 5
0 -2 5 3 5

3 2
-2 -4
3 -1
5 0

2 3
5 2 -1
4 -3 -3

1 5
3 0 0 4 5`

func parseTestcases(raw string) []string {
    blocks := strings.Split(strings.TrimSpace(raw), "\n\n")
    res := make([]string, 0, len(blocks))
    for _, block := range blocks {
        block = strings.TrimSpace(block)
        if block == "" {
            continue
        }
        res = append(res, block)
    }
    return res
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    testcases := parseTestcases(testcasesRaw)
    for idx, tc := range testcases {
        expected, err := solve(strings.NewReader(tc))
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
            os.Exit(1)
        }

        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(tc + "\n")
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }

        got := strings.TrimSpace(string(out))
        if got != strings.TrimSpace(expected) {
            fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, strings.TrimSpace(expected), got)
            os.Exit(1)
        }
    }

    fmt.Printf("All %d tests passed\n", len(testcases))
}
