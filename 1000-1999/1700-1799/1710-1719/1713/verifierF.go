package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `3 12 28 30
2 26 30
4 1 18 24 1
4 2 8 6 23
3 5 3 5
6 11 0 7 30 1 0
4 17 0 4 3
5 19 29 29 6 2
6 14 6 27 0 16 28
5 27 13 2 17 5
2 7 29
4 12 15 25 0
4 6 12 24 1
5 30 8 0 18 11
6 11 26 10 21 14 20
2 19 16
1 8
1 22
1 8
1 22
2 19 24
6 4 12 6 18 21 10
2 13 16
5 3 17 3 22 29
4 3 16 28 14
4 5 14 17 10
2 13 8
4 2 18 16 10
2 14 7
3 15 24 29
4 0 14 23 28
1 17
4 30 24 24 14
2 13 7
3 15 15 4
2 14 9
3 20 15 19
2 26 16
6 2 6 9 29 16 20
6 3 30 1 4 10 1
3 27 19 25
2 24 28
4 12 21 6 13
4 6 5 12 1
3 28 24 21
5 6 23 10 5 16
5 22 27 19 18 13
2 20 22
5 15 28 21 18 29
2 19 14
1 7
4 28 19 26 10
5 23 6 0 1 1
2 7 14
6 19 7 21 21 3 27
5 13 29 20 20 9
6 10 15 21 6 5 11
5 26 11 22 14 12
4 21 10 27 14
1 4
2 24 3
2 30 29
4 30 14 22 17
2 28 11
6 1 0 12 6 2 27
4 18 18 27 25
4 17 6 20 0
6 4 16 20 29 15 9
4 17 17 2 2
5 23 15 24 0 6
4 19 24 11 9
1 23
1 24
2 19 30
1 18
5 11 4 14 8 3
5 14 8 28 26 7
4 12 23 16 29
3 7 26 18
3 29 2 1
2 10 20
5 18 15 2 30 16
4 13 10 12 23
1 3
4 19 21 29 30
2 24 4
1 3
6 13 7 1 20 8 8
4 3 11 9 11
5 25 20 6 23 26
5 17 1 10 16 17
1 5
4 26 0 20 27
3 10 4 11
1 19
5 13 30 25 24 2
3 18 3 16
1 27
5 4 15 1 16 26
2 11 3`

type testCase struct {
	n int
	a []int
}

// solveF mirrors the current 1713F stub: always outputs -1.
func solveF(tc testCase) string {
	return "-1"
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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[i+1])
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')

		expected := solveF(tc)
		got, err := runExe(bin, input.String())
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", idx+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
