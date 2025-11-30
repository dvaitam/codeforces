package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD = 1000000007

type group struct {
	A, B  []int
	cross bool
}

// Embedded testcases (previously in testcasesB.txt) so the verifier stays self contained.
const rawTestcasesB = `
4 9 12 8 9 14
3 10 20 4 5
1 2 12
4 10 14 4 17 9
3 8 13 6 15
4 9 10 15 20 1
4 5 4 15 1 19
2 1 11 15
4 1 16 2 2 12
1 6 1
1 10 0
3 5 7 4 18
3 4 3 13 14
3 7 5 10 13
4 3 14 4 16 10
2 4 5 14
3 7 13 15 12
2 4 14 6
5 1 12 1 7 20 2
2 6 1 20
2 4 19 9
5 2 16 9 11 13 14
1 9 20
5 7 18 14 15 8 15
2 6 8 1
1 1 5
3 1 9 20 0
2 2 13 7
5 7 17 7 14 6 10
5 2 19 2 10 10 17
4 6 8 0 16 1
2 6 2 6
5 6 6 6 8 9 9
5 7 8 15 11 7 1
3 9 2 0 14
4 8 1 13 15 14
4 2 2 2 7 3
2 7 6 14
5 2 13 17 12 1 5
2 8 7 4
3 6 10 13 3
5 5 19 17 6 9 14
5 10 14 17 20 8 8
2 1 3 19
1 3 13
2 4 9 0
5 9 13 1 3 12 20
3 2 18 11 7
5 5 7 7 2 16 9
3 4 11 20 15
3 10 5 4 0
5 9 10 11 18 20 0
2 7 4 5
5 2 4 6 15 18 6
2 3 7 12
3 10 18 4 20
4 2 19 0 16 19
3 8 14 9 0
2 9 20 5
4 8 17 10 2 8
2 10 12 6
3 5 12 1 6
1 6 7
3 8 7 8 11
2 5 0 11
5 9 1 20 4 11 0
4 1 0 7 1 0
2 6 2 1
3 7 4 6 14
4 3 11 9 5 20
3 7 12 0 13
3 9 17 14 1
5 2 13 12 5 0 16
2 10 16 4
1 6 7
2 4 0 5
5 3 2 13 19 3 19
4 3 19 19 1 8
3 7 0 20 1
4 2 11 9 4 14
2 9 11 5
4 6 8 15 12 0
3 9 9 17 15
1 9 18
5 5 1 14 12 3 12
3 8 1 0 8
1 5 18
3 4 16 16 10
4 5 6 3 18 10
2 10 17 11
2 3 10 0
5 1 18 4 11 11 9
3 6 15 12 19
4 3 0 4 18 1
4 3 10 0 15 8
5 4 2 17 13 8 5
5 3 2 20 5 18 3
5 9 19 12 13 8 9
3 1 13 8 8
5 9 17 10 10 6 13
`

func loadTestcases() ([]struct {
	n int
	m int64
	a []int
}, error) {
	lines := strings.Split(rawTestcasesB, "\n")
	var cases []struct {
		n int
		m int64
		a []int
	}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if len(fields) != n+2 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n+1, len(fields)-1)
		}
		m, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", idx+1, err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %w", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, struct {
			n int
			m int64
			a []int
		}{n: n, m: m, a: arr})
	}
	return cases, nil
}

// solveAllLower computes the sum of XOR lower bits (0..bit) for all pairs in groups.
func solveAllLower(groups []group, bit int) int64 {
	var sum int64
	for b := bit; b >= 0; b-- {
		bitw := int64(1) << uint(b)
		var next []group
		for _, g := range groups {
			if !g.cross {
				var z, o []int
				for _, v := range g.A {
					if v>>uint(b)&1 == 1 {
						o = append(o, v)
					} else {
						z = append(z, v)
					}
				}
				if len(z) > 0 && len(o) > 0 {
					sum = (sum + int64(len(z))*int64(len(o))%MOD*bitw) % MOD
					next = append(next, group{A: z, B: o, cross: true})
				}
				if len(z) > 1 {
					next = append(next, group{A: z, cross: false})
				}
				if len(o) > 1 {
					next = append(next, group{A: o, cross: false})
				}
			} else {
				var a0, a1, b0, b1 []int
				for _, v := range g.A {
					if v>>uint(b)&1 == 1 {
						a1 = append(a1, v)
					} else {
						a0 = append(a0, v)
					}
				}
				for _, v := range g.B {
					if v>>uint(b)&1 == 1 {
						b1 = append(b1, v)
					} else {
						b0 = append(b0, v)
					}
				}
				if len(a0) > 0 && len(b1) > 0 {
					sum = (sum + int64(len(a0))*int64(len(b1))%MOD*bitw) % MOD
					next = append(next, group{A: a0, B: b1, cross: true})
				}
				if len(a1) > 0 && len(b0) > 0 {
					sum = (sum + int64(len(a1))*int64(len(b0))%MOD*bitw) % MOD
					next = append(next, group{A: a1, B: b0, cross: true})
				}
				if len(a0) > 0 && len(b0) > 0 {
					next = append(next, group{A: a0, B: b0, cross: true})
				}
				if len(a1) > 0 && len(b1) > 0 {
					next = append(next, group{A: a1, B: b1, cross: true})
				}
			}
		}
		groups = next
		if len(groups) == 0 {
			break
		}
	}
	return sum
}

// solve241B mirrors 241B.go to compute expected output for a single test case.
func solve241B(n int, m int64, a []int) int64 {
	groups := []group{{A: a, cross: false}}
	var ans int64
	for b := 30; b >= 0 && m > 0; b-- {
		bitw := int64(1) << uint(b)
		var highCnt int64
		var highGroups, lowGroups []group
		for _, g := range groups {
			if !g.cross {
				var z, o []int
				for _, v := range g.A {
					if v>>uint(b)&1 == 1 {
						o = append(o, v)
					} else {
						z = append(z, v)
					}
				}
				if len(z) > 0 && len(o) > 0 {
					highCnt += int64(len(z)) * int64(len(o))
					highGroups = append(highGroups, group{A: z, B: o, cross: true})
				}
				if len(z) > 1 {
					lowGroups = append(lowGroups, group{A: z, cross: false})
				}
				if len(o) > 1 {
					lowGroups = append(lowGroups, group{A: o, cross: false})
				}
			} else {
				var a0, a1, b0, b1 []int
				for _, v := range g.A {
					if v>>uint(b)&1 == 1 {
						a1 = append(a1, v)
					} else {
						a0 = append(a0, v)
					}
				}
				for _, v := range g.B {
					if v>>uint(b)&1 == 1 {
						b1 = append(b1, v)
					} else {
						b0 = append(b0, v)
					}
				}
				if len(a0) > 0 && len(b1) > 0 {
					highCnt += int64(len(a0)) * int64(len(b1))
					highGroups = append(highGroups, group{A: a0, B: b1, cross: true})
				}
				if len(a1) > 0 && len(b0) > 0 {
					highCnt += int64(len(a1)) * int64(len(b0))
					highGroups = append(highGroups, group{A: a1, B: b0, cross: true})
				}
				if len(a0) > 0 && len(b0) > 0 {
					lowGroups = append(lowGroups, group{A: a0, B: b0, cross: true})
				}
				if len(a1) > 0 && len(b1) > 0 {
					lowGroups = append(lowGroups, group{A: a1, B: b1, cross: true})
				}
			}
		}
		if highCnt >= m {
			ans = (ans + m*bitw) % MOD
			groups = highGroups
		} else {
			ans = (ans + highCnt*bitw) % MOD
			if highCnt > 0 {
				ans = (ans + solveAllLower(highGroups, b-1)) % MOD
			}
			m -= highCnt
			groups = lowGroups
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		expect := solve241B(tc.n, tc.m, tc.a)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
