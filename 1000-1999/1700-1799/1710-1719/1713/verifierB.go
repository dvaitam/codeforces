package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `5 5 15 13 6 14
4 6 8 15 11
5 5 12 15 3 16
2 10 1
6 15 20 15 1 7 10
1 10
5 20 5 14 16 3
6 16 8 18 13 9 1
1 9
6 2 1 9 13 17 19
6 13 15 4 9 12 10
6 7 20 3 2 3 9
3 18 11 4
5 8 6 3 14 10
3 17 5 19
5 7 18 4 14 18
4 9 10 15 12
5 5 6 4 4 13
4 19 15 5 18
6 10 12 16 14 7 16
4 17 11 16 2
4 10 5 16 2
5 7 1 12 16 13
1 17
1 3
6 13 1 12 2 4 20
1 9
6 10 8 5 19 10 7
1 14
4 11 13 6 11
4 14 5 15 5
5 11 5 7 6 15
3 13 14 16
4 8 7 15 7
5 2 13 2 8 3
2 12 2
6 6 8 20 10 20 3
6 17 10 12 14 15 2
6 17 18 14 19 15 16
3 16 7 11
3 2 2 2
2 12 1
3 1 5 3
4 8 20 13 18
2 15 7
3 20 4 20
1 11
3 18 15 11
3 1 17 2
2 12 3
2 17 12
2 7 9
6 10 10 17 13 9 16
3 8 2 10
5 3 1 15 16 15
1 14
4 15 15 4 3
1 8
1 5
4 7 15 20 3
4 18 13 2 6
2 16 8
2 9 12
3 14 4 18
3 20 18 7
6 10 15 17 20 15 18
6 9 9 8 1 4 20
6 4 6 14 8 7 10
6 1 18 17 14 2 4
4 9 4 19 12
2 18 10
2 8 3
5 10 11 8 12 16
3 19 6 5
1 18
5 11 12 19 1 5
4 5 6 17 3
2 7 16
5 7 8 5 8 13
3 20 19 5
6 16 4 20 1 17 20
3 16 15 10
1 8
5 6 16 16 18 11
6 3 9 5 20 13 7
3 10 13 2
2 2 11
6 8 11 15 8 9 12
6 6 10 1 12 19 18
1 5
3 1 16 2
1 8
1 1
2 11 3
1 12
6 14 5 7 15 14 5
3 10 6 11
6 14 13 1 14 9 18
5 15 2 19 4 14
4 6 1 17 5`

type testCase struct {
	n int
	a []int64
}

// solveB mirrors 1713B.go logic.
func solveB(a []int64) string {
	if len(a) == 0 {
		return "YES"
	}
	var ops int64 = a[0]
	var maxVal int64
	for _, v := range a {
		if v > maxVal {
			maxVal = v
		}
	}
	for i := 1; i < len(a); i++ {
		if a[i] > a[i-1] {
			ops += a[i] - a[i-1]
		}
	}
	if ops == maxVal {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(parts) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n, len(parts)-1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(parts[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, a: arr})
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

func runExe(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')

		expected := solveB(tc.a)
		got, err := runExe(bin, input.String())
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", idx+1, err, input.String())
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
