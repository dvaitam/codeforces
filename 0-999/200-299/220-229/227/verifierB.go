package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const embeddedTestcases = `1 1 2 1 1
3 1 2 3 5 3 1 3 1 3
3 1 3 2 9 2 3 2 3 2 1 1 2 2
6 5 1 3 2 6 4 4 2 1 2 3
3 3 2 1 9 3 3 1 2 2 3 3 2 3
6 1 6 5 2 4 3 9 2 4 3 4 5 5 3 6 4
8 3 1 2 4 7 8 5 6 6 3 5 8 5 5 7
5 1 5 3 4 2 10 1 3 1 2 1 1 5 1 3 5
4 2 4 3 1 5 2 2 1 4 1
1 1 6 1 1 1 1 1 1
2 2 1 1 1
6 4 1 6 5 2 3 9 6 1 4 5 1 2 2 1 1
6 6 4 2 3 1 5 1 3
8 7 8 6 2 5 4 3 1 2 6 2
1 1 8 1 1 1 1 1 1 1 1
10 2 5 4 10 8 9 6 3 1 7 3 2 8 4
9 3 5 7 9 8 6 2 4 1 10 4 6 5 7 5 9 1 3 1 7
7 7 3 4 6 5 1 2 2 1 2
4 2 3 4 1 9 4 4 3 4 2 2 4 4 1
10 4 6 3 9 2 8 5 7 1 10 6 1 9 2 10 6 5
6 2 6 5 4 1 3 5 2 6 1 4 1
7 5 1 3 2 7 4 6 1 3
1 1 6 1 1 1 1 1 1
10 10 5 9 1 4 3 2 8 7 6 2 2 10
6 5 1 3 2 4 6 1 5
8 8 5 2 3 4 7 6 1 3 6 5 8
9 2 5 1 9 3 6 4 7 8 8 5 9 7 2 2 2 6 3
9 5 9 4 2 6 8 1 7 3 4 6 8 3 9
5 5 4 3 2 1 2 3 5
4 2 4 1 3 8 2 4 3 2 4 4 1 4
3 3 1 2 8 2 2 2 3 3 2 3 3
8 4 7 8 2 1 3 5 6 7 7 1 6 8 8 3 2
1 1 7 1 1 1 1 1 1 1
10 9 6 3 1 2 5 7 8 4 10 8 5 9 10 3 8 4 2 6
1 1 8 1 1 1 1 1 1 1 1
10 8 5 4 1 2 10 9 7 3 6 7 8 5 3 1 5 9 8
1 1 6 1 1 1 1 1 1
3 1 3 2 2 2 2
5 1 2 4 5 3 9 1 5 2 4 5 5 2 5 1
5 3 5 4 2 1 9 3 1 1 1 4 3 2 3 3
2 1 2 8 2 1 2 2 2 2 1 2
4 4 1 2 3 2 2 4
4 4 2 1 3 3 2 2 3
9 1 2 4 5 3 6 8 9 7 7 5 9 2 6 5 7 8
3 1 3 2 8 2 1 1 2 2 1 1 1
6 4 6 5 1 3 2 1 5
6 1 6 3 4 5 2 8 6 2 3 3 2 4 1 2
4 4 1 2 3 7 3 3 1 3 2 2 2
10 4 3 7 2 10 5 8 9 6 1 7 8 5 3 6 9 10 2
6 2 3 1 4 5 6 7 4 4 5 3 5 5 5
2 1 2 7 2 1 2 2 2 1 1
8 8 4 7 2 5 6 1 3 5 8 1 5 2 6
4 4 3 1 2 7 1 3 4 1 4 2 1
8 7 8 4 6 2 1 5 3 9 1 6 4 3 6 8 1 3 2
4 3 4 2 1 1 2
7 6 1 7 2 4 5 3 3 6 2 2
7 5 7 1 4 6 3 2 3 4 3 7
5 4 3 2 5 1 5 5 4 3 2 4
3 2 1 3 1 3
9 7 6 3 1 5 8 2 9 4 1 5
8 5 4 8 3 7 2 6 1 2 3 3
2 1 2 5 1 1 2 2 1
7 3 5 2 7 1 6 4 5 1 2 4 3 3
9 1 3 5 9 6 2 8 4 7 7 8 5 4 8 8 4 8
10 7 10 1 4 9 3 8 2 5 6 5 4 9 5 2 1
1 1 4 1 1 1 1
6 3 5 2 6 1 4 5 5 5 2 2 2
3 1 2 3 2 3 3
8 2 3 7 8 4 1 6 5 10 5 4 4 6 7 7 4 4 1 5
4 1 4 3 2 7 1 4 4 4 4 4 3
4 4 3 1 2 9 1 1 1 4 4 4 2 3 1
6 2 1 6 4 3 5 8 6 6 2 3 1 1 4 1
3 2 1 3 4 2 1 3 1
10 8 4 10 1 7 9 6 2 5 3 7 3 1 5 9 5 8 1
9 4 2 7 3 5 1 8 9 6 6 5 8 5 9 3 1
1 1 6 1 1 1 1 1 1
7 7 1 3 5 6 2 4 1 5
6 5 4 1 3 2 6 4 5 2 4 1
6 2 3 4 6 5 1 9 2 1 1 6 4 1 5 3 5
4 4 2 3 1 7 2 2 4 2 4 3 1
10 10 1 6 2 7 5 4 8 9 3 6 3 1 5 3 3 7
10 3 1 6 10 7 2 4 9 8 5 6 1 7 10 1 7 5
1 1 9 1 1 1 1 1 1 1 1 1
9 2 5 1 7 6 8 4 3 9 7 2 7 8 9 8 2 7
8 4 8 1 2 3 5 7 6 3 7 8 8
3 1 3 2 6 3 2 2 3 3 2
6 2 5 3 4 1 6 4 4 1 1 1
6 1 2 4 6 5 3 3 6 3 1
9 6 9 3 1 4 2 5 7 8 6 1 5 8 3 1 1
6 3 5 2 1 6 4 4 6 5 3 2
6 1 6 2 5 4 3 2 5 2
10 9 1 8 5 2 3 7 10 4 6 4 10 8 4 7
5 4 1 3 2 5 1 4
8 3 1 4 5 7 2 8 6 3 4 4 1
10 7 1 5 9 3 10 6 4 8 2 1 6
2 1 2 1 2
8 3 5 6 2 4 8 1 7 2 5 3
8 8 4 1 6 3 5 7 2 2 7 3
3 1 2 3 8 1 1 3 2 2 1 1 1
10 3 4 1 9 10 8 7 2 5 6 2 5 2`

type testCase struct {
	n       int
	perm    []int
	queries []int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	res := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d too short", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d n: %v", idx+1, err)
		}
		if len(fields) < 1+n+1 {
			return nil, fmt.Errorf("line %d missing permutation entries", idx+1)
		}
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			perm[i], err = strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d perm %d: %v", idx+1, i, err)
			}
		}
		mIdx := 1 + n
		m, err := strconv.Atoi(fields[mIdx])
		if err != nil {
			return nil, fmt.Errorf("line %d m: %v", idx+1, err)
		}
		if len(fields) != 1+n+1+m {
			return nil, fmt.Errorf("line %d length mismatch", idx+1)
		}
		queries := make([]int, m)
		for i := 0; i < m; i++ {
			queries[i], err = strconv.Atoi(fields[mIdx+1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d query %d: %v", idx+1, i, err)
			}
		}
		res = append(res, testCase{n: n, perm: perm, queries: queries})
	}
	return res, nil
}

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n int, perm []int, queries []int) (int64, int64) {
	pos := make([]int, n+1)
	for i, v := range perm {
		if v >= 1 && v <= n {
			pos[v] = i + 1
		}
	}
	var vasya, petya int64
	for _, q := range queries {
		if q >= 1 && q <= n {
			p := pos[q]
			vasya += int64(p)
			petya += int64(n - p + 1)
		}
	}
	return vasya, petya
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range cases {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.perm {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		input.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
		for i, v := range tc.queries {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		vasya, petya := solve(tc.n, tc.perm, tc.queries)
		want := fmt.Sprintf("%d %d", vasya, petya)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
