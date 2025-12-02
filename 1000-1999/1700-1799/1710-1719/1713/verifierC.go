package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `20
17
5
3
11
8
6
8
1
6
18
6
3
14
20
4
20
15
5
20
20
2
9
11
13
1
2
16
3
12
10
5
15
8
17
12
6
13
11
9
16
13
1
10
17
10
18
16
2
18
19
18
9
2
15
13
4
13
12
16
2
1
9
2
9
19
10
7
17
17
11
13
9
7
4
19
11
8
19
18
12
6
5
11
1
19
2
19
5
12
12
10
10
11
16
13
20
14
6
1`

type testCase struct {
	n int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		cases = append(cases, testCase{n: n})
	}
	return cases, nil
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

func runExe(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isSquare(x int) bool {
	if x < 0 {
		return false
	}
	// Only simple integer check needed
	r := 0
	for r*r < x {
		r++
	}
	return r*r == x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
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

	for idx, tc := range cases {
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", tc.n))

		got, err := runExe(bin, input.String())
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", idx+1, err, input.String())
			os.Exit(1)
		}

		// Validate 'got'
		// It should be n distinct integers from 0 to n-1.
		// And for each i, p[i]+i is square.
		fields := strings.Fields(got)
		if len(fields) != tc.n {
			fmt.Printf("test %d failed: expected %d numbers, got %d\ninput:\n%s\ngot:\n%s\n", idx+1, tc.n, len(fields), input.String(), got)
			os.Exit(1)
		}

		seen := make(map[int]bool)
		p := make([]int, tc.n)
		for i, s := range fields {
			val, err := strconv.Atoi(s)
			if err != nil {
				fmt.Printf("test %d failed: invalid number %q at index %d\n", idx+1, s, i)
				os.Exit(1)
			}
			if val < 0 || val >= tc.n {
				fmt.Printf("test %d failed: number %d out of range [0, %d)\n", idx+1, val, tc.n)
				os.Exit(1)
			}
			if seen[val] {
				fmt.Printf("test %d failed: duplicate number %d\n", idx+1, val)
				os.Exit(1)
			}
			seen[val] = true
			p[i] = val
		}

		for i, val := range p {
			if !isSquare(val + i) {
				fmt.Printf("test %d failed: p[%d] + %d = %d is not a perfect square\n", idx+1, i, i, val+i)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}