package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `2 3 0 740
4 4 1 1 1 1 413
5 3 0 229
5 5 2 3 2 1 1 3 -1 200
2 3 1 2 1 -1 689
4 5 1 2 2 -1 288
1 5 2 1 3 -1 1 4 -1 615
3 4 3 1 2 -1 2 1 1 1 4 -1 533
5 6 3 3 2 1 1 4 1 4 3 1 366
4 6 4 3 6 1 3 1 1 2 3 1 1 3 1 299
4 1 0 367
6 1 2 6 1 1 3 1 -1 991
2 6 3 1 3 1 2 3 1 2 4 1 341
5 1 2 1 1 1 3 1 -1 587
1 4 1 1 2 1 62
1 1 0 611
6 2 4 1 2 1 3 1 1 5 2 -1 6 1 -1 208
2 4 3 2 1 1 2 4 1 2 2 -1 194
1 1 0 261
2 5 1 1 4 -1 147
3 1 1 3 1 -1 987
6 6 5 6 1 -1 4 1 -1 2 5 1 3 3 -1 6 3 -1 542
2 6 5 2 3 -1 2 1 -1 1 1 -1 1 3 1 1 6 -1 584
4 6 3 4 2 1 2 2 1 3 1 -1 831
4 2 0 901
2 2 1 2 2 1 82
1 6 4 1 4 -1 1 1 -1 1 5 1 1 3 1 89
5 4 3 2 3 -1 2 4 -1 1 1 1 939
5 3 4 4 2 -1 1 3 1 3 2 -1 2 2 1 182
6 2 2 4 2 -1 5 1 1 7
6 3 0 560
1 4 3 1 2 -1 1 3 1 1 1 1 626
3 5 3 3 3 -1 3 5 -1 2 2 -1 791
5 3 2 2 2 1 5 1 1 294
3 2 2 2 1 1 2 2 -1 180
3 1 2 1 1 -1 3 1 1 178
6 4 5 6 2 -1 6 4 1 4 2 -1 3 2 -1 5 4 -1 520
4 3 2 4 2 1 3 1 -1 711
5 5 2 4 4 1 2 5 1 471
6 3 4 3 1 -1 4 1 -1 5 1 1 2 1 -1 137
1 2 1 1 1 1 163
1 2 0 974
3 1 0 279
1 6 2 1 6 1 1 2 1 575
3 6 0 743
5 2 4 4 1 1 3 2 1 2 2 -1 5 2 -1 338
2 4 0 435
1 3 2 1 2 1 1 1 1 870
5 1 1 2 1 -1 934
6 5 1 4 4 -1 869
2 5 4 1 1 -1 2 3 -1 1 2 -1 2 4 -1 828
1 2 1 1 1 -1 541
4 3 0 183
6 1 3 5 1 1 2 1 -1 6 1 1 521
5 1 0 866
2 5 2 1 1 1 2 2 1 370
2 6 5 2 3 -1 2 1 1 2 4 1 2 5 -1 1 6 -1 331
2 2 0 492
4 3 2 4 3 -1 3 3 1 264
3 5 0 32
4 3 2 2 3 -1 1 1 1 719
2 3 2 1 3 -1 2 1 -1 585
2 2 0 252
5 1 0 818
3 3 2 3 3 1 3 1 1 922
6 6 0 776
2 6 1 2 1 -1 648
3 6 4 2 2 1 1 1 1 1 6 1 3 3 -1 108
2 2 1 2 2 -1 372
4 6 3 1 1 1 2 3 1 3 2 -1 589
2 5 2 2 3 1 1 2 -1 191
5 2 0 29
1 2 1 1 2 1 975
3 6 3 1 2 -1 3 4 -1 2 1 -1 716
2 4 0 763
6 6 3 2 1 -1 3 5 1 1 4 1 999
6 5 0 187
5 2 4 3 1 1 4 2 1 5 1 1 2 1 -1 241
4 6 2 3 2 -1 3 3 -1 990
5 1 4 2 1 1 4 1 1 1 1 1 5 1 -1 897
1 6 0 805
5 2 4 2 2 1 3 2 -1 3 1 1 1 1 -1 544
5 3 0 207
4 6 0 821
4 1 2 2 1 1 1 1 1 980
3 5 4 3 4 1 1 5 -1 3 2 1 3 5 1 708
1 1 0 448
2 1 1 2 1 1 603
5 1 4 5 1 1 2 1 -1 3 1 1 1 1 -1 654
2 2 1 1 2 1 535
1 1 0 899
6 6 4 2 2 1 3 6 -1 4 2 -1 6 4 -1 521
1 3 0 100
2 1 1 1 1 1 446
6 6 0 573
2 3 1 2 3 1 2
3 3 0 779
4 5 1 4 3 -1 568
6 1 3 2 1 -1 4 1 1 1 1 1 131
2 4 3 2 4 -1 2 3 -1 1 3 -1 523`

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 != 0 {
			res = (res * a) % mod
		}
		a = (a * a) % mod
		e >>= 1
	}
	return res
}

func solve40E(n, m, k int, cells [][3]int, p int64) int64 {
	known := make([][]bool, n)
	for i := 0; i < n; i++ {
		known[i] = make([]bool, m)
	}
	bRow := make([]int, n)
	bCol := make([]int, m)
	for i := 0; i < n; i++ {
		bRow[i] = 1
	}
	for j := 0; j < m; j++ {
		bCol[j] = 1
	}
	knownR := make([]int, n)
	knownC := make([]int, m)
	for _, c := range cells {
		a := c[0] - 1
		b := c[1] - 1
		y := 0
		if c[2] < 0 {
			y = 1
		}
		known[a][b] = true
		knownR[a]++
		knownC[b]++
		bRow[a] ^= y
		bCol[b] ^= y
	}
	for i := 0; i < n; i++ {
		if knownR[i] == m && bRow[i] != 0 {
			return 0
		}
	}
	for j := 0; j < m; j++ {
		if knownC[j] == n && bCol[j] != 0 {
			return 0
		}
	}
	sumR, sumC := 0, 0
	for i := 0; i < n; i++ {
		sumR ^= bRow[i] & 1
	}
	for j := 0; j < m; j++ {
		sumC ^= bCol[j] & 1
	}
	if sumR != sumC {
		return 0
	}
	visitedR := make([]bool, n)
	visitedC := make([]bool, m)
	comps := 0
	queue := []int{}
	for i := 0; i < n; i++ {
		if !visitedR[i] && knownR[i] < m {
			comps++
			visitedR[i] = true
			queue = queue[:0]
			queue = append(queue, i)
			for qi := 0; qi < len(queue); qi++ {
				u := queue[qi]
				if u < n {
					r := u
					for cj := 0; cj < m; cj++ {
						if known[r][cj] || visitedC[cj] {
							continue
						}
						visitedC[cj] = true
						queue = append(queue, n+cj)
					}
				} else {
					cj := u - n
					for rr := 0; rr < n; rr++ {
						if known[rr][cj] || visitedR[rr] {
							continue
						}
						visitedR[rr] = true
						queue = append(queue, rr)
					}
				}
			}
		}
	}
	for j := 0; j < m; j++ {
		if knownC[j] == n {
			comps++
		}
	}
	U := int64(n)*int64(m) - int64(k)
	F := U - int64(n+m) + int64(comps)
	return modPow(2, F, p)
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesE), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			fmt.Fprintf(os.Stderr, "case %d: not enough fields\n", idx+1)
			os.Exit(1)
		}
		n := mustAtoi(fields[0])
		m := mustAtoi(fields[1])
		k := mustAtoi(fields[2])
		if len(fields) != 3+3*k+1 {
			fmt.Fprintf(os.Stderr, "case %d: expected %d fields, got %d\n", idx+1, 3+3*k+1, len(fields))
			os.Exit(1)
		}
		cells := make([][3]int, k)
		pos := 3
		for i := 0; i < k; i++ {
			a := mustAtoi(fields[pos])
			b := mustAtoi(fields[pos+1])
			c := mustAtoi(fields[pos+2])
			cells[i] = [3]int{a, b, c}
			pos += 3
		}
		pVal, err := strconv.ParseInt(fields[pos], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid modulus\n", idx+1)
			os.Exit(1)
		}
		want := strconv.FormatInt(solve40E(n, m, k, cells, pVal), 10)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, m)
		fmt.Fprintf(&input, "%d\n", k)
		for _, c := range cells {
			fmt.Fprintf(&input, "%d %d %d\n", c[0], c[1], c[2])
		}
		fmt.Fprintf(&input, "%d\n", pVal)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
