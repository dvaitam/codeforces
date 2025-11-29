package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1101ASource = `package main

import "fmt"

func main() {
	var q int
	if _, err := fmt.Scan(&q); err != nil {
		return
	}
	for i := 0; i < q; i++ {
		var l, r, d int64
		fmt.Scan(&l, &r, &d)
		if d < l {
			fmt.Println(d)
		} else {
			fmt.Println((r/d + 1) * d)
		}
	}
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1101ASource

type testCase struct {
	l int64
	r int64
	d int64
}

var testcases = []testCase{
	{l: 50, r: 98, d: 54},
	{l: 6, r: 22, d: 66},
	{l: 63, r: 88, d: 39},
	{l: 62, r: 84, d: 75},
	{l: 28, r: 60, d: 18},
	{l: 37, r: 45, d: 97},
	{l: 13, r: 52, d: 33},
	{l: 69, r: 114, d: 78},
	{l: 19, r: 38, d: 13},
	{l: 94, r: 98, d: 88},
	{l: 43, r: 73, d: 72},
	{l: 13, r: 35, d: 56},
	{l: 41, r: 80, d: 82},
	{l: 27, r: 62, d: 62},
	{l: 57, r: 90, d: 34},
	{l: 8, r: 43, d: 2},
	{l: 12, r: 58, d: 52},
	{l: 91, r: 141, d: 86},
	{l: 81, r: 81, d: 79},
	{l: 64, r: 85, d: 32},
	{l: 94, r: 114, d: 91},
	{l: 9, r: 21, d: 73},
	{l: 29, r: 44, d: 19},
	{l: 70, r: 98, d: 12},
	{l: 11, r: 31, d: 66},
	{l: 63, r: 69, d: 39},
	{l: 71, r: 89, d: 91},
	{l: 16, r: 51, d: 43},
	{l: 70, r: 83, d: 78},
	{l: 71, r: 108, d: 37},
	{l: 57, r: 62, d: 77},
	{l: 50, r: 70, d: 74},
	{l: 31, r: 49, d: 24},
	{l: 25, r: 36, d: 5},
	{l: 79, r: 121, d: 34},
	{l: 61, r: 65, d: 12},
	{l: 87, r: 135, d: 17},
	{l: 20, r: 22, d: 11},
	{l: 90, r: 124, d: 88},
	{l: 51, r: 96, d: 68},
	{l: 36, r: 69, d: 31},
	{l: 28, r: 71, d: 76},
	{l: 54, r: 91, d: 36},
	{l: 58, r: 89, d: 85},
	{l: 83, r: 127, d: 46},
	{l: 11, r: 31, d: 79},
	{l: 15, r: 46, d: 76},
	{l: 81, r: 102, d: 25},
	{l: 32, r: 33, d: 94},
	{l: 35, r: 42, d: 91},
	{l: 29, r: 52, d: 22},
	{l: 43, r: 70, d: 8},
	{l: 13, r: 63, d: 19},
	{l: 90, r: 104, d: 6},
	{l: 74, r: 114, d: 69},
	{l: 78, r: 121, d: 10},
	{l: 4, r: 11, d: 82},
	{l: 25, r: 63, d: 74},
	{l: 16, r: 41, d: 12},
	{l: 48, r: 55, d: 5},
	{l: 78, r: 79, d: 25},
	{l: 24, r: 69, d: 16},
	{l: 62, r: 75, d: 94},
	{l: 8, r: 51, d: 3},
	{l: 70, r: 97, d: 80},
	{l: 13, r: 29, d: 9},
	{l: 29, r: 33, d: 83},
	{l: 39, r: 61, d: 56},
	{l: 24, r: 27, d: 65},
	{l: 60, r: 62, d: 77},
	{l: 13, r: 57, d: 51},
	{l: 26, r: 42, d: 46},
	{l: 94, r: 124, d: 73},
	{l: 22, r: 66, d: 87},
	{l: 27, r: 76, d: 8},
	{l: 87, r: 97, d: 21},
	{l: 44, r: 77, d: 33},
	{l: 16, r: 54, d: 57},
	{l: 86, r: 97, d: 2},
	{l: 61, r: 104, d: 53},
	{l: 73, r: 105, d: 40},
	{l: 84, r: 106, d: 50},
	{l: 85, r: 101, d: 20},
	{l: 72, r: 116, d: 2},
	{l: 59, r: 106, d: 11},
	{l: 43, r: 90, d: 6},
	{l: 70, r: 87, d: 18},
	{l: 31, r: 79, d: 62},
	{l: 46, r: 85, d: 37},
	{l: 87, r: 109, d: 76},
	{l: 82, r: 121, d: 17},
	{l: 92, r: 111, d: 50},
	{l: 96, r: 122, d: 84},
	{l: 11, r: 11, d: 77},
	{l: 25, r: 69, d: 43},
	{l: 21, r: 36, d: 29},
	{l: 82, r: 110, d: 49},
	{l: 91, r: 134, d: 73},
	{l: 54, r: 56, d: 52},
	{l: 90, r: 126, d: 54},
}

func solveCase(l, r, d int64) int64 {
	if d < l {
		return d
	}
	return (r/d + 1) * d
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
	for idx, tc := range testcases {
		input := fmt.Sprintf("1\n%d %d %d\n", tc.l, tc.r, tc.d)
		expected := solveCase(tc.l, tc.r, tc.d)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d failed: execution error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
		if err != nil {
			fmt.Printf("test %d failed: bad output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
