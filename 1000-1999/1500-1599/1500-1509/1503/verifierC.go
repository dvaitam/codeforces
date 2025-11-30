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

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `100
6
1 1
9 10
5 5
9 8
8 3
6 5
7
1 9
9 1
3 9
6 5
4 5
1 10
7 10
6
4 8
6 10
3 1
8 1
2 4
2 1
10
9 8
9 8
1 6
8 1
7 7
1 8
6 8
2 5
4 3
4 8
3
1 7
10 3
5 9
2
8 9
8 3
10
1 2
2 3
1 6
9 1
2 6
2 4
5 5
5 4
6 6
10 4
6
4 10
9 1
3 6
7 5
9 6
8 2
3
2 9
3 8
2 3
3
6 4
7 3
3 3
7
1 3
2 2
9 9
8 1
10 6
1 5
6 5
4
9 2
5 4
2 9
5 9
8
4 2
3 2
6 7
2 3
2 4
10 9
3 1
4 9
2
7 3
10 6
9
2 7
8 10
9 9
9 1
3 1
1 5
10 9
8 1
10 2
6
7 6
5 9
10 9
3 8
10 3
7 1
5
10 2
9 4
1 5
4 6
10 10
7
8 8
4 3
6 2
6 7
3 2
3 3
2 10
4
6 9
7 9
5 3
10 4
3
6 4
9 3
1 6
2
5 6
4 7
5
2 3
2 5
3 3
5 5
8 6
8
5 9
5 4
6 9
10 1
6 4
2 1
5 2
9 7
5
6 3
9 9
4 2
2 5
4 8
10
2 10
10 8
10 5
8 5
6 1
3 1
10 9
1 7
1 2
10 4
6
5 5
10 1
4 7
6 10
1 6
5 6
4
7 3
7 9
7 7
6 9
10
1 10
7 5
9 7
4 9
3 9
1 7
7 7
6 1
4 8
3 9
1
6 1
9
5 6
7 8
2 1
5 7
9 4
4 10
10 9
10 1
4 8
1
6 4
8
1 9
2 5
7 9
5 9
1 5
4 10
7 2
6 5
6
6 7
1 1
10 10
10 10
8 4
6 6
9
5 7
2 3
10 7
6 5
1 9
1 2
4 1
4 1
4 3
3
9 2
9 3
10 4
5
8 6
10 3
7 6
5 5
7 5
10
6 2
5 4
9 5
6 5
5 9
1 10
2 9
4 2
7 5
6 7
4
3 2
9 9
4 4
7 10
7
4 9
1 4
6 8
8 6
7 4
9 4
2 1
8
6 2
6 5
1 2
1 10
10 6
10 5
3 10
4 3
7
3 2
6 3
4 5
5 3
1 5
10 10
6 9
4
6 1
3 1
10 8
3 5
8
6 3
5 6
7 8
7 2
2 10
5 5
9 4
7 9
9
10 5
1 8
5 8
7 9
6 1
9 1
10 10
5 9
8 10
8
10 5
5 10
4 9
10 7
2 6
6 5
6 5
8 5
6
7 5
3 8
5 4
1 10
9 2
4 1
5
4 8
8 7
4 3
3 7
5 1
3
9 3
5 5
9 9
3
1 4
1 4
7 6
8
8 10
1 4
3 3
2 5
9 6
7 9
3 4
8 10
7
6 10
7 5
3 9
2 5
9 5
10 3
10 2
2
1 8
6 5
7
1 5
5 3
10 5
5 6
7 9
6 8
10 5
8
9 5
6 10
2 8
2 5
9 10
7 9
4 1
2 2
3
9 5
6 8
4 4
1
9 5
1
3 8
3
10 2
9 2
6 1
9
8 7
3 6
3 7
1 10
5 6
9 10
6 1
9 8
3 2
3
1 9
3 8
7 7
7
7 10
3 10
8 4
1 8
6 3
1 3
8 6
2
5 6
10 7
8
3 10
2 2
6 7
2 3
1 5
2 2
1 1
3 5
1
3 6
7
10 4
6 5
2 2
6 1
6 3
8 8
2 6
5
9 5
3 9
2 3
1 5
1 2
4
6 2
7 10
4 3
8 1
1
3 5
3
10 10
8 1
6 6
7
8 4
9 6
5 7
2 2
1 4
4 4
9 10
8
7 6
4 7
8 10
6 1
10 8
7 7
7 7
10 7
8
9 10
4 9
9 5
9 8
8 9
4 1
8 10
6 4
10
8 8
3 3
2 10
5 1
6 9
6 4
10 4
2 9
8 8
8 2
7
8 10
4 10
3 10
1 2
4 7
6 5
2 8
9
6 9
2 6
8 6
1 5
3 5
6 3
9 9
5 3
10 10
8
7 1
6 2
3 10
4 3
9 3
2 10
8 6
6 2
4
6 5
8 2
5 9
4 4
1
6 2
5
1 3
3 8
1 8
5 6
9 3
5
9 8
5 5
10 1
9 3
4 9
8
7 3
3 4
5 4
9 4
8 8
8 6
7 9
4 7
9
6 3
1 8
4 8
7 5
1 3
7 7
1 9
8 10
5 8
9
7 10
9 2
5 5
4 3
4 3
9 10
9 3
9 1
7 10
4
8 2
3 4
10 7
4 9
3
3 4
4 10
2 6
3
10 2
3 4
9 7
3
4 1
9 5
4 10
1
9 8
1
10 1
5
3 1
3 6
9 6
1 10
6 9
5
10 1
6 2
8 2
5 8
7 4
5
4 10
8 2
9 4
5 4
2 6
6
5 2
8 8
3 3
9 7
8 9
2 10
8
10 8
8 1
7 4
5 1
2 1
7 8
7 5
7 10
10
10 5
4 8
10 8
4 4
4 10
6 5
10 7
8 8
1 3
9 7
5
9 8
2 7
3 5
3 10
5 3
2
3 3
1 10
5
6 6
10 3
9 3
7 7
5 8
2
9 4
8 2
2
10 4
5 5`

type testCase struct {
	n      int
	cities []struct{ a, c int }
}

func solveCase(tc testCase) int64 {
	cities := make([]struct{ a, c int }, len(tc.cities))
	copy(cities, tc.cities)
	sort.Slice(cities, func(i, j int) bool { return cities[i].a < cities[j].a })
	if len(cities) == 0 {
		return 0
	}
	var ans int64
	ans += int64(cities[0].c)
	reach := int64(cities[0].a + cities[0].c)
	for i := 1; i < len(cities); i++ {
		if int64(cities[i].a) > reach {
			ans += int64(cities[i].a) - reach
			reach = int64(cities[i].a)
		}
		ans += int64(cities[i].c)
		if int64(cities[i].a+cities[i].c) > reach {
			reach = int64(cities[i].a + cities[i].c)
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	res := make([]testCase, 0, t)
	idx := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseIdx+1, err)
		}
		idx++
		if idx+2*n > len(fields) {
			return nil, fmt.Errorf("case %d: not enough city data", caseIdx+1)
		}
		cities := make([]struct{ a, c int }, n)
		for i := 0; i < n; i++ {
			a, err := strconv.Atoi(fields[idx+2*i])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a%d: %v", caseIdx+1, i+1, err)
			}
			c, err := strconv.Atoi(fields[idx+2*i+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse c%d: %v", caseIdx+1, i+1, err)
			}
			cities[i] = struct{ a, c int }{a: a, c: c}
		}
		idx += 2 * n
		res = append(res, testCase{n: n, cities: cities})
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra data after parsing")
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

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, city := range tc.cities {
			sb.WriteString(strconv.Itoa(city.a))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(city.c))
			sb.WriteByte('\n')
		}
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil || val != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
