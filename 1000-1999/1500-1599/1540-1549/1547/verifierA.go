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
const testcasesRaw = `100
7 7
1 5
9 8
7 5
8 6
10 4
9 3
5 3
2 10
5 9
10 3
5 2
7 6
10 4
9 8
8 9
5 1
9 1
2 7
1 10
8 6
4 6
2 4
10 4
4 3
9 8
2 2
6 9
8 2
5 9
5 2
9 6
9 4
10 9
10 5
8 2
10 7
6 10
4 5
3 4
3 1
10 5
8 2
2 3
3 1
2 9
7 9
5 9
4 4
10 7
10 5
8 8
6 2
6 10
2 8
10 6
4 4
1 5
2 4
6 3
6 7
1 2
3 4
1 10
9 10
2 1
2 4
10 10
2 7
2 6
2 1
10 1
4 3
2 8
4 1
1 9
7 10
2 5
2 4
2 5
6 7
3 1
9 8
1 10
2 7
4 5
6 8
10 3
4 1
3 3
6 9
5 2
10 8
3 1
8 7
10 9
5 6
7 5
3 9
1 8
2 6
1 9
5 3
4 8
6 10
5 6
10 10
3 5
7 7
2 1
10 4
6 3
4 4
8 7
1 3
8 2
5 3
1 8
6 5
8 1
7 4
9 2
3 1
7 7
6 1
4 1
1 9
10 2
4 2
10 4
5 5
3 2
8 7
2 1
5 8
2 5
3 9
6 2
3 5
1 1
1 4
5 9
6 6
10 1
10 8
8 7
6 9
3 4
7 10
5 1
3 3
5 6
6 6
2 6
10 1
1 5
3 3
10 5
6 7
9 3
5 2
8 4
1 5
3 9
2 5
7 6
5 7
2 2
9 8
8 6
6 2
8 2
8 7
1 5
6 3
3 10
7 2
2 2
4 4
1 7
1 2
7 9
9 5
8 8
10 4
7 2
6 4
5 10
3 7
4 6
2 2
1 9
8 4
2 8
7 5
4 1
4 10
3 2
4 8
7 6
9 3
2 10
8 3
10 7
7 9
8 6
8 8
4 9
10 4
1 6
6 6
1 9
3 5
10 3
7 10
5 8
2 2
9 1
2 4
3 1
5 1
8 6
3 3
8 6
9 7
9 9
1 10
2 9
10 2
7 4
5 9
10 7
8 7
10 10
4 1
1 3
5 9
10 5
6 2
8 5
5 7
7 7
1 3
3 4
5 6
1 1
8 7
3 8
10 2
3 6
7 1
10 8
7 8
1 2
8 3
1 1
10 10
3 6
2 9
6 4
7 8
2 1
10 8
10 6
2 10
5 3
7 5
2 9
4 1
7 8
6 4
8 6
2 1
1 8
5 1
9 10
10 4
4 2
9 9
7 9
5 2
3 7
10 7
2 2
7 2
2 7
3 1
8 7
7 1
8 6
5 2
6 2
2 6
1 6
6 3
1 4
6 2
10 3
4 1
4 2
`

type testCase struct {
	xA, yA int
	xB, yB int
	xF, yF int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(tc testCase) int {
	dist := abs(tc.xA-tc.xB) + abs(tc.yA-tc.yB)
	if tc.xA == tc.xB && tc.xA == tc.xF {
		minY, maxY := tc.yA, tc.yB
		if minY > maxY {
			minY, maxY = maxY, minY
		}
		if tc.yF > minY && tc.yF < maxY {
			dist += 2
		}
	} else if tc.yA == tc.yB && tc.yA == tc.yF {
		minX, maxX := tc.xA, tc.xB
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		if tc.xF > minX && tc.xF < maxX {
			dist += 2
		}
	}
	return dist
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	expectedLen := 1 + t*6
	if len(fields) != expectedLen {
		return nil, fmt.Errorf("expected %d fields, got %d", expectedLen, len(fields))
	}
	res := make([]testCase, t)
	for i := 0; i < t; i++ {
		offset := 1 + i*6
		vals := [6]int{}
		for j := 0; j < 6; j++ {
			v, err := strconv.Atoi(fields[offset+j])
			if err != nil {
				return nil, fmt.Errorf("parse case %d field %d: %v", i+1, j+1, err)
			}
			vals[j] = v
		}
		res[i] = testCase{
			xA: vals[0], yA: vals[1],
			xB: vals[2], yB: vals[3],
			xF: vals[4], yF: vals[5],
		}
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
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d\n", tc.xA, tc.yA, tc.xB, tc.yB, tc.xF, tc.yF))

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
