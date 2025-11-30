package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 25
2 45
1 41
5 50
6 15
2 45
6 11
3 42
2 32
1 33
6 38
6 31
1 38
3 13
6 41
4 42
1 3
6 8
4 30
1 47
4 20
4 37
6 5
1 36
2 33
5 6
5 21
6 11
2 15
2 1
5 31
4 46
5 13
1 23
5 12
6 24
3 24
3 44
6 39
1 18
4 41
1 31
6 38
6 29
2 35
1 33
3 27
3 14
4 29
5 5
2 46
1 12
4 43
2 36
4 43
6 14
4 28
6 29
6 33
1 2
5 4
3 13
6 14
6 12
3 35
3 34
6 21
3 26
5 26
1 37
1 14
2 42
2 11
3 21
3 36
4 25
6 23
1 18
6 36
2 6
5 30
2 46
4 38
1 19
2 30
2 40
4 24
6 36
5 37
3 20
2 40
4 11
3 19
1 49
5 1
1 8
1 45
5 9
1 11
5 40`

type testCase struct {
	n int
	k int64
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func modPow(a, e int64) int64 {
	const MOD int64 = 998244353
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func sumPow(k int64, p int) int64 {
	const MOD int64 = 998244353
	m := p + 1
	y := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		y[i] = (y[i-1] + modPow(int64(i), int64(p))) % MOD
	}
	if k <= int64(m) {
		return y[k]
	}
	fact := make([]int64, m+1)
	ifact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	ifact[m] = modPow(fact[m], MOD-2)
	for i := m; i >= 1; i-- {
		ifact[i-1] = ifact[i] * int64(i) % MOD
	}
	pre := make([]int64, m+2)
	suf := make([]int64, m+2)
	pre[0] = 1
	for i := 0; i <= m; i++ {
		pre[i+1] = pre[i] * ((k - int64(i) + MOD) % MOD) % MOD
	}
	suf[m+1] = 1
	for i := m; i >= 0; i-- {
		suf[i] = suf[i+1] * ((k - int64(i) + MOD) % MOD) % MOD
	}
	ans := int64(0)
	for i := 0; i <= m; i++ {
		num := pre[i] * suf[i+1] % MOD
		den := ifact[i] * ifact[m-i] % MOD
		if (m-i)%2 == 1 {
			den = (MOD - den) % MOD
		}
		ans = (ans + y[i]*num%MOD*den) % MOD
	}
	return ans
}

func polyMul(a, b []int64) []int64 {
	const MOD int64 = 998244353
	c := make([]int64, len(a)+len(b)-1)
	for i := 0; i < len(a); i++ {
		if a[i] == 0 {
			continue
		}
		for j := 0; j < len(b); j++ {
			if b[j] == 0 {
				continue
			}
			c[i+j] = (c[i+j] + a[i]*b[j]) % MOD
		}
	}
	return c
}

func solve(n int, k int64) int64 {
	const MOD int64 = 998244353
	if n == 0 {
		return 0
	}
	C := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int64, i+1)
		C[i][0] = 1
		C[i][i] = 1
		for j := 1; j < i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % MOD
		}
	}

	sumPowList := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		sumPowList[i] = sumPow(k, i)
	}
	powK := make([]int64, n+1)
	powK1 := make([]int64, n+1)
	powK[0] = 1
	powK1[0] = 1
	for i := 1; i <= n; i++ {
		powK[i] = powK[i-1] * (k % MOD) % MOD
		powK1[i] = powK1[i-1] * ((k + 1) % MOD) % MOD
	}

	ans := int64(0)
	for i := 1; i <= n; i++ {
		for l := 0; l <= i-1; l++ {
			for g := l + 1; g <= n-i; g++ {
				coef := C[i-1][l] * C[n-i][g] % MOD
				base := n - i - g + 1
				poly := make([]int64, base+1)
				poly[base] = 1

				a := make([]int64, l+1)
				for t := 0; t <= l; t++ {
					val := C[l][t]
					if (l-t)%2 == 1 {
						val = (MOD - val) % MOD
					}
					a[t] = val
				}
				poly = polyMul(poly, a)

				bdeg := i - 1 - l
				b := make([]int64, bdeg+1)
				for s := 0; s <= bdeg; s++ {
					val := C[bdeg][s] * powK1[bdeg-s] % MOD
					if s%2 == 1 {
						val = (MOD - val) % MOD
					}
					b[s] = val
				}
				poly = polyMul(poly, b)

				c := make([]int64, g+1)
				for r := 0; r <= g; r++ {
					val := C[g][r] * powK[g-r] % MOD
					if r%2 == 1 {
						val = (MOD - val) % MOD
					}
					c[r] = val
				}
				poly = polyMul(poly, c)
				sum := int64(0)
				for d := 0; d < len(poly); d++ {
					if poly[d] == 0 {
						continue
					}
					sum = (sum + poly[d]*sumPowList[d]) % MOD
				}
				ans = (ans + coef*sum) % MOD
			}
		}
	}
	return ans % MOD
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		k, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d parse k: %v", idx+1, err)
		}
		cases = append(cases, testCase{n: n, k: k})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		want := solve(tc.n, tc.k)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil || gotVal != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
