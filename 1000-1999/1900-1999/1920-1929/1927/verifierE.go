package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `8 5
4 4
16 5
3 1
1 1
18 10
26 25
2 1
17 12
9 3
27 4
9 4
30 1
27 21
26 9
26 9
7 2
10 5
21 12
3 3
11 11
13 9
8 3
8 8
9 2
30 27
28 18
27 10
1 1
19 10
28 25
17 7
14 7
20 10
14 8
6 2
10 5
27 26
2 1
2 2
21 9
17 16
23 11
5 2
3 2
30 7
21 21
15 5
6 3
14 12
19 11
21 18
7 3
4 1
23 8
9 4
4 3
30 6
10 8
1 1
12 12
3 2
24 22
11 1
11 5
11 3
25 21
14 14
28 20
22 3
10 10
7 4
10 3
9 7
20 6
11 10
1 1
2 2
6 3
26 26
12 5
19 4
15 4
14 4
4 1
2 1
24 6
20 5
20 2
18 16
19 8
11 1
4 3
25 14
21 7
16 7
8 8
14 8
2 1
14 8`

func runProg(bin, input string) (string, error) {
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

func generate(n, k int) []int {
	a := make([]int, n)
	start, end := 1, n
	for i := 0; i < k; i++ {
		for j := i; j < n; j += k {
			if i&1 == 1 {
				a[j] = end
				end--
			} else {
				a[j] = start
				start++
			}
		}
	}
	return a
}

type testCase struct {
	n int
	k int
}

func parseCases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid testcase at line %d", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n at line %d: %v", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("invalid k at line %d: %v", idx+1, err)
		}
		res = append(res, testCase{n: n, k: k})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases(testcases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		wantArr := generate(tc.n, tc.k)
		var expected strings.Builder
		for i, v := range wantArr {
			if i > 0 {
				expected.WriteByte(' ')
			}
			expected.WriteString(strconv.Itoa(v))
		}
		want := expected.String()

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))

		got, err := runProg(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
