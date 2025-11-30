package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `50 98 27 6
34 66 32 52
39 62 23 38
28 65 5 37
18 97 4 80
33 69 10 40
13 94 2 88
43 61 36 7
46 56 21 40
82 27 71 16
57 67 17 8
71 2 12 2
91 86 81 1
79 64 43 32
94 42 91 5
25 73 8 31
19 70 15 12
11 41 9 32
14 39 9 19
91 16 71 11
70 27 37 15
12 77 7 41
74 31 38 6
25 24 2 20
85 34 61 5
12 87 3 20
5 11 5 11
51 91 34 36
67 31 28 29
87 76 54 75
36 58 32 43
83 90 46 11
42 79 8 63
76 81 43 25
32 3 18 1
91 29 48 26
22 43 14 4
13 19 12 8
6 74 6 69
78 88 10 4
16 82 7 78
74 16 51 3
48 15 3 10
3 25 1 23
16 62 7 47
8 87 1 70
55 80 7 34
9 29 2 21
39 45 28 12
8 65 8 6
77 13 51 4
34 46 31 37
22 90 22 27
99 8 87 3
21 44 17 17
16 77 15 23
2 61 2 58
73 66 40 46
50 85 17 20
72 89 2 59
95 11 43 1
70 36 18 16
98 62 46 40
37 87 23 76
82 80 17 40
50 96 27 84
11 1 10 1
90 43 21 16
29 82 15 49
91 87 73 54
5 52 5 27
99 85 91 6
22 58 3 17
90 21 58 17
63 72 39 1
5 64 3 40
60 7 52 7
54 25 36 21
11 93 3 2
52 87 27 41
1 28 1 23
97 1 87 1
25 16 20 7
39 36 12 7
61 51 41 6
3 36 2 8
33 18 23 4
20 36 1 3
6 27 6 9
72 41 47 37
6 96 6 78
84 64 83 59
82 56 48 56
69 23 27 13
76 38 2 9
20 35 11 22
48 92 6 44
100 80 5 6
35 21 10 19
38 47 26 36`

func maxDistance(n, m, r, c int) int {
	d1 := (r - 1) + (c - 1)
	d2 := (r - 1) + (m - c)
	d3 := (n - r) + (c - 1)
	d4 := (n - r) + (m - c)
	ans := d1
	if d2 > ans {
		ans = d2
	}
	if d3 > ans {
		ans = d3
	}
	if d4 > ans {
		ans = d4
	}
	return ans
}

func parseTestcases() ([][4]int, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields)%4 != 0 {
		return nil, fmt.Errorf("testcase data malformed, got %d fields", len(fields))
	}
	cnt := len(fields) / 4
	res := make([][4]int, cnt)
	for i := 0; i < cnt; i++ {
		for j := 0; j < 4; j++ {
			v, err := strconv.Atoi(fields[i*4+j])
			if err != nil {
				return nil, fmt.Errorf("parse value %d: %v", i*4+j+1, err)
			}
			res[i][j] = v
		}
	}
	return res, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d %d %d %d\n", tc[0], tc[1], tc[2], tc[3])
	}

	var expected strings.Builder
	for i, tc := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(strconv.Itoa(maxDistance(tc[0], tc[1], tc[2], tc[3])))
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("verifier failed\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
