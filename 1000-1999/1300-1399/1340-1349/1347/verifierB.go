package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
35 15
91 29
48 22
43 55
8 13
19 90
29 6
74 82
69 78
88 10
4 16
82 25
78 74
16 51
12 48
15 5
78 3
25 24
92 16
62 27
94 8
87 3
70 55
80 13
34 9
29 10
83 39
45 56
24 8
65 60
6 77
13 90
51 26
34 46
94 61
73 22
90 87
27 99
8 87
21 21
44 68
33 16
77 57
86 23
2 61
88 53
73 66
40 84
46 50
85 33
20 72
89 2
59 95
11 43
95 6
70 36
18 31
98 62
46 79
37 87
46 76
82 80
17 92
40 50
96 54
84 11
1 77
25 90
43 21
31 29
82 58
49 91
87 73
54 5
52 90
73 54
99 85
91 6
22 58
9 34
90 21
58 68
63 72
78 97
1 5
64 42
40 60
7 54
25 71
82 11
93 17
2 52
87 54
41 1
28 2
92 97
1 87
68 79
13 25
16 78
84 26
39 36
89 24
13 61
51 81
11 3
36 58
15 33
18 84
67 84
83 45
15 20
36 3
6 6
27 88
34 72
41 47
73 6
96 90
78 84
64 92
83 59
82 56
48 69
23 27
49 76
38 2
18 20
35 43
44 48
92 12
44 100
80 5
6 35
21 20
75 38
47 51
71 17
38 15
62 94
31 7
40 23
67 94
10 39
52 43
39 54
14 13
72 62
61 44
44 16
62 15
90 64
55 5
39 43
95 88
20 22
81 73
49 82
12 9
11 26
96 29
8 50
2 13
51 72
67 38
58 63
75 92
87 28
55 11
48 29
34 75
100 22
56 25
46 15
9 90
4 68
58 97
87 26
16 64
51 33
27 83
6 28
80 19
14 26
59 49
47 70
20 14
77 63
19 73
52 82
88 55
67 64
87 42
64 64
82 86
26 70
79 29
2 44
91 96
41 42
`

func solve(a1, b1, a2, b2 int) string {
	ok := false
	if a1 == a2 && b1+b2 == a1 {
		ok = true
	}
	if a1 == b2 && b1+a2 == a1 {
		ok = true
	}
	if b1 == a2 && a1+b2 == b1 {
		ok = true
	}
	if b1 == b2 && a1+a2 == b1 {
		ok = true
	}
	if ok {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	a1 int
	b1 int
	a2 int
	b2 int
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	if len(fields) != 1+4*t {
		return nil, fmt.Errorf("expected %d numbers got %d", 1+4*t, len(fields))
	}
	tests := make([]testCase, t)
	idx := 1
	for i := 0; i < t; i++ {
		a1, _ := strconv.Atoi(fields[idx])
		b1, _ := strconv.Atoi(fields[idx+1])
		a2, _ := strconv.Atoi(fields[idx+2])
		b2, _ := strconv.Atoi(fields[idx+3])
		tests[i] = testCase{a1: a1, b1: b1, a2: a2, b2: b2}
		idx += 4
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("1\n%d %d\n%d %d\n", tc.a1, tc.b1, tc.a2, tc.b2)
}

func runBinary(bin, input string) (string, error) {
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

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	expected := solve(tc.a1, tc.b1, tc.a2, tc.b2)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
