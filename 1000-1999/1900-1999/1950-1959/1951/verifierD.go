package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n int64
	k int64
}

var testcases = []testCase{
	{61, 152},
	{140, 34},
	{95, 155},
	{122, 161},
	{149, 17},
	{156, 4},
	{121, 67},
	{142, 60},
	{50, 184},
	{121, 139},
	{141, 122},
	{102, 164},
	{39, 60},
	{163, 39},
	{134, 100},
	{190, 4},
	{172, 199},
	{17, 41},
	{195, 152},
	{11, 78},
	{200, 8},
	{69, 122},
	{153, 185},
	{100, 183},
	{110, 102},
	{187, 148},
	{114, 35},
	{94, 25},
	{10, 35},
	{127, 56},
	{67, 173},
	{112, 200},
	{161, 78},
	{108, 130},
	{99, 147},
	{90, 137},
	{150, 105},
	{150, 60},
	{87, 175},
	{8, 72},
	{156, 172},
	{179, 42},
	{179, 84},
	{139, 147},
	{146, 27},
	{183, 168},
	{55, 163},
	{147, 69},
	{73, 32},
	{17, 124},
	{164, 124},
	{23, 89},
	{18, 106},
	{39, 6},
	{76, 110},
	{197, 107},
	{31, 12},
	{155, 158},
	{195, 12},
	{97, 184},
	{151, 85},
	{142, 72},
	{130, 61},
	{10, 80},
	{2, 20},
	{28, 154},
	{138, 9},
	{51, 105},
	{75, 157},
	{68, 40},
	{177, 11},
	{87, 81},
	{93, 36},
	{97, 97},
	{118, 134},
	{99, 165},
	{153, 175},
	{144, 27},
	{159, 130},
	{70, 111},
	{163, 185},
	{184, 61},
	{78, 112},
	{67, 134},
	{78, 141},
	{87, 3},
	{107, 149},
	{81, 6},
	{97, 158},
	{151, 162},
	{35, 16},
	{163, 161},
	{86, 120},
	{91, 174},
	{91, 156},
	{181, 72},
	{189, 126},
	{6, 151},
	{16, 174},
	{6, 95},
}

// solve1951D mirrors the reference solution logic embedded from 1951D.go.
func solve1951D(n, k int64) string {
	if k > n {
		return "NO"
	}
	if k == n {
		return "YES\n1\n1"
	}
	if k <= (n+1)/2 {
		return fmt.Sprintf("YES\n2\n%d %d", n-k+1, 1)
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	for idx, tc := range testcases {
		exp := solve1951D(tc.n, tc.k)

		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx+1, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != exp {
			fmt.Printf("Test %d failed: expected %q got %q\n", idx+1, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
