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

const MOD845 = 998244353
const G845 = 3

func pow845(a, b int64) int64 {
	res := int64(1)
	a %= MOD845
	for b > 0 {
		if b%2 == 1 {
			res = res * a % MOD845
		}
		a = a * a % MOD845
		b /= 2
	}
	return res
}

func ntt845(a []int64, invert bool) {
	n := len(a)
	j := 0
	for i := 1; i < n; i++ {
		bit := n >> 1
		for j&bit != 0 {
			j ^= bit
			bit >>= 1
		}
		j ^= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for l := 2; l <= n; l <<= 1 {
		half := l >> 1
		wLen := pow845(G845, (MOD845-1)/int64(l))
		if invert {
			wLen = pow845(wLen, MOD845-2)
		}
		for i := 0; i < n; i += l {
			w := int64(1)
			for k := 0; k < half; k++ {
				u := a[i+k]
				v := a[i+k+half] * w % MOD845
				a[i+k] = (u + v) % MOD845
				a[i+k+half] = (u - v + MOD845) % MOD845
				w = w * wLen % MOD845
			}
		}
	}
	if invert {
		nInv := pow845(int64(n), MOD845-2)
		for i := 0; i < n; i++ {
			a[i] = a[i] * nInv % MOD845
		}
	}
}

func divisors845(n int64) []int64 {
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

func primeFactors845(n int64) []int64 {
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
	return res
}

func countCoprimeUpTo845(n, m int64) int64 {
	primes := primeFactors845(m)
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

func solveCase(l, T int64, speeds []int) string {
	maxV := 0
	for _, v := range speeds {
		if v > maxV {
			maxV = v
		}
	}

	size := 1
	for size <= 2*maxV {
		size <<= 1
	}

	A := make([]int64, size)
	B := make([]int64, size)
	for _, v := range speeds {
		A[v] = 1
		B[maxV-v] = 1
	}

	ntt845(A, false)
	ntt845(B, false)

	sumFFT := make([]int64, size)
	diffFFT := make([]int64, size)
	for i := 0; i < size; i++ {
		sumFFT[i] = A[i] * A[i] % MOD845
		diffFFT[i] = A[i] * B[i] % MOD845
	}

	ntt845(sumFFT, true)
	ntt845(diffFFT, true)

	maxS := 2 * maxV
	S := make([]bool, maxS+1)

	for _, v := range speeds {
		sumFFT[2*v] = (sumFFT[2*v] - 1 + MOD845) % MOD845
	}

	for i := 1; i <= maxS; i++ {
		if sumFFT[i] > 0 {
			S[i] = true
		}
	}
	for i := 1; i <= maxV; i++ {
		if diffFFT[maxV+i] > 0 {
			S[i] = true
		}
	}

	hasMultiple := make([]bool, maxS+1)
	for i := 1; i <= maxS; i++ {
		if S[i] {
			hasMultiple[i] = true
		}
	}

	inSPrime := make([]bool, maxS+1)
	for i := 1; i <= maxS; i++ {
		for j := i; j <= maxS; j += i {
			if hasMultiple[j] {
				inSPrime[i] = true
				break
			}
		}
	}

	minPrime := make([]int, maxS+1)
	for i := 2; i <= maxS; i++ {
		if minPrime[i] == 0 {
			for j := i; j <= maxS; j += i {
				if minPrime[j] == 0 {
					minPrime[j] = i
				}
			}
		}
	}

	var ans int64 = 0
	for b := 1; b <= maxS; b++ {
		if !inSPrime[b] {
			continue
		}

		M := (T * int64(b)) / (2 * l)
		if M == 0 {
			continue
		}

		var primes []int64
		temp := b
		for temp > 1 {
			p := minPrime[temp]
			primes = append(primes, int64(p))
			for temp%p == 0 {
				temp /= p
			}
		}

		k := len(primes)
		var count int64 = 0
		for mask := 0; mask < (1 << k); mask++ {
			prod := int64(1)
			bits := 0
			for i := 0; i < k; i++ {
				if (mask & (1 << i)) != 0 {
					prod *= primes[i]
					bits++
				}
			}
			term := M / prod
			if bits%2 == 1 {
				count -= term
			} else {
				count += term
			}
		}
		count %= 1000000007
		ans = (ans + count) % 1000000007
	}

	return fmt.Sprintf("%d", ans)
}

type testCase struct {
	input    string
	expected string
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
