package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `
1 4
3 32 29 18
1 87
5 12 76 55 5 4
1 28
2 65 78
1 72
2 92 84
5 54 29 58 76 36
1 98
2 90 55
3 36 20 28
3 14 12 49
1 46
3 78 34 6
4 69 16 49 11
5 38 81 80 47 74
2 91 9
1 85
2 99 38
1 30
1 49
3 59 82 47
2 48 46
2 86 35
1 78
2 69 94
2 21 60
4 35 82 89 72
2 88 42
1 30
1 41
4 35 9 28 73
3 28 84 64
4 83 59 19 34
2 32 96
5 69 34 96 75 55
5 52 47 29 18 66
4 12 97 7 15
2 81 21
4 77 9 50 49
5 60 68 33 71 2
1 88
5 97 35 99 83 44
1 38
4 21 59 1 93
3 65 98 23
5 14 81 39 82 65
5 26 20 48 98 21
5 100 68 1 77 42
4 3 15 47 40
2 8 31
5 11 11 94 63 9
5 99 17 17 85 61
5 22 34 68 78 55
2 70 97
2 92 40
4 86 84 48 57
5 58 16 32 29 9
3 3 76 71
2 76 29
1 10
1 30
1 5
3 10 66 31
3 86 63 28
5 17 93 74 74 61
2 61 53
2 13 13
4 46 55 53 60
1 87
1 8
4 94 44 14 32
2 25 69
4 18 55 24 36
4 32 10 57 71
1 7
5 2 12 97 31 22
4 63 62 28 52
1 22
4 1 50 34 59
3 55 90 94
5 85 92 63 20 25
3 28 8 75
5 8 96 41 8 7
5 62 65 68 21 8
5 11 24 9 77 9
2 52 16
5 32 75 77 6 80
1 54
5 73 67 41 34 27
3 31 34 51
2 86 83
3 59 41 97
1 2
4 80 73 13 10
5 28 65 34 17 45
1 32
3 37 21 57
5 91 39 79 84 68
`

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) int {
	if tc.n == 0 {
		return 0
	}
	minV := tc.arr[0]
	maxV := tc.arr[0]
	for _, v := range tc.arr {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	return maxV - minV
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n+1, len(fields))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
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

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for idx, v := range tc.arr {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(expected) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
