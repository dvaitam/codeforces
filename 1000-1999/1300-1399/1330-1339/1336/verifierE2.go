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
	mod          = 998244353
	testcasesRaw = `4 6 12 53 7 37
3 4 8 3 10
1 2 0
2 5 21 28
4 3 5 7 4 2
8 6 1 13 0 13 42 26 39 14
5 1 0 1 0 0 0
3 7 42 42 119
5 2 3 0 3 0 0
2 4 2 12
4 4 12 7 14 11
3 6 1 53 29
5 5 20 30 29 16 22
5 2 1 1 3 2 0
6 7 73 63 5 26 117 46
4 6 42 17 4 62
7 7 112 94 3 50 52 17 125
3 1 1 1 0
2 6 2 61
6 0 0 0 0 0 0 0
1 1 1
7 6 48 29 36 62 51 0 55
5 7 39 102 110 101 77
8 8 189 155 214 87 177 19 217 224
8 4 2 7 10 7 12 0 10 4
3 7 17 25 45
1 3 2
1 8 223
7 2 2 3 3 2 0 1 3
3 1 0 1 1
5 4 12 1 8 10 14
2 0 0 0
8 4 6 12 2 13 1 1 2 11
1 4 9
1 0 0
1 2 1
8 1 0 0 0 0 0 0 0 0
2 7 110 20
2 7 114 43
2 2 2 3
3 1 0 0 0
7 5 20 25 19 15 5 9 25
6 4 11 11 0 0 0 4
5 3 2 2 1 7 7
3 4 8 4 1
3 5 20 12 6
8 2 2 2 1 0 1 2 3 0
5 6 45 4 42 29 50
1 6 34
3 4 1 8 2
6 6 0 34 29 22 32 47
2 2 3 1
2 6 36 5
5 0 0 0 0 0 0
4 8 101 16 97 232
8 5 22 23 7 1 26 23 11 30
4 6 1 24 63 52
7 2 1 3 3 3 1 2 0
3 0 0 0 0
3 4 5 12 8
7 5 27 31 29 7 16 20 16
5 4 12 0 15 6 13
5 6 4 2 53 10 35
1 0 0
7 8 179 241 227 243 94 12 182
4 4 13 2 7 3
3 7 32 127 80
8 2 3 2 3 2 1 0 2 2
7 8 49 157 66 71 74 32 174
4 0 0 0 0 0
1 2 3
3 4 2 9 11
3 7 94 10 84
2 1 0 1
5 5 11 0 30 25 25
4 1 1 1 1 1
8 6 18 54 26 54 24 34 9 17
8 5 23 18 28 21 15 10 29 8
6 3 7 2 1 7 5 2
2 5 10 17
6 1 0 1 1 1 1 1
1 3 4
3 2 1 3 1
2 8 132 54
1 6 11
7 7 124 36 36 14 28 38 30
7 1 1 0 0 0 1 1 0
8 4 10 5 12 8 5 5 2 6
6 7 23 84 108 77 117 14
3 5 11 31 14
4 1 0 1 1 0
6 3 1 5 2 3 3 5
8 5 30 17 18 1 11 17 12 7
5 7 44 28 12 116 96
7 5 21 14 21 31 9 1 30
3 0 0 0 0
2 0 0 0
4 1 0 0 0 1
7 8 188 214 107 200 20 114 3
8 0 0 0 0 0 0 0 0 0`
)

type testCase struct {
	n int
	m int
	a []uint64
}

func powMod(a, e int) int {
	res := 1
	base := a
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		e >>= 1
	}
	return res
}

type state struct {
	idx int
	val uint64
}

// Embedded reference logic from 1336E2.go.
func solve(tc testCase) []int {
	n, m := tc.n, tc.m
	arr := tc.a

	enumerate := func(vects []uint64) []int {
		res := make([]int, m+1)
		stack := []state{{0, 0}}
		k := len(vects)
		for len(stack) > 0 {
			s := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if s.idx == k {
				res[bits.OnesCount64(s.val)]++
			} else {
				stack = append(stack, state{s.idx + 1, s.val})
				stack = append(stack, state{s.idx + 1, s.val ^ vects[s.idx]})
			}
		}
		return res
	}

	basis := make([]uint64, m)
	for _, v := range arr {
		x := v
		for i := m - 1; i >= 0; i-- {
			if ((x >> uint(i)) & 1) == 0 {
				continue
			}
			if basis[i] != 0 {
				x ^= basis[i]
			} else {
				basis[i] = x
				break
			}
		}
	}

	for i := 0; i < m; i++ {
		if basis[i] != 0 {
			for j := i + 1; j < m; j++ {
				if (basis[j]>>uint(i))&1 == 1 {
					basis[j] ^= basis[i]
				}
			}
		}
	}

	var vectors []uint64
	var pivots []int
	for i := 0; i < m; i++ {
		if basis[i] != 0 {
			vectors = append(vectors, basis[i])
			pivots = append(pivots, i)
		}
	}
	r := len(vectors)
	d := m - r

	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = int(int64(pow2[i-1]) * 2 % mod)
	}

	if r <= d {
		cnt := enumerate(vectors)
		for i := 0; i <= m; i++ {
			cnt[i] = int(int64(cnt[i]) * int64(pow2[n-r]) % mod)
		}
		return cnt
	}

	used := make([]bool, m)
	for _, p := range pivots {
		used[p] = true
	}
	var dual []uint64
	for f := 0; f < m; f++ {
		if used[f] {
			continue
		}
		x := uint64(1) << uint(f)
		for idx, p := range pivots {
			if (vectors[idx]>>uint(f))&1 == 1 {
				x |= uint64(1) << uint(p)
			}
		}
		dual = append(dual, x)
	}

	cntDual := enumerate(dual)

	comb := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		comb[i] = make([]int, i+1)
		comb[i][0] = 1
		comb[i][i] = 1
		for j := 1; j < i; j++ {
			comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % mod
		}
	}

	inv := powMod(pow2[d], mod-2)
	ans := make([]int, m+1)
	for i := 0; i <= m; i++ {
		sum := 0
		for j := 0; j <= m; j++ {
			kcoef := 0
			maxT := i
			if j < maxT {
				maxT = j
			}
			for t := 0; t <= maxT; t++ {
				if i-t > m-j {
					continue
				}
				term := int(int64(comb[j][t]) * int64(comb[m-j][i-t]) % mod)
				if t%2 == 1 {
					term = (mod - term) % mod
				}
				kcoef += term
				if kcoef >= mod {
					kcoef -= mod
				}
			}
			sum = (sum + int(int64(cntDual[j])*int64(kcoef)%mod)) % mod
		}
		ans[i] = int(int64(sum) * int64(inv) % mod)
		ans[i] = int(int64(ans[i]) * int64(pow2[n-r]) % mod)
	}

	return ans
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid m", idx+1)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d values, got %d", idx+1, n, len(fields)-2)
		}
		arr := make([]uint64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid value", idx+1)
			}
			arr[i] = uint64(v)
		}
		res = append(res, testCase{n: n, m: m, a: arr})
	}
	return res, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')

		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		outFields := strings.Fields(got)
		if len(outFields) != tc.m+1 {
			fmt.Fprintf(os.Stderr, "case %d wrong number of outputs: expected %d got %d\n", idx+1, tc.m+1, len(outFields))
			os.Exit(1)
		}
		for i := 0; i <= tc.m; i++ {
			v, err := strconv.Atoi(outFields[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d invalid output\n", idx+1)
				os.Exit(1)
			}
			if v%mod != expect[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at index %d: expected %d got %d\n", idx+1, i, expect[i], v)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
