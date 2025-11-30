package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

const testcaseData = `100
10 11
63 98 34 5 1 19 85 76 61 98
6 41
99 3 35 63 26 94
7 69
70 88 13 25 73 71 90
5 85
79 88 12 55 43
2 47
53 33
8 90
13 97 26 90 82 38 13 6
10 26
84 47 63 25 66 74 83 90 65 4
6 32
78 56 39 46 76 16
2 65
87 68
4 15
78 85 35 40
4 49
62 29 18 77
4 90
67 2 25 100
3 3
83 43 72
10 80
40 48 49 68 50 38 17 87 63 7
3 55
77 96 51
2 57
32 12
10 86
58 58 49 10 67 55 61 39 91 53
2 25
96 88
5 58
63 93 23 3 3
9 16
33 76 47 25 33 65 58 43 67
5 53
54 79 63 35 78
8 85
62 63 19 93 49 64 40 82
8 42
47 85 21 80 49 89 78 35
6 83
51 63 94 21 38 72
1 80
58
1 24
4
10 100
73 15 88 92 48 47 64 76 8 25
3 35
79 2 55
9 63
10 61 30 13 48 47 19 88 81
4 78
41 19 5 81
2 92
14 6
8 60
93 9 81 4 90 71 17 79
2 21
80 81
5 84
58 64 89 3 82
3 25
57 62 100
8 58
72 36 62 92 75 9 38 47
5 47
5 11 65 36 36
5 61
56 57 48 5 92
2 32
78 82
4 69
96 5 95 23
6 48
5 7 95 98 97 27
7 44
10 21 94 87 15 52 71
3 60
96 64 82
2 17
10 45
9 4
24 68 91 20 24 37 95 20 65
10 50
29 29 90 73 92 17 54 89 21 97
6 31
72 74 1 73 50 25
6 49
68 57 2 51 68 71
3 92
30 83 70
6 35
86 24 20 22 89 81
3 50
10 8 2
1 61
98
1 15
84
9 13
87 48 77 5 8 70 29 53 29
10 89
38 12 61 10 93 74 8 22 85 57
2 15
90 23
1 88
4
7 84
91 37 33 57 30 26 84
9 83
28 36 76 2 33 68 94 96 91
6 100
67 12 75 77 83 5
6 47
89 15 29 98 4 47
6 56
88 7 45 18 4 2
6 9
63 1 60 34 12 87
10 79
35 30 17 17 54 70 24 32 85 50
10 55
73 67 81 62 35 33 18 39 5 25
8 14
23 74 14 1 66 21 16 20
10 51
61 63 68 99 89 70 61 28 83 50
3 99
59 2 44
7 33
41 18 37 48 76 87 47
4 1
100 40 23 11
7 80
80 99 30 83 12 25 49
7 17
97 19 67 7 67 68 100
8 47
43 92 76 50 35 84 85 9
5 17
6 15 76 78 79
3 65
65 15 89
2 12
86 97
2 28
45 44
3 94
32 23 24
4 81
17 62 95 59
10 15
93 79 18 59 26 74 48 94 19 43
5 9
36 97 1 8 84
1 62
81
8 41
8 2 36 93 13 50 80 44
2 95
48 19
3 48
86 98 61
7 3
69 20 14 58 48 2 17
9 62
34 10 62 16 52 91 26 70 22
8 43
33 70 100 45 38 70 89 35
8 72
92 99 6 65 96 32 85 80
1 68
64
9 69
33 22 5 42 2 27 74 39 25
3 36
8 35 94
5 90
48 28 91 94 22
10 1
12 46 74 24 45 24 62 82 29 61
7 71
80 46 16 84 94 12 56
1 75
80
4 78
57 35 29 82
1 62
85
9 68
33 67 35 82 47 39 55 96 48
2 98
70 54
`

func solveCase(n int, c int64, a []int64) int {
    costs := make([]int64, n)
    for i, v := range a {
        costs[i] = v + int64(i+1)
    }
    // sort costs ascending
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            if costs[j] < costs[i] {
                costs[i], costs[j] = costs[j], costs[i]
            }
        }
    }
    cnt := 0
    for _, v := range costs {
        if c >= v {
            c -= v
            cnt++
        } else {
            break
        }
    }
    return cnt
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
            return "", fmt.Errorf("case %d: missing n or c", caseNum+1)
        }
        n, err := strconv.Atoi(fields[pos])
        if err != nil {
            return "", fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
        }
        pos++
        cVal, err := strconv.ParseInt(fields[pos], 10, 64)
        if err != nil {
            return "", fmt.Errorf("case %d: bad c: %w", caseNum+1, err)
        }
        pos++
        a := make([]int64, n)
        for i := 0; i < n; i++ {
            if pos >= len(fields) {
                return "", fmt.Errorf("case %d: missing array value", caseNum+1)
            }
            v, err := strconv.ParseInt(fields[pos], 10, 64)
            if err != nil {
                return "", fmt.Errorf("case %d: bad array value: %w", caseNum+1, err)
            }
            a[i] = v
            pos++
        }
        res := solveCase(n, cVal, a)
        out.WriteString(strconv.Itoa(res))
        if caseNum+1 < t {
            out.WriteByte('\n')
        }
    }
    return out.String(), nil
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
        fmt.Println("usage: verifierG1 /path/to/binary")
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
