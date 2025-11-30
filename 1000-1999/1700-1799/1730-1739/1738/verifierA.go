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

// Embedded test cases from testcasesA.txt. Each line: n followed by n types and n damages.
const testcasesRaw = `7 1 0 1 1 1 1 1 12 19 7 17 5 10 5
2 1 0 10 4
2 1 1 18 4
6 1 1 0 1 1 1 2 18 1 3 13 1
8 1 0 1 0 0 0 0 0 18 15 3 3 11 17 16 4
5 1 0 1 0 1 15 3 20 13 11
4 1 0 0 0 2 20 9 16
2 0 0 5 2
2 1 1 17 8
4 1 1 1 1 12 3 11 20
2 1 1 7 8
1 1 4
4 1 0 1 1 2 4 5 8
1 0 1
2 0 0 13 3
6 0 0 0 0 0 0 16 7 2 1 18 14
2 1 0 8 3
5 1 1 0 0 1 2 20 4 13 7
5 1 1 0 0 0 6 6 11 17 9
2 1 0 1 16
7 1 1 1 1 0 0 1 3 11 2 18 9 5 8
8 1 1 1 0 1 1 1 0 1 20 7 11 6 8 8 15
7 1 0 1 1 0 0 1 3 9 6 15 17 16 18
1 0 16
6 1 1 0 1 0 0 5 1 13 14 11 1
4 0 0 0 0 4 20 7 10
5 0 0 1 1 0 1 9 15 4 9
3 1 0 0 9 1 2
1 0 9
6 1 0 1 1 1 1 18 6 7 13 19 10
1 0 5
5 1 1 1 0 1 20 2 2 9 6
3 1 1 1 18 5 10
2 1 0 2 10
3 0 1 1 11 10 14
2 0 1 16 11
6 0 1 0 1 1 0 10 11 5 6 19 13
2 0 0 7 8
1 1 1
2 1 1 15 16
4 1 0 1 0 9 19 6 14
4 1 0 0 0 17 15 7 4
8 1 1 0 0 0 0 0 0 15 13 12 18 5 4 20 16
3 1 1 1 11 16 16
4 0 0 1 1 11 2 17 5
5 0 1 1 1 0 3 17 2 3 8
3 0 1 0 15 11 6
3 1 1 1 17 17 2
2 0 1 7 10
7 1 1 0 0 0 0 1 17 19 9 11 3 16 9
5 1 1 1 0 0 5 8 10 11 2
1 1 14
3 1 0 0 12 14 2
8 1 1 0 0 1 0 0 0 20 20 5 11 4 18 12 7
7 1 0 0 1 1 0 1 5 13 10 4 17 7 2
7 1 1 0 1 1 0 0 2 16 9 1 17 19 19
4 0 0 1 1 4 5 14 19
7 0 0 1 0 0 1 0 1 15 14 14 1 16 11
5 0 1 0 0 1 1 12 12 6 1
4 1 0 0 0 1 7 4 1
5 1 0 0 0 0 15 4 16 12 9
3 0 0 1 11 16 10
5 1 0 0 0 1 6 13 5 5 8
6 0 0 0 1 1 1 2 5 20 1 13 3
2 0 1 10 18
7 0 1 1 1 0 0 1 12 17 2 13 14 1 14
6 1 0 1 1 1 0 6 4 9 4 18 20
3 1 1 0 14 14 6
4 1 1 0 1 15 3 16 7
5 0 1 1 0 0 10 4 10 18 20
3 1 1 0 16 8 18
7 1 0 0 1 0 0 1 13 17 19 13 15 4 9
6 1 0 0 0 0 1 10 18 11 4 17 8
3 0 1 1 10 17 5
4 0 1 1 1 10 15 12 19
3 0 0 0 13 13 19
8 0 1 1 1 1 0 1 1 17 11 16 2 15 10 5 16
1 0 1
6 1 1 0 0 0 1 1 12 2 4 20 1
5 1 0 0 1 0 4 14 15 11 13
3 1 1 1 5 15 5
6 0 0 0 1 1 1 14 16 13 8 7 15
4 0 1 0 0 3 6 12 2
3 0 1 0 17 10 12
7 1 0 1 1 1 1 1 7 11 9 2 2 2 6
6 0 1 0 0 0 1 8 20 13 18 8 15
4 1 0 0 1 11 18 15 11
5 0 0 0 1 0 7 17 12 7 7
5 1 1 1 1 1 12 8 2 10 18
2 0 1 16 15
1 1 16
8 1 0 0 0 0 0 0 1 7 15 20 3 14 18 13 2
3 0 1 0 5 9 12
6 1 0 1 0 1 1 17 20 15 18 9 9
4 0 0 0 0 14 8 7 10
1 1 2
2 1 1 4 19
6 0 1 0 0 0 1 11 8 12 16 10 19
3 0 0 1 12 19 1
3 1 0 0 17 3 5`

type testCase struct {
	n   int
	tp  []int
	dmg []int64
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
		expected := 1 + 2*n
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, expected, len(fields))
		}
		tp := make([]int, n)
		dmg := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d parse tp[%d]: %v", idx+1, i, err)
			}
			tp[i] = v
		}
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+n+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d parse dmg[%d]: %v", idx+1, i, err)
			}
			dmg[i] = v
		}
		cases = append(cases, testCase{n: n, tp: tp, dmg: dmg})
	}
	return cases, nil
}

// Embedded solver logic from 1738A.go.
func solve(tc testCase) int64 {
	var fire, frost []int64
	for i, t := range tc.tp {
		if t == 0 {
			fire = append(fire, tc.dmg[i])
		} else {
			frost = append(frost, tc.dmg[i])
		}
	}
	sort.Slice(fire, func(i, j int) bool { return fire[i] > fire[j] })
	sort.Slice(frost, func(i, j int) bool { return frost[i] > frost[j] })
	var sumFire, sumFrost int64
	for _, x := range fire {
		sumFire += x
	}
	for _, x := range frost {
		sumFrost += x
	}
	if len(fire) == len(frost) {
		if len(fire) == 0 {
			return 0
		}
		minVal := fire[len(fire)-1]
		if frost[len(frost)-1] < minVal {
			minVal = frost[len(frost)-1]
		}
		return 2*(sumFire+sumFrost) - minVal
	}
	if len(fire) < len(frost) {
		fire, frost = frost, fire
		sumFire, sumFrost = sumFrost, sumFire
	}
	m := len(frost)
	var extra int64
	for i := 0; i < m; i++ {
		extra += fire[i] + frost[i]
	}
	return sumFire + sumFrost + extra
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		for i, v := range tc.tp {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for i, v := range tc.dmg {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.FormatInt(want, 10) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
