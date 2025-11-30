package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcaseData = `2 10 5 13 7 18 3 10
6 2 4 2 11 17 10
8 7 6 3 8 11 1 17 13
2 10 5 12 6 10 16 3
8 5 2 16 9
5 5 6 7 16 19 14 18 4
8 9 2 15 2
9 10 3 7 4 15
7 7 6 16 3 14 9 15 6
6 9 2 9 5
5 8 6 8 18 14 15 16 17
3 1 6 15 13 7 10 1 14
5 1 5 15 2 19 8 16
4 6 2 3 6
1 1 2 16 10
9 8 2 2 6
4 4 6 3 12 10 4 6 11
3 3 3 16 9 12
2 2 4 8 17 10 7
7 9 3 11 6 18
2 8 6 18 12 3 7 19 10
5 10 2 12 11
7 1 5 4 9 17 1 5
8 5 3 6 19 15
2 8 2 8 19
2 10 2 6 9
10 2 6 1 5 14 6 10 16
5 8 3 16 15 12
2 4 6 2 14 16 8 10 15
9 6 4 4 9 7 2
2 8 2 6 13
4 8 4 10 1 15 11
7 10 2 16 1
7 8 2 11 10
10 7 5 3 1 2 16 7
2 10 5 7 12 5 16 6
3 3 2 11 3
8 1 2 18 8
3 7 5 12 18 13 10 5
3 9 6 1 3 13 18 11 6
8 1 5 19 15 4 18 7
9 5 5 10 15 11 9 2
5 10 5 6 9 1 17 13
7 1 2 18 9
8 7 2 19 15
10 3 3 7 16 14
9 5 5 14 4 16 15 13
5 4 6 16 19 4 13 1 15
8 2 4 7 13 15 12
4 7 6 13 9 5 1 18 2
8 6 6 9 6 18 10 8 17
10 5 4 11 8 2 16
9 7 6 10 8 6 17 19 12
4 9 2 13 6
4 2 2 13 14
3 2 6 6 8 17 13 11 18
6 3 6 18 19 11 16 7 9
7 2 6 17 6 15 5 14 16
10 6 4 16 10 15 11
5 3 4 3 9 15 7
1 1 4 13 8 6 5
6 5 2 6 15
8 10 5 1 11 17 5 18
6 4 4 18 13 10 4
9 10 3 7 19 8
10 6 4 19 6 13 18
4 5 3 1 17 13
2 4 5 15 5 8 13 14
9 5 4 18 8 11 3
1 1 4 8 12 5 1
9 3 3 13 9 19
7 9 6 10 12 5 14 19 1
3 10 5 19 4 15 13 2
8 10 2 11 15
10 9 5 1 10 9 15 19
1 8 3 19 11 15
4 2 3 15 6 12
4 1 5 4 18 5 10 3
4 1 4 10 15 2 9
6 2 5 4 1 13 15 19
4 8 5 9 14 8 16 6
2 4 6 10 13 5 16 8 17
8 7 5 14 13 8 4 2
7 8 4 15 4 5 12
7 1 5 19 4 10 14 18
4 6 3 12 16 7
4 4 2 13 6
10 3 4 12 19 6 8
2 6 3 16 1 17
6 5 5 15 3 13 6 4
8 9 3 16 10 6
3 4 3 4 8 18
3 8 4 9 2 13 7
2 5 5 13 11 16 19 14
8 4 5 12 14 11 15 9
6 8 4 19 2 6 10
6 3 4 17 2 16 11
6 3 5 3 11 1 15 18
4 1 6 11 18 9 15 6 10
6 2 4 14 1 7 2
4 4 6 4 10 1 19 3 18
4 8 5 3 6 19 2 12
10 5 4 9 2 15 13
4 4 6 1 3 19 4 7 15
3 4 3 5 14 1
4 9 4 19 14 18 1
8 3 5 9 15 19 2 7
10 10 6 17 1 7 10 4 5
4 9 3 5 4 12
8 3 5 4 11 9 12 13
3 8 5 15 6 17 19 1
5 4 6 13 10 7 4 18 19
3 1 2 1 9
3 4 3 16 12 17
9 3 6 8 5 2 15 18 10
1 8 6 19 2 15 18 1 13
4 3 6 10 13 2 5 16 8
10 10 4 11 15 10 9
3 9 6 18 17 8 14 7 9
4 6 4 8 7 9 6`

const (
	mod1 int64 = 998244353
	mod2 int64 = 1004535809
	root int64 = 3
)

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(a, mod int64) int64 {
	return modPow(a, mod-2, mod)
}

func ntt(a []int64, invert bool, mod, root int64) {
	n := len(a)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j |= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		wlen := modPow(root, (mod-1)/int64(length), mod)
		if invert {
			wlen = modInv(wlen, mod)
		}
		for i := 0; i < n; i += length {
			w := int64(1)
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := a[i+j+half] * w % mod
				a[i+j] = u + v
				if a[i+j] >= mod {
					a[i+j] -= mod
				}
				a[i+j+half] = u - v
				if a[i+j+half] < 0 {
					a[i+j+half] += mod
				}
				w = w * wlen % mod
			}
		}
	}
	if invert {
		invN := modInv(int64(n), mod)
		for i := range a {
			a[i] = a[i] * invN % mod
		}
	}
}

func convolution(a, b []int64, mod, root int64) []int64 {
	need := len(a) + len(b) - 1
	n := 1
	for n < need {
		n <<= 1
	}
	fa := make([]int64, n)
	fb := make([]int64, n)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false, mod, root)
	ntt(fb, false, mod, root)
	for i := 0; i < n; i++ {
		fa[i] = fa[i] * fb[i] % mod
	}
	ntt(fa, true, mod, root)
	return fa[:need]
}

var invMod1 int64 = modInv(mod1%mod2, mod2)

func crt(a1, a2 int64) int64 {
	diff := (a2 - a1) % mod2
	if diff < 0 {
		diff += mod2
	}
	x := diff * invMod1 % mod2
	return a1 + x*mod1
}

func divisors(n int64) []int64 {
	res := []int64{}
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			res = append(res, i)
			if i*i != n {
				res = append(res, n/i)
			}
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return res
}

var factorCache = map[int64][]int64{}

func primeFactors(n int64) []int64 {
	if pf, ok := factorCache[n]; ok {
		return pf
	}
	tmp := n
	res := []int64{}
	for p := int64(2); p*p <= tmp; p++ {
		if tmp%p == 0 {
			res = append(res, p)
			for tmp%p == 0 {
				tmp /= p
			}
		}
	}
	if tmp > 1 {
		res = append(res, tmp)
	}
	factorCache[n] = res
	return res
}

func phi(n int64) int64 {
	res := n
	for _, p := range primeFactors(n) {
		res = res / p * (p - 1)
	}
	return res
}

func countCoprimeUpTo(n, m int64) int64 {
	primes := primeFactors(m)
	res := n
	l := len(primes)
	for mask := 1; mask < (1 << l); mask++ {
		mult := int64(1)
		bits := 0
		for i := 0; i < l; i++ {
			if mask&(1<<i) != 0 {
				mult *= primes[i]
				bits++
			}
		}
		if bits%2 == 1 {
			res -= n / mult
		} else {
			res += n / mult
		}
	}
	return res
}

type testCase struct {
	input    string
	expected string
}

func solveCase(l, T int64, speeds []int) string {
	L := l * 2
	maxV := 0
	for _, v := range speeds {
		if v > maxV {
			maxV = v
		}
	}
	size := 1
	for size < maxV+1 {
		size <<= 1
	}
	size <<= 1

	A := make([]int64, size)
	for _, v := range speeds {
		A[v] = 1
	}
	B := make([]int64, size)
	for i := 0; i < size; i++ {
		B[i] = A[size-1-i]
	}

	conv1 := convolution(A, B, mod1, root)
	conv2 := convolution(A, B, mod2, root)

	maxDiff := maxV
	diffPresence := make([]bool, maxDiff+1)
	for d := 1; d <= maxDiff; d++ {
		idx := size - 1 + d
		if idx < len(conv1) && (conv1[idx] != 0 || conv2[idx] != 0) {
			diffPresence[d] = true
		}
	}

	convS1 := convolution(A, A, mod1, root)
	convS2 := convolution(A, A, mod2, root)
	maxSum := maxV * 2
	sumPresence := make([]bool, maxSum+1)
	for s := 0; s <= maxSum && s < len(convS1); s++ {
		cnt := crt(convS1[s], convS2[s])
		if s%2 == 0 {
			v := s / 2
			if v < len(A) && A[v] == 1 {
				if cnt > 0 {
					cnt--
				}
			}
		}
		if cnt > 0 {
			sumPresence[s] = true
		}
	}

	divs := divisors(L)
	diffDivs := []int64{}
	sumDivs := []int64{}
	for _, g := range divs {
		if int64(maxDiff) >= g {
			diffDivs = append(diffDivs, g)
		}
		if int64(maxSum) >= g {
			sumDivs = append(sumDivs, g)
		}
	}

	diffMultiple := map[int64]bool{}
	for _, g := range diffDivs {
		for d := int(g); d <= maxDiff; d += int(g) {
			if diffPresence[d] {
				diffMultiple[g] = true
				break
			}
		}
	}

	sumMultiple := map[int64]bool{}
	for _, g := range sumDivs {
		for s := int(g); s <= maxSum; s += int(g) {
			if sumPresence[s] {
				sumMultiple[g] = true
				break
			}
		}
	}

	sort.Slice(diffDivs, func(i, j int) bool { return diffDivs[i] > diffDivs[j] })
	sort.Slice(sumDivs, func(i, j int) bool { return sumDivs[i] > sumDivs[j] })

	diffExact := map[int64]bool{}
	for i, g := range diffDivs {
		if !diffMultiple[g] {
			continue
		}
		ok := true
		for j := 0; j < i; j++ {
			h := diffDivs[j]
			if h%g == 0 && diffExact[h] {
				ok = false
				break
			}
		}
		if ok {
			diffExact[g] = true
		}
	}

	sumExact := map[int64]bool{}
	for i, g := range sumDivs {
		if !sumMultiple[g] {
			continue
		}
		ok := true
		for j := 0; j < i; j++ {
			h := sumDivs[j]
			if h%g == 0 && sumExact[h] {
				ok = false
				break
			}
		}
		if ok {
			sumExact[g] = true
		}
	}

	flag := map[int64]bool{}
	for g := range diffExact {
		m := L / g
		for _, d := range divisors(g) {
			flag[m*d] = true
		}
	}
	for g := range sumExact {
		m := L / g
		for _, d := range divisors(g) {
			flag[m*d] = true
		}
	}

	countUpTo := func(r int64) int64 {
		if r <= 0 {
			return 0
		}
		var ans int64
		for _, g := range divs {
			if !flag[g] {
				continue
			}
			m := L / g
			ans += countCoprimeUpTo(r/g, m)
		}
		return ans
	}

	totalInPeriod := countUpTo(L)
	ans := (T / L % mod1) * (totalInPeriod % mod1) % mod1
	ans = (ans + countUpTo(T%L)%mod1) % mod1
	return fmt.Sprintf("%d", ans)
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("case %d: not enough tokens", idx+1)
		}
		l, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad l: %w", idx+1, err)
		}
		T, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad T: %w", idx+1, err)
		}
		n, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", idx+1, err)
		}
		if len(fields) != 3+n {
			return nil, fmt.Errorf("case %d: expected %d speeds got %d", idx+1, n, len(fields)-3)
		}
		speeds := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[3+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad speed: %w", idx+1, err)
			}
			speeds[i] = v
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n%d\n", l, T, n)
		for i, v := range speeds {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solveCase(l, T, speeds),
		})
	}
	return cases, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
