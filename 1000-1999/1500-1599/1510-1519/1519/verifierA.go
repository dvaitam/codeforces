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
const testcasesRaw = `25 49 3
3 17 4
32 26 2
31 23 4
14 33 1
19 9 0
40 17 4
46 39 1
20 7 5
5 44 2
31 36 0
23 28 2
40 41 1
36 31 3
34 17 0
36 1 0
47 26 5
43 41 0
40 32 2
16 47 2
46 5 1
37 15 1
10 35 3
6 6 2
33 32 0
20 36 2
46 8 4
22 35 1
39 36 4
19 29 0
39 25 2
37 16 2
12 13 1
3 40 5
17 31 0
6 44 1
10 3 0
45 35 5
26 46 4
18 34 1
14 44 4
27 38 2
29 32 5
42 45 2
6 21 4
8 32 4
41 22 1
16 2 5
18 8 5
15 24 1
22 28 0
7 10 5
15 3 4
41 35 4
44 5 0
8 41 1
39 37 0
26 6 2
8 3 4
2 13 1
46 8 3
14 47 0
44 2 4
28 40 0
17 5 1
5 42 2
23 28 1
4 33 3
3 39 0
45 26 1
17 23 5
31 37 1
45 44 1
50 4 5
11 11 2
34 17 0
39 29 5
12 1 3
44 27 4
33 20 5
23 25 5
17 10 4
45 1 3
48 6 2
48 3 4
18 9 1
49 31 2
40 19 5
23 38 5
40 9 5
20 25 5
27 42 0
1 39 1
45 22 1
16 15 5
29 25 5
44 37 3
3 26 5
37 27 5
46 3 1`

type testCase struct {
	r int64
	b int64
	d int64
}

func solveCase(tc testCase) string {
	r, b, d := tc.r, tc.b, tc.d
	if r < b {
		r, b = b, r
	}
	if r <= b*(d+1) {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("case %d: malformed", idx+1)
		}
		r, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse r: %v", idx+1, err)
		}
		b, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse b: %v", idx+1, err)
		}
		d, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse d: %v", idx+1, err)
		}
		res = append(res, testCase{r: r, b: b, d: d})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
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

	for i, tc := range cases {
		expected := solveCase(tc)
		input := fmt.Sprintf("1\n%d %d %d\n", tc.r, tc.b, tc.d)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
