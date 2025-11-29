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
const testcasesRaw = `2 4 -4
3 -4 2
4 2 5
4 -2 -4
4 -5 1
4 4 -5
4 -1 -2
5 -4 0
1 -5 -5
5 -5 1
2 1 -5
5 -2 2
4 3 -2
3 -2 5
2 2 -1
1 1 3
1 -3 5
3 -4 0
5 1 3
2 -1 -1
5 2 3
4 4 -5
4 -2 1
4 5 -3
3 3 5
3 -4 2
5 -4 -3
5 1 0
4 -5 2
1 -1 4
5 4 1
2 -3 3
2 -5 -2
5 3 -2
4 3 0
5 0 2
3 5 3
5 -5 1
5 -3 3
5 -2 1
1 2 0
5 3 -2
5 1 2
3 1 0
1 3 3
5 4 0
4 4 -5
2 5 -3
5 4 -3
1 3 -1
1 5 -4
1 -5 2
1 -1 -2
3 -4 4
2 0 -1
1 -3 -3
3 3 -3
3 5 -1
4 0 2
4 -4 -5
3 1 0
4 -2 -1
1 -1 3
2 4 1
1 -2 -5
4 -3 -5
2 2 3
4 3 -2
5 2 -2
5 5 -5
4 5 4
3 5 5
4 -5 -1
2 -2 -5
3 -4 -4
3 -1 -3
4 4 -1
2 -5 3
1 4 -2
5 2 -3
5 3 -5
4 -2 0
1 -2 4
4 4 -2
4 -4 5
4 -1 3
4 -5 0
5 1 -1
1 -3 -2
3 4 -3
3 1 -2
3 5 -4
4 3 0
5 2 3
2 -4 -5
1 -3 -3
2 3 -2
3 0 4
5 -1 0
3 0 -4`

type testCase struct {
	a int64
	x int64
	y int64
}

func solveCase(tc testCase) int64 {
	a2 := 2 * tc.a
	x2 := 2 * tc.x
	y2 := 2 * tc.y
	if y2 <= 0 || y2%a2 == 0 {
		return -1
	}
	row := y2/a2 + 1
	var w int64
	if row == 1 || row%2 == 0 {
		w = 1
	} else {
		w = 2
	}
	var pos int64
	switch w {
	case 1:
		if x2 <= -tc.a || x2 >= tc.a {
			return -1
		}
		pos = 1
	case 2:
		if x2 <= -a2 || x2 >= a2 || x2 == 0 {
			return -1
		}
		if x2 < 0 {
			pos = 1
		} else {
			pos = 2
		}
	}
	L := row - 1
	var sumPrev int64
	switch {
	case L <= 0:
		sumPrev = 0
	case L == 1:
		sumPrev = 1
	default:
		countEven := L / 2
		countOdd := (L+1)/2 - 1
		sumPrev = 1 + countEven + countOdd*2
	}
	return sumPrev + pos
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields got %d", idx+1, len(fields))
		}
		a, err1 := strconv.ParseInt(fields[0], 10, 64)
		x, err2 := strconv.ParseInt(fields[1], 10, 64)
		y, err3 := strconv.ParseInt(fields[2], 10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("line %d: parse error: %v %v %v", idx+1, err1, err2, err3)
		}
		cases = append(cases, testCase{a: a, x: x, y: y})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
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
		input := fmt.Sprintf("%d %d %d\n", tc.a, tc.x, tc.y)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solveCase(tc)
		var parsed int64
		if _, err := fmt.Sscan(got, &parsed); err != nil {
			fmt.Printf("test %d: could not parse output %q: %v\n", i+1, got, err)
			os.Exit(1)
		}
		if parsed != expect {
			fmt.Printf("test %d failed\ninput: %sexpected: %d\ngot: %d\n", i+1, input, expect, parsed)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
