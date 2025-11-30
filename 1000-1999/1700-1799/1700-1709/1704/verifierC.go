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

const testcasesRaw = `6 5 5 2 1 6 4
11 2 9 8
8 1 1
10 2 7 3
9 2 2 4
8 3 4 7 6
10 4 6 8 4 9
10 4 9 2 8 3
5 2 1 4
9 1 1
4 2 4 2
11 3 3 11 9
4 3 1 2 3
6 5 6 4 1 3 5
9 2 3 6
6 1 5
11 2 3 7
12 1 9
6 4 2 1 6 3
11 5 11 9 2 4 8
10 1 10
3 2 1 3
10 1 9
6 1 1
7 4 3 6 4 2
12 2 9 6
11 4 9 7 11 3
9 4 4 8 7 3
8 2 5 8
7 2 7 6
12 1 12
8 3 3 8 7
7 3 4 3 5
10 1 3
5 3 2 5 1
12 5 10 4 9 7 11
12 2 9 8
9 2 2 9
5 1 1
9 4 7 3 5 9
5 5 5 1 4 2 3
7 2 6 7
9 3 3 4 9
5 3 4 3 1
11 3 4 8 1
7 5 5 1 7 3 2
7 5 1 7 3 2 6
4 1 4
12 2 12 4
11 5 7 2 4 11 6
11 2 10 5
3 3 1 3 2
10 5 10 4 5 1 6
5 5 5 2 4 3 1
9 3 8 2 6
5 2 5 1
10 1 8
6 4 6 5 3 1
4 1 1
9 4 4 3 5 7
6 5 4 5 3 1 6
8 5 2 3 1 4 6
12 2 3 5
10 1 10
11 1 10
9 1 7
11 5 11 5 7 10 3
10 1 9
10 1 7
7 5 6 3 2 7 1
12 3 7 12 9
3 3 3 1 2
12 5 1 2 6 10 9
11 1 11
8 5 2 4 6 1 7
8 5 1 2 3 6 7
5 5 2 1 5 3 4
8 3 5 8 3
3 3 1 3 2
9 5 5 2 1 6 8
7 4 1 2 3 5
6 2 1 3
10 1 9
7 2 7 5
4 3 3 2 1
3 2 2 1
3 2 2 1
11 1 10
11 4 3 4 10 2
12 2 10 9
4 3 4 1 3
8 4 6 5 8 3
6 1 1
10 1 3
7 5 1 7 2 6 3
5 1 5
6 2 4 3
6 4 5 3 6 2
4 2 2 1
12 3 10 7 12`

type testCase struct {
	n   int64
	m   int
	arr []int64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: invalid", i+1)
		}
		n, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", i+1, err)
		}
		if len(parts) != 2+m {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, 2+m, len(parts))
		}
		arr := make([]int64, m)
		for j := 0; j < m; j++ {
			v, err := strconv.ParseInt(parts[2+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse arr[%d]: %v", i+1, j, err)
			}
			arr[j] = v
		}
		cases = append(cases, testCase{n: n, m: m, arr: arr})
	}
	return cases, nil
}

func solve(tc testCase) int64 {
	arr := append([]int64(nil), tc.arr...)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	gaps := make([]int64, 0, tc.m)
	for i := 0; i < tc.m-1; i++ {
		gap := arr[i+1] - arr[i] - 1
		if gap > 0 {
			gaps = append(gaps, gap)
		}
	}
	lastGap := tc.n - arr[tc.m-1] + arr[0] - 1
	if lastGap > 0 {
		gaps = append(gaps, lastGap)
	}
	sort.Slice(gaps, func(i, j int) bool { return gaps[i] > gaps[j] })

	infected := int64(tc.m)
	days := int64(0)
	for _, g := range gaps {
		g -= 2 * days
		if g <= 0 {
			continue
		}
		if g == 1 {
			infected += 1
			days += 1
		} else {
			infected += g - 1
			days += 2
		}
	}
	return infected
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')

		expect := fmt.Sprint(solve(tc))
		got, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
