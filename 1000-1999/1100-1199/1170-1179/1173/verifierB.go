package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution1173BSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var n int64
	if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
		return
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	// number of points to output
	cnt := n/2 + 1
	fmt.Fprintln(w, cnt)
	r, c := int64(1), int64(1)
	for i := int64(0); i < n; i++ {
		fmt.Fprintf(w, "%d %d\n", r, c)
		if i%2 == 1 {
			r++
		} else {
			c++
		}
	}
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1173BSource

var testcases = []int{
	138,
	583,
	868,
	822,
	783,
	65,
	262,
	121,
	508,
	780,
	461,
	484,
	668,
	389,
	808,
	215,
	97,
	500,
	30,
	915,
	856,
	400,
	444,
	623,
	781,
	786,
	3,
	713,
	457,
	273,
	739,
	822,
	235,
	606,
	968,
	105,
	924,
	326,
	32,
	23,
	27,
	666,
	555,
	10,
	962,
	903,
	391,
	703,
	222,
	993,
	433,
	744,
	30,
	541,
	228,
	783,
	449,
	962,
	508,
	567,
	239,
	354,
	237,
	694,
	225,
	780,
	471,
	976,
	297,
	949,
	23,
	427,
	858,
	939,
	570,
	945,
	658,
	103,
	191,
	645,
	742,
	881,
	304,
	124,
	761,
	341,
	918,
	739,
	997,
	729,
	513,
	959,
	991,
	433,
	520,
	850,
	933,
	687,
	195,
	311,
}

func solveCase(n int64) []string {
	cnt := n/2 + 1
	lines := make([]string, 0, cnt+1)
	lines = append(lines, fmt.Sprintf("%d", cnt))
	r, c := int64(1), int64(1)
	for i := int64(0); i < n; i++ {
		lines = append(lines, fmt.Sprintf("%d %d", r, c))
		if i%2 == 1 {
			r++
		} else {
			c++
		}
	}
	return lines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, n := range testcases {
		input := fmt.Sprintf("%d\n", n)
		wantLines := solveCase(int64(n))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		rawLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(rawLines) != len(wantLines) {
			fmt.Printf("case %d failed: expected %d lines got %d\n", idx+1, len(wantLines), len(rawLines))
			os.Exit(1)
		}
		for i := range wantLines {
			if strings.TrimSpace(rawLines[i]) != wantLines[i] {
				fmt.Printf("case %d failed at line %d: expected %s got %s\n", idx+1, i+1, wantLines[i], strings.TrimSpace(rawLines[i]))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
