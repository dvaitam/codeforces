package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `49 97 53
5 33 65
62 51 100
38 61 45
74 27 64
17 36 17
96 12 79
32 68 90
77 18 39
12 93 9
87 42 60
71 12 45
55 40 78
81 26 70
61 56 66
33 7 70
1 11 92
51 90 100
85 80 0
78 63 42
31 93 41
90 8 24
72 28 30
18 69 57
11 10 40
65 62 13
38 70 37
90 15 70
42 69 26
77 70 75
36 56 11
76 49 40
73 30 37
23 24 23
4 78 84
33 60 8
11 86 96
16 19 4
10 89 69
87 50 90
67 35 66
30 27 86
75 53 74
35 57 63
84 82 89
45 10 41
78 14 62
75 80 42
24 31 2
93 34 14
90 28 47
21 42 54
7 12 100
18 89 28
5 73 81
68 77 87
9 3 15
81 24 77
73 15 50
11 47 14
4 77 2
24 23 91
15 61 26
93 7 86
2 69 54
79 12 33
8 28 9
82 38 44
55 23 7
64 59 5
76 12 89
50 25 33
45 93 60
72 21 89
86 26 98
7 100 86
20 20 43
67 32 15
76 56 85
22 1 60
87 52 72
65 39 83
45 49 84
32 19 71
88 1 58
94 10 42
94 5 69
35 17 30
97 61 45
78 36 86
45 75 81
79 16 91
39 49 95
53 83 10
0 76 24
89 42 20
30 28 81
57 48 90
86 72 53
4 51 89
`

// solve mirrors 1173A.go logic.
func solve(x, y, z int) string {
	dMin := x - y - z
	dMax := x + z - y
	switch {
	case dMin > 0:
		return "+"
	case dMax < 0:
		return "-"
	case dMin == 0 && dMax == 0:
		return "0"
	default:
		return "?"
	}
}

type testCase struct {
	x int
	y int
	z int
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	var tests []testCase
	vals := make([]int, 0)
	for scan.Scan() {
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	if len(vals)%3 != 0 {
		return nil, fmt.Errorf("invalid test file length")
	}
	for i := 0; i < len(vals); i += 3 {
		tests = append(tests, testCase{x: vals[i], y: vals[i+1], z: vals[i+2]})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d %d\n", tc.x, tc.y, tc.z)
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	expected := solve(tc.x, tc.y, tc.z)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
