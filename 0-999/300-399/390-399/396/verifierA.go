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
34 17
1
50
3
42 49 43
1
19
1
49
3
31 50 24
1
49
1
17
3
39 23 34
2
23 34
4
26 21 36 34
3
30 29 37
2
32 31
2
16 6
5
45 49 44 35 4
1
28
2
32 49
3
35 10 29
3
25 35 42
3
47 11 44
1
22
4
5 40 18 26
1
19
2
25 9
2
38 43
2
16 43
1
47
4
3 21 14 22
4
18 24 20 18
5
13 36 25 12 30
5
9 34 9 32 13
1
1
3
3 37 8
5
16 36 15 33 11
4
12 46 35 22
2
43 10
2
23 33
3
17 8 25
4
19 18 18 34
4
7 36 26 5
1
7
5
8 6 9 13 38
3
10 12 14
3
28 3 23
5
41 26 32 5 25
5
7 11 49 15 12
2
42 23
2
14 26
4
24 48 35 37
4
4 14 1 34
3
32 29 44
4
48 27 41 14
5
44 26 8 44 46
4
10 25 29 20
5
4 36 17 30 2
4
20 26 32 35
4
43 10 33 4
4
23 20 32 6
3
32 14 44
3
32 35 24
5
6 20 8 21 34
2
14 31
4
26 33 5 27
1
25
1
13
1
37
2
38 21
3
20 46 9
5
35 14 34 31 20
1
5
4
31 6 42 4
1
2
2
35 39
3
19 39 47
3
30 22 15
1
22
3
39 6 25
4
40 1 1 47
3
40 46 20
1
47
3
12 38 16
3
11 27 29
4
46 11 36 45
3
11 20 31
2
43 28
3
47 5 19
1
38
2
9 22
4
43 19 47 35
2
27 9
3
28 35 36
3
45 27 14
3
1 20 22
2
3 15
3
32 7 36
3
35 43 8
3
46 6 13
5
29 41 1 2 35
3
35 43 2
2
15 43
1
32
4
23 45 25 34
4
37 7 6 23
3
18 18 2
2
19 4
3
48 40 5
3
30 46 47
1
25
1
36
2
7 39
3
29 6 11
5
27 14 17 30 19
3
10 28 19
5
32 42 36 1 47
2
24 48
1
40
5
12 5 6 12 27
4
17 7 36 28
4
49 8 37 18
2
49 34
5
35 46 37 48 42
2
21 20
4
4 17 9 3
2
46 15
3
10 37 23
4
28 35 3 15
5
6 18 30 46 24
3
46 42 32
1
16
2
46 32
4
34 14 3 3
5
26 31 28 44 6
5
37 5 11 24 24
3
18 40 46
4
6 39 11 41
3
20 26 8
2
11 3
3
7 7 1
3
40 25 17
2
42 18
4
34 25 9 32
1
41
3
7 20 29
2
42 8
1
41
2
23 26
1
16
1
16
5
19 14 6 31 31
2
1 26
3
16 28 44
1
5
4
39 1 18 45
2
26 27
1
8
1
11
3
36 31 44
3
10 43 12
1
4
4
18 35 37 12
3
40 30 49
3
5 5 20
2
13 40
5
1 37 14 33 8
5
26 39 40 2 40
4
35 14 36 38
3
1 26 47
2
38 19
5
20 15 30 8 22
1
36
3
24 31 49
1
48
4
22 32 12 3
5
16 49 20 33 18
4
43 9 24 36
3
39 32 40
2
24 22
1
20
5
42 49 44 49 25
2
31 21
2
33 30
2
42 8
4
3 13 31 15
5
35 49 9 47 48
1
42
1
46
3
23 29 42
2
10 11
3
28 34 36
4
49 12 41 30
2
45 46
1
35
3
11 43 18
1
18
5
7 18 28 12 45
5
36 12 12 17 5
2
5 16
2
27 45
3
15 41 16
4
23 11 6 7
5
36 2 34 6 27
3
32 8 2
2
21 40
4
20 44 21 16
3
10 5 20
1
14
5
2 32 26 9 42
2
35 39
2
15 38
5
38 16 17 27 28
1
24
2
35 11
5
17 32 33 40 16
5
41 49 42 15 34
3
7 42 49
3
32 17 24
5
1 21 30 30 25
5
48 37 28 16 2
4
6 23 49 17
3
5 34 40
4
48 1 48 49
3
42 7 45
3
47 14 16
2
32 15
1
48
3
46 47 40
5
48 16 20 22 10
2
22 20
5
49 13 20 8 23
5
14 39 39 41 7
4
45 8 11 10
4
49 49 38 21
4
32 27 10 5
2
45 9
2
47 21
1
48
3
8 41 1
4
20 2 24 43
5
45 32 14 16 10
2
42 13
4
48 10 37 35
2
38 19
4
2 40 17 36
4
4 16 47 19
4
18 5 26 22
5
31 29 19 36 40
2
32 7
1
26
2
7 30
1
7
3
32 36 11
5
43 8 25 9 31
1
37
3
10 27 22
5
5 2 30 6 17
4
33 11 35 4
4
20 20 36 41
4
14 31 11 2
5
38 8 29 27 14
3
34 49 35
5
35 13 14 46 11
4
38 7 13 16
5
10 16 32 17 33
1
35
4
8 21 2 6
5
17 14 12 10 32
4
38 10 32 28
1
15
4
36 35 22 17
2
19 35
4
20 16 19 15
3
4 16 29
2
48 49
5
16 38 2 45 2
1
8
3
5 34 44
3
1 13 21
4
14 48 24 8
1
45
2
32 49
1
4
1
1`

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
