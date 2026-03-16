package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func expect(x int) string {
    if x%3==0 || x%5==0 {
        return "YES"
    }
    return "NO"
}

const testcasesERaw = `100
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
99
100`

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    data := []byte(testcasesERaw)
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() { fmt.Println("invalid test file"); os.Exit(1) }
    var t int
    fmt.Sscan(scan.Text(), &t)
    expected := make([]string, t)
    for i:=0;i<t;i++{
        if !scan.Scan(){fmt.Println("bad test file"); os.Exit(1)}
        var x int
        fmt.Sscan(scan.Text(), &x)
        expected[i]=expect(x)
    }
    cmd := exec.Command(bin)
    cmd.Stdin = bytes.NewReader(data)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        fmt.Println("execution failed:", err)
        os.Exit(1)
    }
    outs := strings.Fields(strings.TrimSpace(out.String()))
    if len(outs)!=t { fmt.Printf("expected %d lines got %d\n", t,len(outs)); os.Exit(1) }
    for i:=0;i<t;i++{
        if outs[i]!=expected[i]{
            fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], outs[i])
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
