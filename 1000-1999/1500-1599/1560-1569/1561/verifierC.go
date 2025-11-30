package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `1;1 1
2;1 1;1 2
2;5 3 9 4 16 15;4 13 7 4 16
1;4 14 20 1 15
3;2 19 4;3 1 1 1;5 1 13 7 14 1
5;2 15 16;5 8 12 8 8 15;3 1 14 18;1 6;3 4 11 17
4;5 7 10 10 19 16;5 13 19 2 16 8;4 14 6 12 18;3 3 15 17
1;2 17 13
3;4 1 16 2 10;5 19 19 13 6 6;5 8 1 7 18 18
2;4 17 12 19 12;4 9 18 20 1
1;4 8 11 16 19
4;2 16 10;1 3;1 10;1 6
5;3 10 14 19;3 3 3 6;3 10 8 18;4 11 7 4 18;5 14 4 1 7 6
5;3 14 11 6;3 6 4 7;2 19 15;2 16 12;5 13 2 18 2 17
2;3 9 10 16;4 8 17 2 10
4;5 9 11 15 20 8;4 1 6 18 16;3 4 2 12;5 18 19 15 15 5
3;3 13 15 19;2 9 20;3 13 9 3
1;5 14 16 10 11 1
3;1 6;1 16;1 13
1;2 19 5
2;5 12 10 18 2 13;3 2 9 3
2;4 19 16 7 12;5 10 3 6 7 2
2;2 18 12;3 13 4 5
5;5 5 3 19 13 4;4 12 17 19 18;3 15 12 11;3 20 7 16;5 10 11 14 5 6
3;3 5 14 6;1 10;3 8 17 18
5;2 18 2;5 7 13 14 11 1;2 15 11;5 6 5 1 12 6;5 13 5 1 9 17
2;1 16;3 1 16 13
4;3 20 16 17;5 9 12 6 13 6;1 12;2 10 17
3;5 8 12 12 12 14;2 19 12;1 16
1;4 2 15 16 15
2;3 16 18 9;2 15 18
5;5 4 14 2 5 20;2 20 16;3 4 8 1;5 10 9 10 19 9;2 15 14
5;1 5;3 9 10 20;2 1 8;1 1;4 10 3 17 1
2;2 6 7;5 13 18 4 5 20
5;4 7 8 15 17;5 2 10 8 14 12;4 18 5 18 10;5 2 9 10 6 16;4 6 1 3 3
2;1 5;3 12 14 7
2;4 8 10 18 12;4 1 14 17 14
5;5 6 3 7 15 9;1 15;2 10 12;4 13 20 8 5;5 9 17 6 9 13
2;4 3 8 9 5;4 6 8 14 20
1;2 5 14
3;4 11 17 2 17;2 3 12;4 2 2 19 13
2;5 10 7 16 12 3;3 8 1 10
5;3 9 8 18;4 9 19 1 6;2 1 12;1 17;3 2 15 17
5;4 12 15 17 12;5 2 7 10 16 18;2 16 6;3 4 7 18;1 20
1;4 7 11 11 12
2;3 20 1 16;2 7 2
3;3 11 14 17;4 5 5 10 17;4 4 1 10 2
2;5 8 3 10 10 14;2 11 16
3;4 8 17 5 1;5 14 19 17 6 11;2 11 20
4;2 13 15;2 2 5;4 11 9 12 9;1 10
3;3 5 14 2;4 12 4 6 13;2 13 13
2;5 20 8 18 9 9;3 16 15 15
2;4 17 11 17 2;2 14 13
5;3 12 1 13;5 4 16 20 14 11;3 4 2 1;4 17 7 14 8;1 17
4;4 16 19 5 11;5 9 6 16 5 6;5 19 16 15 19 10;3 13 12 14
5;4 10 9 6 12;5 10 10 7 13 9;3 5 8 17;5 10 2 2 19 10;1 2
3;1 8;4 19 2 10 2;2 1 10
5;4 9 5 8 14;4 10 6 6 15;4 1 18 4 19;3 19 15 3;4 18 13 14 19
3;1 6;4 13 17 16 1;3 4 15 7
1;2 6 10
5;4 15 7 4 15;3 7 1 12;1 14;2 16 15;3 17 15 19
2;4 7 1 20 5;2 4 3
4;5 12 9 17 8 8;2 1 7;5 14 8 9 3 16;1 1
3;1 12;1 4;5 15 14 8 16 12
1;5 7 18 17 5 12
4;2 2 15;5 5 1 3 6 17;5 8 11 17 6 3;2 9 18
2;3 1 9 10;4 2 16 7 7
4;2 6 19;4 8 8 12 17;5 12 1 20 12 1;5 14 17 14 3 11
3;4 10 12 1 4;4 12 3 3 13;1 9
2;3 17 18 11;5 20 3 6 16 7
2;2 13 9;1 20
4;2 2 9;1 2;1 9;1 5
5;2 10 8;4 13 2 7 5;4 9 19 14 5;3 12 12 9;5 17 7 2 11 19
4;4 17 10 8 12;1 3;5 2 16 18 11 15;1 15
5;4 4 6 5 10;1 3;2 11 1;2 12 18;1 14
3;5 12 6 19 8 18;5 1 14 17 1 14;4 7 14 19 19
2;1 13;5 1 10 4 14 10
4;2 12 2;4 12 12 5 7;5 5 1 15 7 10;2 16 6
4;1 14;2 4 2;2 7 6;3 3 15 1
4;2 9 2;5 15 5 3 15 8;5 6 11 14 14 9;4 7 6 20 14
4;3 11 19 11;4 4 16 4 1;5 17 11 8 19 5;5 8 8 1 7 12
4;2 18 17;3 4 19 19;2 9 1;2 3 8
3;2 10 2;2 2 4;5 14 12 13 9 4
1;4 2 13 4 13
4;4 8 5 9 7;1 19;5 15 10 11 17 20;4 12 5 13 6
2;1 3;2 2 8
2;3 15 18 14;1 19
1;1 1
2;4 9 6 19 4;4 11 18 16 18
5;2 17 16;3 20 14 8;2 16 13;3 1 16 20;2 8 12
4;4 13 1 2 10;4 8 9 4 7;1 20;3 14 1 13
4;5 19 12 5 4 18;2 5 8;4 9 6 16 15;1 19
5;1 3;2 2 17;3 17 10 19;3 17 12 19;3 5 2 1
3;5 2 13 12 6 10;3 19 4 16;2 3 13
1;3 8 10 12
4;3 11 19 10;5 20 11 10 15 17;1 5;1 6
2;1 5;3 12 16 1
3;4 11 8 17 15;4 12 6 16 6;5 19 13 12 14 3
2;5 1 19 16 20 12;4 7 18 12 16
5;4 5 15 12 7;3 11 20 8;5 9 20 5 8 1;1 3;1 2
1;4 7 18 11 18
4;5 10 1 18 18 6;4 10 18 18 18;3 13 4 2;4 11 3 18 9
5;3 10 9 19;2 14 9;4 8 3 3 17;3 6 3 18;2 5 7
2;2 20 4;5 13 11 16 3 2
1;1 1
1;4 17 5 3 14
4;1 14;5 3 12 17 12 15;1 17;1 17
1;2 1 7
4;2 11 20;5 11 5 16 5 14;4 10 2 13 20;1 4
1;4 9 18 17 8
2;2 5 4;3 13 8 8
3;4 3 11 13 1;3 16 12 12;2 19 8
4;4 14 10 9 12;5 3 2 7 4 11;1 5;4 5 7 10 1
5;1 10;1 6;3 4 12 19;5 7 5 15 14 18;1 19
4;1 6;3 3 19 2;5 19 7 10 3 16;3 17 18 1
4;3 5 3 10;2 6 14;1 11;3 15 20 4
1;4 11 17 10 7
1;5 8 8 12 11 11
1;4 13 11 15 11
1;5 3 11 6 4 2`

type testCase struct {
	n     int
	caves [][2][]int // each cave: need calc from array; stored as k, arr
}

func solveCase(tc testCase) int {
	type caveInfo struct {
		need int
		k    int
	}
	caves := make([]caveInfo, tc.n)
	for i := 0; i < tc.n; i++ {
		k := tc.caves[i][0][0]
		arr := tc.caves[i][1]
		maxNeed := 0
		for j, a := range arr {
			if val := a - j; val > maxNeed {
				maxNeed = val
			}
		}
		caves[i] = caveInfo{need: maxNeed + 1, k: k}
	}
	// insertion sort by need then k to avoid importing sort
	for i := 1; i < len(caves); i++ {
		j := i
		for j > 0 {
			swap := false
			if caves[j].need < caves[j-1].need {
				swap = true
			} else if caves[j].need == caves[j-1].need && caves[j].k < caves[j-1].k {
				swap = true
			}
			if !swap {
				break
			}
			caves[j], caves[j-1] = caves[j-1], caves[j]
			j--
		}
	}
	start := caves[0].need
	cur := start + caves[0].k
	for i := 1; i < len(caves); i++ {
		if cur < caves[i].need {
			start += caves[i].need - cur
			cur = caves[i].need
		}
		cur += caves[i].k
	}
	return start
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: not enough parts", idx+1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(parts)-1 != n {
			return nil, fmt.Errorf("line %d: expected %d caves got %d", idx+1, n, len(parts)-1)
		}
		tc := testCase{n: n, caves: make([][2][]int, n)}
		for i := 0; i < n; i++ {
			fields := strings.Fields(parts[i+1])
			if len(fields) == 0 {
				return nil, fmt.Errorf("line %d: empty cave data", idx+1)
			}
			k, err := strconv.Atoi(fields[0])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse k: %v", idx+1, err)
			}
			if len(fields) != 1+k {
				return nil, fmt.Errorf("line %d: expected %d numbers for cave got %d", idx+1, 1+k, len(fields))
			}
			arr := make([]int, k)
			for j := 0; j < k; j++ {
				v, err := strconv.Atoi(fields[1+j])
				if err != nil {
					return nil, fmt.Errorf("line %d: parse value: %v", idx+1, err)
				}
				arr[j] = v
			}
			tc.caves[i] = [2][]int{{k}, arr}
		}
		cases = append(cases, tc)
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
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, cave := range tc.caves {
			k := cave[0][0]
			sb.WriteString(strconv.Itoa(k))
			for _, v := range cave[1] {
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(v))
			}
			sb.WriteByte('\n')
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		vals := strings.Fields(got)
		if len(vals) != 1 {
			fmt.Printf("case %d: expected single integer output, got %q\n", i+1, got)
			os.Exit(1)
		}
		gotVal, err := strconv.Atoi(vals[0])
		if err != nil {
			fmt.Printf("case %d: non-integer output %q\n", i+1, vals[0])
			os.Exit(1)
		}
		if gotVal != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %d\n", i+1, expected, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
