package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution1174DSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

// f computes the value for position i, adjusting for the avoid threshold
func f(i, avoid int) int {
	ans := i & -i
	if ans >= avoid {
		ans <<= 1
	}
	return ans
}

func main() {
	var n, x int
	if _, err := fmt.Scan(&n, &x); err != nil {
		return
	}
	total := 1 << n
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if x >= total {
		cnt := total - 1
		fmt.Fprintln(w, cnt)
		for i := 1; i < total; i++ {
			if i > 1 {
				w.WriteByte(' ')
			}
			fmt.Fprint(w, i&-i)
		}
		return
	}

	l := (1 << (n - 1)) - 1
	fmt.Fprintln(w, l)
	if l == 0 {
		return
	}
	avoid := x & -x
	for i := 1; i <= l; i++ {
		if i > 1 {
			w.WriteByte(' ')
		}
		fmt.Fprint(w, f(i, avoid))
	}
	fmt.Fprintln(w)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1174DSource

type testCase struct {
	n int
	x int
}

// f computes the value for position i, adjusting for the avoid threshold.
func f(i, avoid int) int {
	ans := i & -i
	if ans >= avoid {
		ans <<= 1
	}
	return ans
}

var testcases = []testCase{
	{n: 10, x: 260818},
	{n: 7, x: 199911},
	{n: 17, x: 135335},
	{n: 11, x: 101207},
	{n: 9, x: 54781},
	{n: 4, x: 148384},
	{n: 11, x: 210974},
	{n: 8, x: 153655},
	{n: 18, x: 178606},
	{n: 12, x: 42612},
	{n: 5, x: 86681},
	{n: 1, x: 153337},
	{n: 2, x: 148078},
	{n: 5, x: 249642},
	{n: 12, x: 95209},
	{n: 10, x: 164027},
	{n: 10, x: 84674},
	{n: 16, x: 205203},
	{n: 13, x: 157471},
	{n: 14, x: 258860},
	{n: 6, x: 266},
	{n: 5, x: 149056},
	{n: 2, x: 115839},
	{n: 5, x: 89482},
	{n: 1, x: 248672},
	{n: 16, x: 238527},
	{n: 9, x: 195896},
	{n: 7, x: 18415},
	{n: 18, x: 245388},
	{n: 14, x: 73062},
	{n: 6, x: 138814},
	{n: 6, x: 16591},
	{n: 6, x: 151962},
	{n: 4, x: 246849},
	{n: 17, x: 165709},
	{n: 18, x: 158610},
	{n: 13, x: 197855},
	{n: 14, x: 69668},
	{n: 10, x: 74702},
	{n: 1, x: 112383},
	{n: 9, x: 248055},
	{n: 9, x: 141000},
	{n: 17, x: 145375},
	{n: 11, x: 259547},
	{n: 11, x: 49805},
	{n: 14, x: 207491},
	{n: 5, x: 221444},
	{n: 1, x: 196622},
	{n: 17, x: 40720},
	{n: 13, x: 94619},
	{n: 15, x: 248912},
	{n: 2, x: 147220},
	{n: 14, x: 166898},
	{n: 8, x: 235881},
	{n: 1, x: 94955},
	{n: 17, x: 257308},
	{n: 6, x: 178326},
	{n: 7, x: 165389},
	{n: 12, x: 164578},
	{n: 16, x: 5018},
	{n: 8, x: 149831},
	{n: 8, x: 71952},
	{n: 6, x: 247712},
	{n: 14, x: 20428},
	{n: 15, x: 62059},
	{n: 15, x: 226431},
	{n: 17, x: 234609},
	{n: 4, x: 251136},
	{n: 7, x: 43120},
	{n: 15, x: 242106},
	{n: 3, x: 243267},
	{n: 14, x: 166396},
	{n: 13, x: 246518},
	{n: 9, x: 66241},
	{n: 14, x: 231156},
	{n: 12, x: 159827},
	{n: 11, x: 23758},
	{n: 10, x: 7900},
	{n: 16, x: 2876},
	{n: 9, x: 53160},
	{n: 13, x: 101352},
	{n: 14, x: 203279},
	{n: 13, x: 182741},
	{n: 2, x: 152749},
	{n: 15, x: 92952},
	{n: 5, x: 147959},
	{n: 9, x: 86014},
	{n: 1, x: 104358},
	{n: 16, x: 241414},
	{n: 17, x: 35479},
	{n: 2, x: 21361},
	{n: 12, x: 94640},
	{n: 1, x: 18186},
	{n: 7, x: 188033},
	{n: 4, x: 176003},
	{n: 18, x: 123497},
	{n: 2, x: 82453},
	{n: 1, x: 82524},
	{n: 13, x: 237311},
	{n: 5, x: 258032},
}

func solveCase(tc testCase) []string {
	n := tc.n
	x := tc.x
	total := 1 << n
	lines := []string{}
	if x >= total {
		cnt := total - 1
		lines = append(lines, fmt.Sprintf("%d", cnt))
		for i := 1; i < total; i++ {
			lines = append(lines, fmt.Sprintf("%d", i&-i))
		}
		return lines
	}
	l := (1 << (n - 1)) - 1
	lines = append(lines, fmt.Sprintf("%d", l))
	if l == 0 {
		return lines
	}
	avoid := x & -x
	for i := 1; i <= l; i++ {
		lines = append(lines, fmt.Sprintf("%d", f(i, avoid)))
	}
	return lines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.x)
		want := solveCase(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		gotLines := strings.Fields(string(out))
		if len(gotLines) != len(want) {
			fmt.Printf("case %d failed: expected %d tokens got %d\n", idx+1, len(want), len(gotLines))
			os.Exit(1)
		}
		for i := range want {
			if gotLines[i] != want[i] {
				fmt.Printf("case %d failed at token %d: expected %s got %s\n", idx+1, i+1, want[i], gotLines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
