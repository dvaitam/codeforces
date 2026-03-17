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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, n := range testcases {
		input := fmt.Sprintf("%d\n", n)
		expectedM := int64(n)/2 + 1
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		rawLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(rawLines) != n+1 {
			fmt.Printf("case %d failed: expected %d lines got %d\n", idx+1, n+1, len(rawLines))
			os.Exit(1)
		}
		// First line: m
		var gotM int64
		if _, err := fmt.Sscan(strings.TrimSpace(rawLines[0]), &gotM); err != nil || gotM != expectedM {
			fmt.Printf("case %d failed: expected m=%d got %s\n", idx+1, expectedM, strings.TrimSpace(rawLines[0]))
			os.Exit(1)
		}
		// Parse coordinates
		type pt struct{ r, c int64 }
		pts := make([]pt, n)
		for i := 0; i < n; i++ {
			var r, c int64
			if _, err := fmt.Sscan(strings.TrimSpace(rawLines[i+1]), &r, &c); err != nil {
				fmt.Printf("case %d failed: invalid line %d: %s\n", idx+1, i+2, rawLines[i+1])
				os.Exit(1)
			}
			if r < 1 || r > gotM || c < 1 || c > gotM {
				fmt.Printf("case %d failed: point %d (%d,%d) out of range [1,%d]\n", idx+1, i+1, r, c, gotM)
				os.Exit(1)
			}
			pts[i] = pt{r, c}
		}
		// Check constraint: for all i,j: |r_i-r_j|+|c_i-c_j| >= |i-j|
		ok := true
		for i := 0; i < n && ok; i++ {
			for j := i + 1; j < n && ok; j++ {
				dr := pts[i].r - pts[j].r
				if dr < 0 {
					dr = -dr
				}
				dc := pts[i].c - pts[j].c
				if dc < 0 {
					dc = -dc
				}
				dist := dr + dc
				diff := int64(j - i)
				if dist < diff {
					fmt.Printf("case %d failed: dist(%d,%d)=%d < %d\n", idx+1, i+1, j+1, dist, diff)
					ok = false
				}
			}
		}
		if !ok {
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
