package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases from testcasesE.txt (one test per line).
const testcaseData = `
4 14 2 9 17 4 8 75 19 16 ? 5 6 ? 4 11 8? 20
2 10 4 1 16 AF3 12
4 11 20 7 18 4 9 408 1 3 ?0? 11 5 2?0 7 16 774 18
4 3 3 11 17 4 3 20 18 7 34? 3 11 59 8 6 1 6
1 20 3 9 1 5 16 1 3 16 C?8 17
2 7 19 4 11 77 12 3 0? 19 12 33 1 13 1B 8
3 6 11 14 1 3 1 2
5 18 20 3 1 4 2 11 161 12 15 1 20
1 7 2 13 7 7 13 A 1
5 14 20 4 9 3 2 3 22? 6 2 10? 4
4 7 9 12 16 5 4 101 6 7 41? 6 2 ?1 19 15 9BC 9 4 030 11
1 18 3 4 3 12 11 A5 19 16 49C 14
1 1 5 5 211 8 12 6B 19 15 1C 19 8 027 3 6 134 16
5 20 1 2 16 11 3 15 1D 7 10 120 13 12 50 7
1 1 5 11 3 4 11 344 6 3 ?0 1 6 ?? 4 15 4B 4
2 9 1 1 2 ? 9
5 11 12 19 2 20 4 13 7A6 12 15 56C 19 6 1 5 6 2? 12
1 11 5 2 1 6 4 223 18 4 03 8 16 9 6 15 29C 11
3 14 4 4 5 16 FA 11 3 0? 14 2 1? 5 16 C 3 3 1 8
1 13 1 3 2? 16
5 7 14 3 12 8 3 11 6 7 7 1 1 16 E63 16
4 9 7 2 7 5 4 1 15 8 ?2 4 11 29 13 12 687 11 15 ?6 18
5 8 1 11 11 11 1 10 4 20
2 13 19 3 13 C71 3 10 1 8 4 2 1
4 11 6 5 15 3 10 88 2 11 A 17 14 1B6 7
3 18 20 14 4 15 70 1 13 4 17 11 51 16 15 9D 13
4 2 6 5 8 3 13 00 16 8 7 20 13 A 5
3 14 2 20 4 8 01 16 14 0 2 11 2A5 4 13 A53 13
4 4 2 20 15 5 12 A1 20 6 3 10 16 3?6 2 14 75 7 9 10 2
4 9 1 17 19 5 16 7 3 14 C8B 17 8 412 14 16 D23 14 3 ? 5
1 15 4 12 07 11 13 15 3 3 02 12 4 1 12
1 20 2 5 1 4 13 4 12
1 20 2 15 5 15 3 22 5
1 7 3 7 44 18 12 29 3 3 21? 5
2 8 11 5 5 1 10 7 02 20 2 0? 3 4 2? 14 13 9 14
3 12 3 8 4 12 A8 2 8 06 11 9 5 10 9 2 4
3 4 18 20 2 14 7?6 6 14 62 8
4 11 17 5 12 4 12 17? 7 6 ? 15 11 03 10 3 21? 16
1 16 2 10 4? 1 3 00 9
4 17 19 13 15 1 13 54 7
5 3 2 3 9 10 5 7 3 6 3 22 17 4 ?1? 4 8 ?64 10 9 ?2 6
1 4 4 8 72? 10 7 ?63 16 9 857 2 9 27 2
5 7 1 12 16 13 1 15 22C 1
3 2 4 20 1 6 525 8
2 19 10 2 3 ?2 13 4 33 5
4 5 17 11 5 2 16 E 12 14 6C 16
4 8 7 15 7 5 13 6 2 5 012 2 13 A23 20 6 05? 17 14 C? 12
4 15 2 17 18 4 11 74 16 5 20 2 2 1 1 6 010 14
2 20 13 5 16 E 7 7 115 11 10 54 1 10 3 12 3 2 7
2 9 10 3 10 47 12 15 719 18 3 ? 16
4 2 14 16 15 4 3 0 8 3 ? 7 9 168 13 2 0 16
2 5 9 3 7 14 20 10 4 15 10 78? 9
3 8 1 4 5 14 12B 14 5 2 1 13 86? 2 3 20 19 7 4 8
2 3 17 3 12 35 16 6 11? 1 10 559 1
2 13 5 2 10 2 7 14 9C 7
2 5 8 4 7 2?1 20 15 B 16 9 03 18 12 A 16
4 18 11 3 9 2 11 ?3 11 14 60 7
1 11 2 7 34 12 12 4 1
3 19 18 2 2 7 ? 2 2 0 1
2 11 3 1 7 623 15
4 5 12 10 6 3 13 660 14 16 E1 19 16 D 13
2 1 17 2 11 8?? 5 3 11 8
1 6 5 4 03? 4 11 7?2 20 11 4 11 14 B60 2 9 5 10
2 15 8 5 7 6 11 6 30 10 10 87 2 14 984 2 9 16 12
4 2 1 9 2 3 12 9B? 10 12 ? 17 10 64 7
1 19 3 14 9 18 12 22 11 15 014 12
3 10 10 11 4 14 96 6 2 ? 2 9 5 1 13 A? 9
5 7 3 18 14 9 2 10 1 6 11 8 18
5 13 14 9 10 10 1 8 44? 17
5 11 11 7 14 5 1 14 2CA 19
4 12 15 2 18 4 12 ??? 8 16 B 17 4 123 1 13 393 9
2 14 3 5 9 7 17 16 365 15 16 D 13 6 3? 12 11 14 1
4 1 9 7 13 4 8 607 12 15 48A 1 8 ?2 2 3 220 3
2 4 18 4 2 01 13 16 8 14 12 9 5 8 ?0 6
2 5 16 1 13 0?8 13
2 12 19 1 3 121 9
3 9 9 17 4 16 E 18 4 ? 6 12 0?5 3 5 341 15
5 6 11 5 16 18 1 10 8 11
1 3 1 8 575 13
5 12 4 5 11 1 2 13 0 11 16 61D 2
3 13 2 20 2 7 6 2 9 ?? 9
3 19 15 14 2 2 10 13 3 00 1
3 1 6 13 5 13 057 18 13 771 2 10 40 17 3 2 12 3 01 17
3 2 1 13 3 4 111 6 12 395 1 9 603 8
3 11 6 8 3 5 3 15 7 260 6 16 0C5 5
1 1 3 10 0 2 14 9 20 4 3 1
4 14 19 11 8 2 7 361 5 8 516 14
2 16 13 2 16 7F 13 16 288 17
3 18 1 20 5 14 ?3 9 2 11? 4 10 2 13 2 ?1 14 16 E 20
4 6 6 11 16 4 4 2?0 12 7 5 16 12 ?03 9 4 ?22 13
1 10 5 8 6 9 8 352 5 3 210 14 11 71 3 14 8B 18
2 6 5 2 4 0 17 4 2? 10
3 4 14 9 2 7 526 18 6 3 12
4 16 17 10 14 4 15 4 5 2 ??? 17 3 102 20 4 00 1
1 12 4 7 ? 12 11 4A7 8 16 2 17 16 0 6
5 14 20 7 15 13 3 2 010 15 11 6 3 12 653 17
4 2 16 20 4 3 15 ?C? 9 4 ? 12 4 ?1 15
1 3 4 8 6?1 9 9 1 10 4 01 2 15 460 2
4 18 20 16 4 4 10 51 1 5 3 10 6 0 16 7 511 10
5 14 8 12 13 5 2 6 5 16 12 46 20
`

var primes = []int{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
	31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
	73, 79, 83, 89, 97,
}

func solveInput(input string) (string, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return "", err
	}
	a := make([]uint64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	pc := len(primes)
	aMod := make([][]int, pc)
	for pi := 0; pi < pc; pi++ {
		p := primes[pi]
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = int(a[j] % uint64(p))
		}
		aMod[pi] = arr
	}
	var m int
	if _, err := fmt.Fscan(in, &m); err != nil {
		return "", err
	}
	type Q struct {
		dIdx int
		s    string
		c    uint64
	}
	qs := make([]Q, m)
	lengths := make([]map[int]int, 15)
	for i := range lengths {
		lengths[i] = make(map[int]int)
	}
	for i := 0; i < m; i++ {
		var d int
		var s string
		var c uint64
		fmt.Fscan(in, &d, &s, &c)
		idx := d - 2
		lengths[idx][len(s)] = 0
		qs[i] = Q{dIdx: idx, s: s, c: c}
	}

	lengthList := make([][]int, 15)
	for dIdx := 0; dIdx < 15; dIdx++ {
		mp := lengths[dIdx]
		if len(mp) == 0 {
			continue
		}
		list := make([]int, 0, len(mp))
		for L := range mp {
			list = append(list, L)
		}
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				if list[i] > list[j] {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
		lengthList[dIdx] = list
		for i, L := range list {
			lengths[dIdx][L] = i
		}
	}

	repDigits := make([][][]uint8, 15)
	repLens := make([][]int, 15)
	for dIdx := 0; dIdx < 15; dIdx++ {
		base := dIdx + 2
		repDigits[dIdx] = make([][]uint8, n)
		repLens[dIdx] = make([]int, n)
		for j := 0; j < n; j++ {
			x := j
			var digs []uint8
			if x == 0 {
				digs = []uint8{0}
			} else {
				for x > 0 {
					digs = append(digs, uint8(x%base))
					x /= base
				}
				for l, r := 0, len(digs)-1; l < r; l, r = l+1, r-1 {
					digs[l], digs[r] = digs[r], digs[l]
				}
			}
			repDigits[dIdx][j] = digs
			repLens[dIdx][j] = len(digs)
		}
	}

	indices := make([][][][]([]int), 15)
	for dIdx := 0; dIdx < 15; dIdx++ {
		list := lengthList[dIdx]
		if len(list) == 0 {
			continue
		}
		base := dIdx + 2
		Lcnt := len(list)
		indices[dIdx] = make([][][][]int, Lcnt)
		for li, L := range list {
			posList := make([][][]int, L)
			for pos := 0; pos < L; pos++ {
				posList[pos] = make([][]int, base)
			}
			indices[dIdx][li] = posList
			for j := 0; j < n; j++ {
				rl := repLens[dIdx][j]
				offset := L - rl
				for pos := 0; pos < L; pos++ {
					var dig int
					if pos < offset {
						dig = 0
					} else {
						dig = int(repDigits[dIdx][j][pos-offset])
					}
					indices[dIdx][li][pos][dig] = append(indices[dIdx][li][pos][dig], j)
				}
			}
		}
	}

	totalProd := make([][]int, len(primes))
	for pi := 0; pi < len(primes); pi++ {
		totalProd[pi] = make([]int, 15)
		p := primes[pi]
		for dIdx := 0; dIdx < 15; dIdx++ {
			prod := 1
			for j := 0; j < n; j++ {
				prod = (prod * aMod[pi][j]) % p
			}
			totalProd[pi][dIdx] = prod
		}
	}

	var out bytes.Buffer
	writer := bufio.NewWriter(&out)
	for _, q := range qs {
		dIdx := q.dIdx
		L := len(q.s)
		li := lengths[dIdx][L]
		type fv struct{ pos, val int }
		var fixed []fv
		for pos := 0; pos < L; pos++ {
			ch := q.s[pos]
			if ch != '?' {
				var v int
				if ch >= '0' && ch <= '9' {
					v = int(ch - '0')
				} else {
					v = int(ch - 'A' + 10)
				}
				fixed = append(fixed, fv{pos, v})
			}
		}
		ans := -1
		for pi := 0; pi < len(primes); pi++ {
			p := primes[pi]
			cmod := int(q.c % uint64(p))
			var prod int
			if len(fixed) == 0 {
				prod = totalProd[pi][dIdx]
			} else {
				best := fixed[0]
				bestList := indices[dIdx][li][best.pos][best.val]
				for _, f := range fixed[1:] {
					lst := indices[dIdx][li][f.pos][f.val]
					if len(lst) < len(bestList) {
						best = f
						bestList = lst
					}
				}
				prod = 1
				for _, j := range bestList {
					ok := true
					rl := repLens[dIdx][j]
					offset := L - rl
					for _, f := range fixed {
						var dig int
						if f.pos < offset {
							dig = 0
						} else {
							dig = int(repDigits[dIdx][j][f.pos-offset])
						}
						if dig != f.val {
							ok = false
							break
						}
					}
					if !ok {
						continue
					}
					prod = (prod * aMod[pi][j]) % p
				}
			}
			if (prod+cmod)%p == 0 {
				ans = p
				break
			}
		}
		if ans <= 100 && ans != -1 {
			fmt.Fprintln(writer, ans)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
	writer.Flush()
	return strings.TrimSpace(out.String()), nil
}

func parseLines() []string {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseLines()

	for idx, line := range tests {
		expect, err := solveInput(line)
		if err != nil {
			fmt.Printf("case %d failed to compute expected: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, line)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
