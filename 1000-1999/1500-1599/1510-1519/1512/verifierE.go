package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `8 5 6 14
6 5 5 1
1 1 1 1
8 8 8 1
7 2 3 3
8 3 8 27
7 2 2 5
1 1 1 1
7 1 4 10
7 2 6 22
4 3 3 2
4 1 1 4
4 1 3 6
7 5 6 7
8 4 6 13
8 8 8 3
8 4 6 13
5 5 5 2
8 1 8 36
4 3 3 2
5 5 5 5
6 5 6 6
1 1 1 1
1 1 1 1
2 2 2 1
3 1 1 3
7 3 7 18
1 1 1 1
5 4 5 5
1 1 1 1
6 3 4 6
2 2 2 1
4 3 3 1
1 1 1 1
5 3 3 3
7 4 4 1
4 3 3 1
2 1 2 3
7 4 7 14
1 1 1 1
7 6 6 5
5 5 5 3
1 1 1 1
8 5 5 4
3 3 3 2
1 1 1 1
7 5 5 1
1 1 1 1
2 1 1 1
6 3 3 2
6 5 5 5
5 3 5 12
8 5 6 8
3 1 3 6
6 4 6 9
7 7 7 6
6 1 1 4
8 1 3 16
5 2 4 11
5 1 1 2
2 2 2 2
2 1 1 2
1 1 1 1
1 1 1 1
4 4 4 1
2 1 1 1
2 2 2 2
1 1 1 1
5 5 5 4
3 2 3 4
1 1 1 1
3 1 2 4
7 6 7 12
7 5 6 3
5 5 5 5
4 1 2 5
7 3 3 7
4 4 4 4
3 2 3 5
7 7 7 7
4 4 4 2
2 1 2 3
1 1 1 1
4 4 4 1
5 5 5 2
8 4 7 17
1 1 1 1
3 1 2 3
3 2 3 5
8 2 6 20
1 1 1 1
3 3 3 3
1 1 1 1
4 2 3 6
4 4 4 2
7 1 2 12
4 3 4 3
7 4 7 15
7 2 5 20
8 5 7 16`

type testCase struct {
	n int
	l int
	r int
	s int
}

func solveCase(tc testCase) []int {
	n, l, r, s := tc.n, tc.l-1, tc.r-1, tc.s
	k := r - l + 1
	for first := 1; first+k-1 <= n; first++ {
		sum := k*first + k*(k-1)/2
		if sum > s {
			break
		}
		extra := s - sum
		if extra <= k {
			ans := make([]int, n)
			used := make([]bool, n+1)
			needAdd := r - extra + 1
			ok := true
			for i := l; i <= r; i++ {
				val := first + (i - l)
				if i >= needAdd {
					val++
				}
				if val > n {
					ok = false
					break
				}
				ans[i] = val
				used[val] = true
			}
			if !ok {
				continue
			}
			cur := 1
			for i := 0; i < n; i++ {
				if i >= l && i <= r {
					continue
				}
				for cur <= n && used[cur] {
					cur++
				}
				ans[i] = cur
				used[cur] = true
			}
			return ans
		}
	}
	return nil
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
		if len(fields) != 4 {
			return nil, fmt.Errorf("case %d: malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", idx+1, err)
		}
		l, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse l: %v", idx+1, err)
		}
		r, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse r: %v", idx+1, err)
		}
		s, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse s: %v", idx+1, err)
		}
		res = append(res, testCase{n: n, l: l, r: r, s: s})
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.l, tc.r, tc.s))

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected == nil {
			if strings.TrimSpace(got) != "-1" {
				fmt.Printf("case %d failed\nexpected: -1\ngot: %s\n", i+1, got)
				os.Exit(1)
			}
			continue
		}
		gotVals := strings.Fields(strings.TrimSpace(got))
		if len(gotVals) != tc.n {
			fmt.Printf("case %d failed\nexpected %d numbers\ngot %d numbers\n", i+1, tc.n, len(gotVals))
			os.Exit(1)
		}
		for idx, exp := range expected {
			val, err := strconv.Atoi(gotVals[idx])
			if err != nil || val != exp {
				fmt.Printf("case %d failed\nexpected: %v\ngot: %s\n", i+1, expected, got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
