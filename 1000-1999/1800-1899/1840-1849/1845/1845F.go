package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	MOD  int64 = 1000000007
	MOD1 int64 = 998244353
	MOD2 int64 = 1004535809
	ROOT int64 = 3
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

var invMod1 int64 = modInv(MOD1%MOD2, MOD2)

func crt(a1, a2 int64) int64 {
	diff := (a2 - a1) % MOD2
	if diff < 0 {
		diff += MOD2
	}
	x := diff * invMod1 % MOD2
	return a1 + x*MOD1
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var l, T int64
	fmt.Fscan(in, &l, &T)
	L := l * 2
	var n int
	fmt.Fscan(in, &n)
	speeds := make([]int, n)
	maxV := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &speeds[i])
		if speeds[i] > maxV {
			maxV = speeds[i]
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

	conv1 := convolution(A, B, MOD1, ROOT)
	conv2 := convolution(A, B, MOD2, ROOT)

	maxDiff := maxV
	diffPresence := make([]bool, maxDiff+1)
	for d := 1; d <= maxDiff; d++ {
		idx := size - 1 + d
		if idx < len(conv1) && (conv1[idx] != 0 || conv2[idx] != 0) {
			diffPresence[d] = true
		}
	}

	convS1 := convolution(A, A, MOD1, ROOT)
	convS2 := convolution(A, A, MOD2, ROOT)
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
	ans := (T / L % MOD) * (totalInPeriod % MOD) % MOD
	ans = (ans + countUpTo(T%L)%MOD) % MOD
	fmt.Fprintln(out, ans)
}
