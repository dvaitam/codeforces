package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	MOD           int64 = 998244353
	testcasesData       = `
3 5 23 17 0
1 0 0
5 6 30 50 63 46 45
2 4 13 3
1 6 22
1 4 8
6 5 13 24 12 20 11 11
8 5 9 26 14 0 3 27 2 8
2 3 4 5
4 5 3 21 8 20
5 3 7 5 4 1 7
3 0 0 0 0
7 5 10 26 8 7 7 11 30
5 5 6 27 16 17 7
6 4 2 0 1 0 4 3
7 3 7 0 5 2 1 6 3
1 1 1
1 0 0
3 4 9 13 12
1 5 28
3 4 4 6 5
2 4 15 9
5 2 1 3 0 0 0
1 5 18
3 6 18 17 14
6 0 0 0 0 0 0 0
4 2 3 2 3 3
2 0 0 0
7 6 51 61 54 9 28 18 51
6 1 0 0 0 0 1 0
2 4 8 9
6 5 30 25 20 29 22 14
1 5 24
6 2 0 1 0 3 0 2
7 6 31 3 14 5 45 5 46
3 2 2 3 3
7 0 0 0 0 0 0 0 0
7 5 18 30 13 1 29 11 28
6 4 10 13 4 13 12 7
6 3 4 3 3 4 3 5
8 1 1 0 1 0 0 1 0 0
8 3 3 2 5 3 1 7 3 4
7 2 1 2 3 2 0 0 1
4 4 5 14 2 8
5 5 31 26 20 24 3
2 0 0 0
5 4 14 9 13 15 11
4 5 5 5 5 1
6 5 6 15 1 2 6 19
2 6 50 16
4 2 2 1 1 1
3 1 0 1 1
5 1 1 0 0 0 1
3 1 0 1 1
3 6 34 10 35
4 0 0 0 0 0
3 1 1 0 1
4 5 11 29 20 28
4 5 17 12 6 22
8 2 2 2 1 2 3 0 1 3
2 4 3 4
8 1 1 1 0 0 0 1 1 1
7 4 1 11 12 12 4 4 1
5 4 14 13 6 14 6
7 0 0 0 0 0 0 0 0
4 1 0 1 1 0
4 4 6 6 15 2
4 6 12 16 24 55
7 6 33 60 3 44 44 26 4
4 3 3 3 0 2
4 3 7 1 4 4
7 2 2 0 0 3 2 0 1
3 3 3 0 7
5 1 1 0 0 1 1
2 2 1 1
5 4 13 12 12 1 4
8 0 0 0 0 0 0 0 0 0
5 5 13 7 10 8 28
7 1 0 1 1 0 1 0 0
8 1 0 0 0 1 1 1 1 0
4 5 7 10 31 15
1 1 0
3 6 19 22 8
2 2 0 2
8 4 8 13 1 2 4 14 1 11
6 2 3 3 0 2 1 1
2 4 4 7
4 3 5 0 4 6
3 5 21 24 12
5 0 0 0 0 0 0
2 2 1 3
1 3 7
1 5 18
7 3 3 2 0 4 1 6 4
4 3 7 4 6 2
8 2 1 0 1 2 2 0 0 3
5 2 1 2 3 1 0
4 2 3 3 1 2
7 6 55 52 46 13 11 57 16
2 0 0 0
`
)

type testCase struct {
	n   int
	m   int
	arr []uint64
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func powmod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

// buildBasis returns basis vectors in RREF order and pivot indices.
func buildBasis(arr []uint64, m int) ([]uint64, []int) {
	basis := make([]uint64, m)
	for _, x := range arr {
		v := x
		for i := m - 1; i >= 0; i-- {
			if (v>>i)&1 == 0 {
				continue
			}
			if basis[i] != 0 {
				v ^= basis[i]
			} else {
				basis[i] = v
				for j := i + 1; j < m; j++ {
					if basis[j] != 0 && ((basis[j]>>i)&1) == 1 {
						basis[j] ^= v
					}
				}
				break
			}
		}
	}
	for i := m - 1; i >= 0; i-- {
		if basis[i] == 0 {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if basis[j] != 0 && ((basis[j]>>i)&1) == 1 {
				basis[j] ^= basis[i]
			}
		}
	}
	pivots := make([]int, 0)
	for i := 0; i < m; i++ {
		if basis[i] != 0 {
			pivots = append(pivots, i)
		}
	}
	vecs := make([]uint64, len(pivots))
	for idx, p := range pivots {
		vecs[idx] = basis[p]
	}
	return vecs, pivots
}

func dualBasis(basis []uint64, pivots []int, m int) []uint64 {
	pivotSet := make(map[int]bool)
	for _, p := range pivots {
		pivotSet[p] = true
	}
	dual := make([]uint64, 0, m-len(pivots))
	for j := 0; j < m; j++ {
		if pivotSet[j] {
			continue
		}
		w := uint64(1) << j
		for idx, p := range pivots {
			if ((basis[idx] >> j) & 1) == 1 {
				w |= uint64(1) << p
			}
		}
		dual = append(dual, w)
	}
	return dual
}

func enumerate(vecs []uint64, idx int, cur uint64, res []int64) {
	if idx == len(vecs) {
		w := bits.OnesCount64(cur)
		res[w]++
		return
	}
	enumerate(vecs, idx+1, cur, res)
	enumerate(vecs, idx+1, cur^vecs[idx], res)
}

func macWilliams(dual []int64, m int, r int) []int64 {
	d := m - r
	comb := make([][]int64, m+1)
	for i := 0; i <= m; i++ {
		comb[i] = make([]int64, m+1)
		comb[i][0] = 1
		for j := 1; j <= i; j++ {
			comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % MOD
		}
	}
	res := make([]int64, m+1)
	powInv := powmod(2, int64(MOD-1)-int64(d)) // inverse of 2^d
	for k := 0; k <= m; k++ {
		var sum int64
		for i := 0; i <= m; i++ {
			if dual[i] == 0 {
				continue
			}
			var kk int64
			for j := 0; j <= k && j <= i; j++ {
				sign := int64(1)
				if j%2 == 1 {
					sign = MOD - 1
				}
				kk = (kk + sign*comb[i][j]%MOD*comb[m-i][k-j]) % MOD
			}
			sum = (sum + dual[i]*kk) % MOD
		}
		res[k] = sum * powInv % MOD
	}
	return res
}

func solve(arr []uint64, m int) []int64 {
	vecs, pivots := buildBasis(arr, m)
	r := len(vecs)
	pow2 := powmod(2, int64(len(arr)-r))

	ans := make([]int64, m+1)
	if r <= m-r {
		enumerate(vecs, 0, 0, ans)
		for i := range ans {
			ans[i] = ans[i] * pow2 % MOD
		}
	} else {
		dualVecs := dualBasis(vecs, pivots, m)
		dualCnt := make([]int64, m+1)
		enumerate(dualVecs, 0, 0, dualCnt)
		ans = macWilliams(dualCnt, m, r)
		for i := range ans {
			ans[i] = ans[i] * pow2 % MOD
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("case %d invalid line", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("case %d invalid number of values", idx+1)
		}
		arr := make([]uint64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, err
			}
			arr[i] = uint64(v)
		}
		cases = append(cases, testCase{n: n, m: m, arr: arr})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		want := solve(tc.arr, tc.m)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatUint(v, 10))
		}
		input.WriteByte('\n')

		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		outFields := strings.Fields(got)
		if len(outFields) != tc.m+1 {
			fmt.Fprintf(os.Stderr, "case %d wrong number of outputs\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i <= tc.m; i++ {
			v, err := strconv.ParseInt(outFields[i], 10, 64)
			if err != nil || (v%MOD+MOD)%MOD != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at index %d: expected %d got %s\n", idx+1, i, want[i], outFields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
