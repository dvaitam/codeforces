package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
3
0 1 0
8
1 1 1 0 0 1 0 1
7
0 1 1 0 0 1 0
1
0
9
0 1 0 1 0 0 1 1 0
6
0 0 1 1 0 1
9
0 0 1 0 1 1 0 1 1
10
1 1 0 1 0 1 1 0 1 1
2
1 0
3
0 1 0
7
1 0 0 0 1 0 0
9
0 0 1 1 0 0 0 1 0
4
1 0 0 0
10
0 1 1 1 1 0 1 1 1 0
8
0 0 0 1 1 1 1 0
2
0 1
10
1 1 0 1 1 0 0 1 0 0
3
1 1 1
7
0 0 0 1 0 0 0
6
1 1 1 0 0 1
10
1 0 1 1 0 1 1 1 0 1
10
1 0 0 1 1 0 1 1 1 1
1
0
1
1
8
0 0 0 1 1 1 0 0
7
1 0 0 0 0 0 0
2
0 0
4
0 0 0 1
9
1 1 0 1 0 1 1 0 0
3
0 0 0
3
1 1 0
4
0 0 0 0
4
0 0 0 1
5
1 0 1 1 1
6
0 1 1 1 0 0
4
0 1 0 0
9
1 1 0 0 1 0 0 1 1
10
1 1 1 0 0 1 0 0 1 1
9
1 1 1 0 0 0 1 0 0
9
0 1 1 0 1 1 0 0 1
7
1 1 1 1 1 0 1
9
0 0 0 0 1 1 1 0 1
6
1 1 1 0 0 1
5
1 0 1 1 0
7
0 1 1 0 1 0 0
3
0 0 1
5
0 0 0 0 1
4
1 0 1 0
7
0 1 1 1 0 0 0
4
0 0 1 0
10
1 0 0 0 0 1 1 0 0 0
2
0 0
6
0 1 0 1 1 0
2
0 1
4
0 1 0 0
10
1 1 0 1 0 0 1 1 0 1
6
0 0 1 1 1 1
5
0 1 1 0 0
7
1 1 0 1 0 1 1
6
0 1 1 0 0 0
3
1 1 1
3
1 0 0
5
0 0 0 0 0
3
1 0 1
10
1 1 1 1 1 1 1 1 1 1
6
1 0 0 1 1 0
7
1 0 1 0 1 1 1
8
0 1 1 1 1 1 1 1
3
0 1 1
4
1 0 0 1
10
0 1 1 1 1 0 1 1 1 1
10
1 0 0 1 1 1 0 0 1 1
8
0 1 0 1 1 0 1 0
5
0 0 1 1 0
8
1 0 1 0 1 1 1 0
6
1 0 1 0 1 0
5
0 0 0 0 1
9
1 1 0 1 1 0 1 0 0
7
0 0 1 0 0 0 1
6
1 1 0 0 1 1
3
1 1 0
10
0 1 1 0 1 1 1 0 1 1
7
0 0 0 1 1 1 0
8
0 0 0 0 1 0 1 1
3
1 0 1
4
0 0 0 0
8
0 0 0 1 1 1 0 1
4
1 0 1 1
6
0 1 0 0 0 1
9
0 0 0 0 1 0 1 1 1
8
1 1 1 1 1 0 1 1
4
1 0 1 1
2
0 0
3
0 0 1
8
0 0 0 1 1 0 0 0
8
1 0 1 0 0 0 1 1
4
1 1 1 0
2
1 1
8
0 1 0 1 1 1 1 1
3
0 1 0
1
1
`

type testCase struct {
	n int
	a []int
}

func solveCase(tc testCase) int64 {
	positions := make([]int, 0)
	for i, v := range tc.a {
		if v == 1 {
			positions = append(positions, i+1)
		}
	}
	total := len(positions)
	if total <= 1 {
		return -1
	}
	// prime factors of total
	factors := make([]int, 0)
	x := total
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			factors = append(factors, p)
			for x%p == 0 {
				x /= p
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	best := int64(-1)
	for _, k := range factors {
		var cost int64
		for i := 0; i < total; i += k {
			median := positions[i+k/2]
			for j := i; j < i+k; j++ {
				if positions[j] > median {
					cost += int64(positions[j] - median)
				} else {
					cost += int64(median - positions[j])
				}
			}
		}
		if best == -1 || cost < best {
			best = cost
		}
	}
	return best
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	idx := 0
	t, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
	if err != nil {
		return nil, err
	}
	idx++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
		if err != nil {
			return nil, err
		}
		idx++
		if idx >= len(lines) {
			return nil, fmt.Errorf("missing array at case %d", i+1)
		}
		arrStr := strings.Fields(strings.TrimSpace(lines[idx]))
		if len(arrStr) != n {
			return nil, fmt.Errorf("case %d expected %d numbers got %d", i+1, n, len(arrStr))
		}
		idx++
		arr := make([]int, n)
		for j, s := range arrStr {
			val, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			arr[j] = val
		}
		cases = append(cases, testCase{n: n, a: arr})
	}
	return cases, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"
		want := strconv.FormatInt(solveCase(tc), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
