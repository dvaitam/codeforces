package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `3 72 97 8
4 15 63 97 57
5 83 48 100 26 12
5 3 49 55 77 97
2 89 57
4 92 29 75 13
4 3 2 3 83
6 1 48 87 27 54 92
2 67 28
5 63 70 29 44 29
3 97 58 37
2 53 71
2 23 80
4 15 95 42 92
6 54 64 85 24 38 36
6 63 64 50 75 4 61
3 95 51 53
3 46 70 89
4 11 56 84 65
2 99 20
6 50 47 62 93 3 60
2 39 90
6 75 74 50 82 21 21
6 29 1 98 25 69 70
3 51 65 44
6 45 58 34 84 70 77
2 49 100
6 16 66 99 71 26 54
2 61 46
6 70 25 64 52 62 45
6 43 76 17 100 22 35
6 88 87 57 67 66 96
2 41 1
6 82 26 21 48 65 34
3 53 91 60
4 22 80 72 87
2 30 9
2 64 26
6 26 16 60 97 67 76
6 25 56 11 80 7 86
2 7 19
3 37 69 100
5 53 85 95 10 63
2 10 12
6 44 61 2 88 27 93
4 85 65 27 91
4 87 51 32 3
3 44 88 3
2 65 35
6 17 47 27 62 33 50
4 81 43 39 8
6 25 6 68 6 77 3
2 7 68
6 30 4 100 37 43 5
6 85 28 18 16 43 24
6 99 15 99 8 18 3
5 2 37 73 87 49
3 40 41 65
2 60 29
4 52 47 100 77
4 32 75 35 43
2 58 46
6 94 48 95 26 82 45
4 58 8 75 40
2 40 75
2 63 65
6 7 81 84 38 95 4
6 18 70 64 81 32 62
4 24 4 9 4
4 91 32 60 2
4 37 90 81 13
2 30 7
5 43 79 47 89 60
5 8 64 83 59 95
4 68 74 58 78
2 60 79
5 79 11 25 23 42
6 19 34 59 58 80 48
3 72 11 6
4 17 80 15 53
2 99 24
3 71 87 63
3 42 23 26
5 22 64 58 60 18
4 27 36 85 68
2 71 7
3 56 5 41
3 54 98 4
6 73 2 17 76 94 41
2 11 60
2 16 53
2 7 63
5 74 14 68 64 68
3 24 12 52
5 95 57 16 22 74
6 42 92 67 59 36 10
3 39 98 53
3 83 38 46
5 87 20 33 47 79`

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) string {
	total := 0
	for _, v := range tc.arr {
		total ^= v
	}
	if total == 0 {
		return "YES"
	}
	prefix := 0
	cnt := 0
	for _, v := range tc.arr {
		prefix ^= v
		if prefix == total {
			cnt++
			prefix = 0
		}
	}
	if cnt >= 2 {
		return "YES"
	}
	return "NO"
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
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d: expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse value %d: %v", idx+1, i+1, err)
			}
			arr[i] = val
		}
		res = append(res, testCase{n: n, arr: arr})
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
