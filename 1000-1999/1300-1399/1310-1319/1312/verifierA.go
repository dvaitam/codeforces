package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded reference solution (1312A.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       if n%m == 0 {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
`

const testcasesRaw = `53 51
57 5
37 35
66 28
42 33
49 40
31 19
21 12
21 6
83 35
72 21
43 9
97 12
91 45
64 38
16 8
59 23
82 29
74 64
60 58
70 36
11 3
15 14
55 48
89 83
4 3
46 18
97 44
94 11
28 21
32 10
22 20
61 8
14 8
69 65
17 7
74 40
94 18
74 45
73 29
81 73
79 39
60 8
80 52
44 39
34 12
27 9
27 4
82 36
64 7
15 13
100 19
23 4
14 11
91 53
94 70
39 36
34 30
31 24
79 56
78 38
61 34
88 85
93 48
14 8
82 17
66 40
84 45
28 10
6 5
38 10
94 31
51 13
46 30
11 4
22 10
9 7
85 71
81 12
7 3
85 27
81 76
19 15
15 8
18 3
81 5
28 8
95 18
65 16
97 10
90 5
73 57
83 15
37 7
32 5
86 41
48 30
27 4
68 62
9 7`

var _ = solutionSource

type testCase struct {
	n int
	m int
}

func runProg(prog, input string) (string, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid test line %d", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n on line %d: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("invalid m on line %d: %v", idx+1, err)
		}
		tests = append(tests, testCase{n: n, m: m})
	}
	return tests, nil
}

func expected(tc testCase) string {
	if tc.n%tc.m == 0 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.m)
		want := expected(tc)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("candidate error on test %d: %v\n%s\n", idx+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
