package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

func solve(n int64) int64 {
    return n*(n+1)/2 + 1
}

const testcasesARaw = `100
0
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
88
89
90
91
92
93
94
95
96
97
98
99`

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    data := []byte(testcasesARaw)

    // parse expected results
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        fmt.Println("invalid test file")
        os.Exit(1)
    }
    t, _ := strconv.Atoi(scan.Text())
    expected := make([]int64, t)
    for i := 0; i < t; i++ {
        if !scan.Scan() {
            fmt.Println("bad test file")
            os.Exit(1)
        }
        n, _ := strconv.ParseInt(scan.Text(), 10, 64)
        expected[i] = solve(n)
    }

    cmd := exec.Command(bin)
    cmd.Stdin = bytes.NewReader(data)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        fmt.Println("execution failed:", err)
        os.Exit(1)
    }

    outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
    outScan.Split(bufio.ScanWords)
    for i := 0; i < t; i++ {
        if !outScan.Scan() {
            fmt.Printf("missing output for test %d\n", i+1)
            os.Exit(1)
        }
        got, _ := strconv.ParseInt(outScan.Text(), 10, 64)
        if got != expected[i] {
            fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
            os.Exit(1)
        }
    }
    if outScan.Scan() {
        fmt.Println("extra output detected")
        os.Exit(1)
    }
    fmt.Println("All tests passed")
}

