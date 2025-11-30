package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

const maxZ = 110
const eps = 1e-6

// rawTestcases embeds the contents of testcasesE.txt.
var rawTestcases = []string{
	"2 2 0 3 2 1 1 1 1 1 2 1 1",
	"2 4 4 2 2 1 2 2 1 1",
	"3 2 1 1 2 2 1 1 3 3 1",
	"4 4 1 1 1 2 3 3 1 1 3 3",
	"2 3 3 3 2 2 2 1 2 1 2 1 1",
	"1 3 3 1 1 3 1 1 1 1 1 2",
	"4 2 1 2 3 3 3 3 1 1 2 1 3 1 1",
	"3 1 2 3 1 1 3 1",
	"1 2 3 1 1 1 1 1 2 1 1 1",
	"2 3 2 1 2 1 2",
	"3 4 0 2 1 2 2 1",
	"3 2 4 0 2 1 1 2 1 1 1",
	"1 5 1 1 1 5",
	"4 4 1 2 0 1 3 2 2",
	"4 1 1 3 3 2 1 4 1 4 4 2",
	"2 3 1 1 1 2 2",
	"2 4 1 1 2 1 1",
	"3 0 2 4 1 3 3 4",
	"1 3 2 1 1 2 1 1 2",
	"4 5 2 3 4 1 3 4 2",
	"4 0 2 5 5 1 1 2 1",
	"3 5 0 1 3 3 3 1 2 2 1 1 1 1",
	"2 0 4 2 1 2 1 1 1 1",
	"2 5 4 2 2 1 4 1 1 1",
	"4 4 2 4 1 1 3 1 1",
	"3 2 0 0 3 2 1 1 2 1 1 2 1 1",
	"1 0 3 1 1 1 1 1 1 1 1 1",
	"1 4 1 1 1 2",
	"2 2 3 2 2 1 3 1 2 1",
	"1 4 1 1 1 4",
	"2 3 2 2 1 1 1 1 2 3",
	"4 4 2 4 3 3 4 2 2 3 3 3 2 1 4",
	"2 5 4 1 2 1 3",
	"3 5 4 0 1 2 2 2",
	"3 2 1 2 1 3 3 1",
	"1 2 1 1 1 1",
	"1 0 1 1 1 1",
	"4 4 1 3 5 1 3 2 2",
	"3 4 5 4 2 2 2 5 2 2 3",
	"4 4 0 2 1 3 4 4 1 4 2 1 1 3 4",
	"3 0 2 3 1 3 3 2",
	"1 1 1 1 1 1",
	"1 1 2 1 1 1 1 1 1",
	"1 2 1 1 1 2",
	"3 4 5 5 1 3 3 2",
	"4 1 5 5 0 3 3 2 1 4 2 1 3 1 4",
	"4 1 1 2 4 2 3 3 2 2 1 1",
	"4 0 2 5 0 2 2 2 1 1 2 1",
	"2 4 0 2 1 2 4 2 1 2",
	"1 3 3 1 1 2 1 1 3 1 1 3",
	"4 5 2 0 1 2 1 3 1 4 1 1",
	"2 5 0 2 1 1 1 1 1 3",
	"3 0 4 4 1 1 3 1",
	"3 4 0 5 1 1 3 4",
	"2 1 2 1 2 2 2",
	"4 0 1 3 3 1 3 3 3",
	"2 5 5 2 2 1 2 1 2 4",
	"3 2 3 5 2 2 1 3 2 3 1",
	"1 0 2 1 1 1 1 1 1",
	"1 5 1 1 1 2",
	"1 5 2 1 1 4 1 1 1",
	"1 0 2 1 1 1 1 1 1",
	"2 2 0 3 2 2 1 1 1 1 1 1 1",
	"1 5 2 1 1 1 1 1 4",
	"4 4 2 2 3 3 4 1 1 2 3 1 2 2 1",
	"3 4 1 4 3 2 2 1 3 1 1 1 1 3",
	"2 0 0 1 1 2 1",
	"4 0 2 5 3 1 2 4 2",
	"3 2 0 3 3 1 1 2 3 3 3 2 1 1",
	"3 2 4 4 1 1 1 2",
	"1 1 3 1 1 1 1 1 1 1 1 1",
	"2 1 4 2 1 2 1 2 2 2",
	"3 2 4 2 3 1 2 1 3 2 1 2 1 1",
	"2 1 1 3 1 2 1 1 1 1 1 1 1",
	"3 3 2 1 3 2 3 2 3 1 1 2 1 1",
	"1 2 3 1 1 1 1 1 1 1 1 2",
	"2 1 1 3 1 1 1 1 2 1 1 2 1",
	"2 5 1 3 1 1 1 1 1 4 1 2 5",
	"3 5 2 1 3 3 3 1 3 3 1 1 1 2",
	"3 4 5 3 2 3 1 2 3 3 1",
	"4 1 4 3 1 2 1 3 1 4 2 1",
	"3 1 2 2 2 3 3 1 2 2 2",
	"3 5 1 2 1 1 3 3",
	"2 1 3 1 2 1 1",
	"1 0 3 1 1 1 1 1 1 1 1 1",
	"3 4 0 2 1 1 1 2",
	"4 0 3 1 3 3 1 3 1 3 2 2 1 3 1",
	"1 3 3 1 1 2 1 1 3 1 1 1",
	"4 3 4 1 0 1 3 2 1",
	"2 3 3 3 2 2 2 2 2 3 2 1 3",
	"3 3 4 0 3 2 3 4 2 3 1 3 2 2",
	"3 0 2 4 1 3 3 2",
	"2 3 1 3 2 2 1 2 2 1 1 2 3",
	"3 3 5 5 2 1 2 1 2 2 2",
	"4 4 1 3 3 3 4 1 2 4 2 1 2 4 1",
	"3 1 0 0 1 2 2 1",
	"3 3 4 1 2 1 1 3 1 1 2",
	"2 1 3 3 1 1 1 1 1 1 1 1 1",
	"4 2 2 2 4 2 2 1 2 3 1 2",
	"1 2 3 1 1 2 1 1 1 1 1 2",
}

// solve248E mirrors the reference solution from 248E.go using in-memory I/O.
func solve248E(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return "", err
	}
	a := make([]int, n+1)
	p := make([][]float64, n+1)
	for i := 1; i <= n; i++ {
		p[i] = make([]float64, maxZ)
	}
	var ans float64
	for i := 1; i <= n; i++ {
		if _, err := fmt.Fscan(reader, &a[i]); err != nil {
			return "", err
		}
		if a[i] < maxZ {
			p[i][a[i]] = 1.0
		}
		if a[i] == 0 {
			ans += 1.0
		}
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return "", err
	}
	tmp := make([]float64, maxZ)
	var sb strings.Builder
	for qi := 0; qi < q; qi++ {
		var u, v, k int
		if _, err := fmt.Fscan(reader, &u, &v, &k); err != nil {
			return "", err
		}
		for j := 0; j < k; j++ {
			ans -= p[u][0]
			au := a[u]
			denom := float64(au)
			for z := 0; z <= au && z+1 < maxZ; z++ {
				tmp[z] = (p[u][z]*(denom-float64(z)) + p[u][z+1]*float64(z+1)) / denom
			}
			for z := 0; z <= au && z < maxZ; z++ {
				p[u][z] = tmp[z]
			}
			a[u]--
			ans += p[u][0]
		}
		a[v] += k
		fmt.Fprintf(&sb, "%.10f\n", ans)
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

func sameFloats(exp, got string) bool {
	eLines := strings.Fields(exp)
	gLines := strings.Fields(got)
	if len(eLines) != len(gLines) {
		return false
	}
	for i := range eLines {
		var a, b float64
		if _, err := fmt.Sscan(eLines[i], &a); err != nil {
			return false
		}
		if _, err := fmt.Sscan(gLines[i], &b); err != nil {
			return false
		}
		if math.Abs(a-b) > eps {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range rawTestcases {
		input := tc + "\n"

		expected, err := solve248E(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if !sameFloats(expected, got) {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
