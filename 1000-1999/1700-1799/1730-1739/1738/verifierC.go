package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesC.txt. Each line: n followed by n integers.
const testcasesRaw = `1 -4
2 0 -3
5 -1 4 -2 4 -5
3 1 5 1
6 3 2 3 -1 -5 -5
6 2 0 1 1 3 -3
3 -2 -2 -5
3 0 -3 -3
6 3 5 3 -3 2 1
6 4 0 0 2 -3 1
8 5 3 -2 2 -1 2 3 3
6 5 2 2 0 4 3
8 2 5 -2 0 -3 4 -1 2
5 -1 3 3 3 3
7 -1 -2 2 3 0 5 4
2 0 -5
4 -4 -5 4 5
1 -1
4 5 -4 3 -3
5 -2 -2 -5 1 -5
1 0
6 -3 -2 5 -5 -4 -4
2 -5 -5
1 0
5 -3 -3 -3 3 -5
7 4 -5 -2 -3 -5 -5 0
2 -1 0
8 -5 -1 2 3 4 -5 -1 1
3 2 -2 -4
6 -4 -5 2 -3 3 4
7 2 3 0 -3 0 -1 -1
7 5 -5 3 -3 5 -5 -1
1 -3
3 -3 -4 2
4 3 -5 -2 -2
8 -4 -1 -4 4 -2 4 4 0
5 5 1 -1 3 -5
3 -5 1 1
3 -4 3 -4
4 -4 -4 -5 -3
4 -4 -2 -5 3
8 2 -1 3 5 1 -2 5 -2
7 1 3 -5 4 4 -5 1
3 -4 5 2
6 -5 3 -4 4 0 -1
6 -1 -5 5 1 -4 -4
5 -2 5 -5 2 -5
7 5 2 2 -2 4 4 -4
1 -1
1 0
5 -4 -2 2 -2 -4
6 1 2 -3 0 1 -4
5 -4 -4 -4 4 0
7 -2 -4 -5 4 5 2 -5
8 -1 0 2 -3 0 -1 2 3
8 1 2 5 -1 1 -2 -3 2
5 3 1 5 -4 4
2 -4 0
3 3 -3 1
2 -4 5
1 -3
5 1 -2 5 5 0
8 -3 3 -1 -4 -3 3 1 -4
6 3 -2 3 -1 -3 -3
8 -2 1 0 4 -3 2 2 -5
7 -3 1 3 -5 2 -1 1
5 1 5 2 0 3
6 5 -4 -2 3 4 -2
7 5 1 5 -5 0 2 3
8 5 -3 -4 -5 1 -2 4 4
7 -2 -4 1 3 -2 -1 4
4 2 4 -3 -5
7 2 -1 3 4 -3 2 -2
2 0 -5
8 3 5 5 -4 4 2 5 0
8 -1 3 2 -5 -4 4 0 -3
7 -1 5 5 -3 -5 -3 2
7 2 5 -1 -3 -5 -1 3
8 -5 0 -5 3 1 4 2 -2
5 2 5 -3 2 3
5 -4 -1 0 -1 0
5 5 5 1 3 -4
4 1 4 3 -3
2 -1 -5
4 2 3 -2 3
5 -5 -4 -4 5 1
6 -2 0 0 -4 0 2
6 -3 2 2 -1 2 -3
8 5 -2 -1 0 -3 -4 -2 2
4 5 0 -3 0
3 -3 -2 -1
7 1 0 -1 4 3 4 0
7 5 -1 3 4 5 5 -4
6 -1 1 2 -3 -1 0
8 2 -4 -3 0 1 -3 -5 -4
6 -3 0 -4 5 1 -5
6 -2 4 1 3 -1 2
3 0 0 -2
8 -4 -3 -2 0 -1 -3 1 0
5 -4 0 -2 -2 -2`

type testCase struct {
	n   int
	arr []int
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d parse arr[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return cases, nil
}

var memo [105][105][2][2]int8

// Embedded solver logic from 1738C.go.
func dfs(e, o, turn, parity int) bool {
	if e == 0 && o == 0 {
		return parity == 0
	}
	m := &memo[e][o][turn][parity]
	if *m != 0 {
		return *m == 1
	}
	var res bool
	if turn == 0 {
		if e > 0 && dfs(e-1, o, 1, parity) {
			res = true
		} else if o > 0 && dfs(e, o-1, 1, parity^1) {
			res = true
		}
	} else {
		res = true
		if e > 0 && !dfs(e-1, o, 0, parity) {
			res = false
		}
		if res && o > 0 && !dfs(e, o-1, 0, parity) {
			res = false
		}
	}
	if res {
		*m = 1
	} else {
		*m = 2
	}
	return res
}

func solve(tc testCase) string {
	e, o := 0, 0
	for _, x := range tc.arr {
		if x%2 == 0 {
			e++
		} else {
			o++
		}
	}
	for i := 0; i <= e; i++ {
		for j := 0; j <= o; j++ {
			memo[i][j][0][0] = 0
			memo[i][j][0][1] = 0
			memo[i][j][1][0] = 0
			memo[i][j][1][1] = 0
		}
	}
	if dfs(e, o, 0, 0) {
		return "Alice"
	}
	return "Bob"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

	for idx, tc := range cases {
		want := solve(tc)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
