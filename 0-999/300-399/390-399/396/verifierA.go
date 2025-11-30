package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

var primes []int

func sieve(limit int) []int {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	var res []int
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			res = append(res, i)
		}
	}
	return res
}

func modPow(x, e int64) int64 {
	res := int64(1)
	x %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		e >>= 1
	}
	return res
}

func modInv(x int64) int64 { return modPow(x, mod-2) }

func solveCase(n int, arr []int64) string {
	exp := make(map[int64]int)
	maxE := 0
	for _, v := range arr {
		x := v
		for _, p := range primes {
			pp := int64(p)
			if pp*pp > x {
				break
			}
			if x%pp == 0 {
				cnt := 0
				for x%pp == 0 {
					x /= pp
					cnt++
				}
				exp[pp] += cnt
				if exp[pp] > maxE {
					maxE = exp[pp]
				}
			}
		}
		if x > 1 {
			exp[x]++
			if exp[x] > maxE {
				maxE = exp[x]
			}
		}
	}

	lim := maxE + n - 1
	fact := make([]int64, lim+1)
	inv := make([]int64, lim+1)
	fact[0] = 1
	for i := 1; i <= lim; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[lim] = modInv(fact[lim])
	for i := lim; i >= 1; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}

	res := int64(1)
	for _, e := range exp {
		top := e + n - 1
		ways := fact[top] * inv[n-1] % mod * inv[e] % mod
		res = res * ways % mod
	}
	if res < 0 {
		res += mod
	}
	return strconv.FormatInt(res, 10)
}

var testcasesRaw = `100
4
49 27 3 17
5
32 26 20 31 23
5
14 33 9 19 9
1
40
3
35 46 39
2
20 7
1
44
3
31 36 7
3
28 21 40
2
36 31
4
34 17 4 36
1
6
4
46 43 41 1
5
32 22 16 47 21
1
13
5
15 16 10 35 29
1
6
3
33 32 7
3
36 19 46
1
36
3
35 14 39
5
38 19 29 6 39
4
21 37 16 19
2
13 12
1
40
3
31 5 6
2
10 3
1
45
5
44 26 46 34 18
5
16 14 44 38 27
5
18 29 32 43 42
3
6 21 40
1
32
5
41 22 13 16 2
3
8 46 15
3
11 22 28
1
7
2
45 15
1
37
5
39 44 5 2 8
2
39 37
1
26
1
24
1
3
5
2 13 12 46 8
4
14 47 4 44
1
35
4
40 7 17 5
2
5 42
3
23 28 12
1
33
4
3 39 7 45
4
13 17 23 47
4
37 11 45 44
2
50 4
2
11 22
5
17 8 39 29 43
2
1 31
4
37 33 20 42
3
25 43 17
2
36 45
1
30
1
22
1
35
3
9 16 49
4
23 40 19 44
3
38 41 40
2
46 20
4
48 27 42 6
1
39
2
45 22
2
16 15
4
25 46 44 37
4
3 26 45 37
4
50 43 46 3
2
29 5
3
45 11 29
5
32 36 39 49 1
1
32
3
20 30 4
4
13 36 41 6
2
1 26
4
21 1 14 1
1
44
5
40 7 13 8 39
2
20 18
2
7 31
4
41 6 2 18
4
8 17 9 42
5
42 42 23 8 10
3
2 3 3
2
44 17
5
21 24 37 3 48
5
42 32 46 42 30
4
24 35 12 14
4
38 19 1 9
2
18 22
3
24 46 6
3
50 40 3
1
18`

type testcase struct {
	n   int
	arr []int64
}

func parseTestcases() ([]testcase, error) {
	r := strings.NewReader(testcasesRaw)
	var t int
	if _, err := fmt.Fscan(r, &t); err != nil {
		return nil, err
	}
	cases := make([]testcase, 0, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(r, &n); err != nil {
			return nil, err
		}
		tc := testcase{n: n, arr: make([]int64, n)}
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(r, &tc.arr[j]); err != nil {
				return nil, err
			}
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func main() {
	primes = sieve(31623)

	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solveCase(tc.n, tc.arr)

		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
