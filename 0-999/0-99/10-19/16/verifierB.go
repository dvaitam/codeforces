package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `11 4 8 5 12 3 10 1 12 10
18 1 5 6
1 4 2 1 8 1 1 4 11 2
2 3 14 3 7 8 14 3
12 3 6 6 14 7 1 7
9 5 18 8 2 10 4 7 13 3 1 9
5 5 17 3 3 6 8 3 8 1 6 9
6 1 14 10
4 5 15 3 20 10 2 5 11 7 1 1
16 1 12 5
5 4 8 9 12 3 13 6 9 8
13 1 10 9
10 5 16 1 18 10 18 5 2 8 13 2
13 3 16 1 1 5 2 5
19 3 7 9 17 6 13 5
7 1 19 6
8 5 18 6 6 3 11 1 19 1 19 3
12 3 10 5 11 8 13 10
14 2 1 3 19 1
15 2 11 1 16 5
20 2 3 9 14 5
6 5 6 2 6 10 4 9 18 10 13 7
9 3 10 1 14 5 9 9
17 5 11 6 7 7 5 1 17 3 19 7
12 4 2 9 14 10 8 1 12 9
6 2 12 8 1 4
19 2 9 3 14 2
19 4 8 8 17 2 7 3 15 2
14 4 9 5 14 6 20 6 3 5
1 4 1 5 7 7 13 7 13 1
19 4 12 10 5 10 9 6 1 7
16 5 5 1 3 10 12 6 1 2 7 2
18 4 2 6 1 6 13 3 9 7
5 5 5 7 10 9 2 3 5 3 16 1
17 1 18 7
6 3 19 2 3 9 6 5
7 3 11 5 9 9 15 3
15 5 5 1 19 3 17 1 11 2 7 8
20 2 15 9 6 6
5 4 18 1 18 2 17 6 1 2
4 4 20 6 19 8 11 7 17 6
4 2 11 1 6 3
1 3 20 4 2 7 2 5
13 1 20 3
12 1 14 1
15 3 20 10 9 5 19 8
14 2 1 8 9 4
13 1 12 2
4 1 12 1
6 4 20 1 11 8 18 8 16 2
2 5 13 5 1 9 4 2 11 6 4 8
2 2 17 5 2 1
13 3 6 9 5 3 6 3
8 5 11 1 16 7 2 4 8 5 11 3
8 3 8 3 14 8 12 10
5 4 19 1 6 10 1 7 6 3
1 1 11 9
1 1 2 2
19 5 5 3 13 1 14 7 19 6 8 3
12 5 7 9 13 2 5 7 19 6 4 7
14 2 16 7 8 7
8 4 13 10 3 5 9 9 12 9
1 5 20 8 8 5 2 10 11 7 4 9
2 2 13 1 14 7
14 1 15 10
15 2 6 6 16 7
6 5 10 9 4 6 12 3 12 8 17 1
7 3 6 10 11 5 13 1
10 5 14 1 14 5 13 4 12 3 5 2
20 3 6 1 14 10 13 8
3 1 14 9
18 2 6 3 7 3
8 1 17 3
16 3 20 5 11 2 14 5
6 3 17 6 18 3 13 9
10 2 13 6 13 8
17 3 14 7 4 3 5 1
19 5 17 2 7 10 20 9 4 5 20 3
13 1 2 1
4 3 16 6 4 8 12 10
9 4 8 3 18 9 3 9 6 1
6 5 14 10 7 8 13 5 1 10 5 7
6 4 19 1 13 2 19 7 11 4
17 4 2 8 20 2 9 9 16 9
13 4 9 3 8 9 12 3 10 10
5 4 3 2 16 7 19 7 18 2
9 4 8 2 10 3 12 2 5 1
5 5 18 4 1 1 13 9 20 8 4 8
18 3 11 2 1 4 8 8
10 3 8 1 16 6 17 6
3 1 10 10
14 2 12 7 5 4
10 2 16 6 10 7
20 2 4 7 12 9
16 2 12 6 14 5
12 4 10 2 16 5 4 8 5 6
8 2 11 8 8 2
13 4 15 9 15 10 20 4 13 9
10 4 8 6 17 1 3 8 11 7
8 4 2 10 2 7 3 5 7 6`

type pair struct {
	a, b int64
}

type testCase struct {
	n int64
	m int
	p []pair
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: invalid", i+1)
		}
		n, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", i+1, err)
		}
		if len(fields) != 2+m*2 {
			return nil, fmt.Errorf("line %d: expected %d values got %d", i+1, 2+m*2, len(fields))
		}
		pairs := make([]pair, m)
		for j := 0; j < m; j++ {
			a, err := strconv.ParseInt(fields[2+2*j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a%d: %v", i+1, j, err)
			}
			b, err := strconv.ParseInt(fields[3+2*j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse b%d: %v", i+1, j, err)
			}
			pairs[j] = pair{a: a, b: b}
		}
		cases = append(cases, testCase{n: n, m: m, p: pairs})
	}
	return cases, nil
}

func solve(tc testCase) int64 {
	p := append([]pair(nil), tc.p...)
	sort.Slice(p, func(i, j int) bool { return p[i].b > p[j].b })
	remaining := tc.n
	var total int64
	for _, pr := range p {
		if remaining == 0 {
			break
		}
		take := pr.a
		if take > remaining {
			take = remaining
		}
		total += take * pr.b
		remaining -= take
	}
	return total
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

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, pr := range tc.p {
			sb.WriteString(fmt.Sprintf("%d %d\n", pr.a, pr.b))
		}
		expect := fmt.Sprint(solve(tc))
		got, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
