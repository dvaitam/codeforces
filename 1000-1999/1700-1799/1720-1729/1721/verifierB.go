package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
4 3 3 1 1
9 2 7 2 0
9 6 4 5 3
7 2 6 1 3
8 2 4 2 3
7 5 6 2 4
2 8 1 3 1
7 10 4 9 6
6 6 5 4 6
2 9 1 7 2
7 10 6 6 2
9 10 2 3 0
9 2 5 2 2
4 10 2 1 0
8 10 3 9 6
8 2 8 2 3
10 8 8 6 11
2 10 2 8 0
5 4 5 2 1
10 6 1 6 2
3 2 2 1 0
4 7 3 1 1
2 6 2 3 1
6 10 2 10 0
5 2 4 1 0
4 9 4 9 3
10 9 4 9 0
8 7 7 1 4
5 2 3 1 1
6 6 6 2 4
4 2 4 1 0
8 5 6 1 3
8 5 8 1 6
6 10 4 1 0
4 5 3 5 2
7 8 2 5 3
8 10 6 9 7
3 2 3 1 0
6 5 5 4 2
10 3 6 1 1
8 4 3 3 1
8 3 4 3 1
6 7 3 5 1
9 6 2 1 0
2 3 2 1 0
5 5 5 4 2
3 9 1 4 2
3 8 2 5 4
9 7 2 2 1
2 2 2 2 0
3 7 3 4 1
6 5 5 5 7
7 6 2 5 3
6 5 2 3 1
6 3 4 1 0
7 4 3 3 3
7 3 5 3 1
5 5 1 2 0
2 6 2 4 2
3 10 2 2 1
10 6 3 2 1
7 10 7 4 5
6 8 5 3 0
5 6 1 6 4
10 9 9 8 0
8 7 3 3 0
8 2 1 2 0
4 10 3 9 7
5 5 3 4 3
8 7 5 6 7
2 3 2 1 0
3 10 3 7 2
4 6 4 2 0
9 8 6 7 5
10 2 9 1 4
3 6 3 1 1
9 5 7 4 6
4 7 4 2 3
3 8 3 7 1
6 6 2 4 0
5 10 4 10 0
2 5 2 2 1
5 8 2 5 1
2 3 1 3 1
3 10 2 10 4
8 10 6 9 10
2 3 2 3 2
6 10 4 6 3
8 8 4 1 0
8 8 7 6 2
9 5 5 1 2
8 6 3 1 0
7 6 7 6 6
10 6 3 4 1
6 10 1 10 2
7 3 6 2 0
4 10 2 2 1
3 10 3 6 0
4 6 3 3 3
8 10 7 3 7`

type testCase struct {
	n  int
	m  int
	sx int
	sy int
	d  int
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
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

func solveCase(n, m, sx, sy, d int) string {
	path1 := (sx-1 > d) && (m-sy > d)
	path2 := (sy-1 > d) && (n-sx > d)
	if path1 || path2 {
		return fmt.Sprintf("%d", n+m-2)
	}
	return "-1"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	if len(lines)-1 != t {
		return nil, fmt.Errorf("expected %d cases got %d", t, len(lines)-1)
	}
	var cases []testCase
	for idx, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) != 5 {
			return nil, fmt.Errorf("line %d: expected 5 numbers got %d", idx+2, len(fields))
		}
		nums := make([]int, 5)
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse field %d: %v", idx+2, i+1, err)
			}
			nums[i] = v
		}
		cases = append(cases, testCase{n: nums[0], m: nums[1], sx: nums[2], sy: nums[3], d: nums[4]})
	}
	return cases, nil
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

	for idx, tc := range cases {
		input := fmt.Sprintf("1\n%d %d %d %d %d\n", tc.n, tc.m, tc.sx, tc.sy, tc.d)
		exp := solveCase(tc.n, tc.m, tc.sx, tc.sy, tc.d)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
