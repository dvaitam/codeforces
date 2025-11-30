package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

const testcaseData = `100
6 20
51 84 7 10 69 13
6 75
8 65 28 5 12 56
7 9
31 12 71 55 8 73 16
4 81
81 75 8 74
10 51
7 29 6 72 18 38 54 19 70 16
10 40
72 88 24 14 75 74 82 25 48 13
9 92
9 73 8 80 27 64 88 69 55
6 60
75 59 47 39 32 24
4 11
74 39 68 64
6 94
58 37 78 10 16 66
7 22
97 44 20 63 54 6 86
2 98
72 74
6 44
89 45 77 64 75 59
2 12
35 61
2 8
94 90
5 83
74 88 58 37 92
7 86
45 3 60 46 22 79 15
8 8
28 99 37 17 95 32 51 51
8 11
22 58 52 71 36 18 56 71
5 91
54 46 88 49 30
3 11
23 20 30
4 2
63 76 24 34
5 1
19 54 69 48 79
10 41
17 89 66 80 84 87 95 7 59 100
9 51
51 52 51 14 62 82 52 8 25
2 27
57 21
2 44
77 7
2 1
73 20
9 13
47 79 4 10 27 79 49 20 82
5 45
78 47 61 16 15
8 60
62 62 40 11 19 14 96 44
5 62
89 21 67 3 27
9 47
19 89 70 4 98 68 39 83 12
5 67
47 22 46 99 29
9 70
100 65 43 82 29 79 98 25 31
7 95
30 26 67 64 46 94 4
1 36
61
5 25
89 78 45 58 93
6 47
11 29 14 30 61 26
6 27
62 80 79 1 62 84
6 83
11 85 16 50 92 97
4 62
23 56 82 43
2 93
51 60
7 96
11 93 21 22 17 4 20
10 60
84 19 79 77 61 85 45 20 71 71
3 3
2 93 84
2 68
96 18
7 25
28 4 33 28 38 65 31
10 42
34 70 54 17 8 95 46 59 85 75
9 54
65 17 69 20 68 66 3 57 100
3 78
1 100 20
3 19
61 80 93
2 72
8 42
9 68
72 62 100 14 72 8 32 25 36
1 99
13
9 58
72 4 98 9 57 42 79 65 78
9 26
89 36 58 66 69 62 65 32 90
9 34
72 26 58 18 54 16 51 57 41
2 86
31 55
2 28
86 39
2 100
20 92
6 19
33 18 60 29 96 13
7 63
21 86 29 21 91 56 66
7 44
54 26 46 41 12 93 47
1 44
71
8 57
91 3 50 43 67 80 38 66
2 15
30 14
2 34
35 6
3 35
97 17 55
5 52
20 69 66 74 64
6 12
36 8 89 24 55 10
5 3
82 12 34 11 78
4 9
34 16 59 2
6 71
54 35 80 17 6 68
4 15
21 34 7 24
4 40
81 40 68 98
4 38
58 65 87 23
5 45
3 33 5 2 3
9 71
25 66 61 32 58 14 85 84 56
8 70
51 65 40 89 28 30 44 26
3 52
45 7 17
1 10
81
5 56
21 8 11 86 49
9 86
37 77 32 89 38 6 59 24 21
5 58
1 34 47 43 71
6 32
5 40 28 46 24 1
6 49
11 61 36 65 84 26
4 65
100 1 12 34
2 19
52 76
1 51
3
5 39
81 30 11 75 68
3 85
92 77 50
6 93
64 20 37 93 80 83
3 6
92 66 81
7 94
90 65 18 68 97 65 73
1 88
75
4 11
4 6 18 82
6 14
49 58 72 7 81 3
9 88
32 63 34 1 59 9 96 65 69
2 85
68 9
`

func solveCase(n int, c int64, a []int64) int {
    type pair struct {
        cost int64
        idx  int
    }
    pairs := make([]pair, n)
    for i := 1; i <= n; i++ {
        cost := int64(i)
        if tmp := int64(n+1-i); tmp < cost {
            cost = tmp
        }
        cost += a[i]
        pairs[i-1] = pair{cost: cost, idx: i}
    }
    sort.Slice(pairs, func(i, j int) bool { return pairs[i].cost < pairs[j].cost })

    pref := make([]int64, n)
    pos := make([]int, n+1)
    for i, p := range pairs {
        if i == 0 {
            pref[i] = p.cost
        } else {
            pref[i] = pref[i-1] + p.cost
        }
        pos[p.idx] = i
    }

    countWithBudget := func(b int64) int {
        if b < 0 {
            return 0
        }
        return sort.Search(n, func(i int) bool { return pref[i] > b })
    }

    ans := countWithBudget(c)
    for i := 1; i <= n; i++ {
        firstCost := int64(i) + a[i]
        if firstCost > c {
            continue
        }
        rem := c - firstCost
        k := countWithBudget(rem)
        if pos[i] < k {
            k--
        }
        if 1+k > ans {
            ans = 1 + k
        }
    }
    return ans
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
        a := make([]int64, n+1)
        for i := 1; i <= n; i++ {
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
        fmt.Println("usage: verifierG2 /path/to/binary")
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
