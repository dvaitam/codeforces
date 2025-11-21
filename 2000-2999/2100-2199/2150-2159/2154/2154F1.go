package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353
const maxN = 3000

var fact []int
var invFact []int

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func precomputeFactorials(limit int) {
	fact = make([]int, limit+1)
	invFact = make([]int, limit+1)
	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	invFact[limit] = modPow(fact[limit], mod-2)
	for i := limit; i >= 1; i-- {
		invFact[i-1] = invFact[i] * i % mod
	}
}

func comb(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

// countForK counts completions where numbers 1..k occupy k slots in order.
// Fixed values pin prefix counts; unconstrained stretches contribute binomial factors.
func countForK(arr []int, n, k int) int {
	prefix := make([]int, n+1)
	for i := range prefix {
		prefix[i] = -1
	}
	assign := func(idx, value int) bool {
		if value < 0 || value > k {
			return false
		}
		if prefix[idx] == -1 {
			prefix[idx] = value
			return true
		}
		return prefix[idx] == value
	}
	if !assign(0, 0) || !assign(n, k) {
		return 0
	}
	for idx := 1; idx <= n; idx++ {
		val := arr[idx-1]
		if val == -1 {
			continue
		}
		if val <= k {
			if !assign(idx-1, val-1) || !assign(idx, val) {
				return 0
			}
		} else {
			required := idx + k - val
			if !assign(idx-1, required) || !assign(idx, required) {
				return 0
			}
		}
	}
	if prefix[0] == -1 || prefix[n] == -1 {
		return 0
	}
	result := 1
	lastIdx := 0
	lastVal := prefix[0]
	for i := 1; i <= n; i++ {
		if prefix[i] == -1 {
			continue
		}
		length := i - lastIdx
		delta := prefix[i] - lastVal
		if delta < 0 || delta > length {
			return 0
		}
		result = result * comb(length, delta) % mod
		lastIdx = i
		lastVal = prefix[i]
	}
	if lastIdx != n {
		return 0
	}
	return result
}

func main() {
	precomputeFactorials(maxN)
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		identityPossible := true
		for i := 0; i < n; i++ {
			if arr[i] != -1 && arr[i] != i+1 {
				identityPossible = false
			}
		}
		total := 0
		for k := 1; k <= n-1; k++ {
			total += countForK(arr, n, k)
			if total >= mod {
				total -= mod
			}
		}
		if identityPossible {
			deduction := n - 2
			if deduction < 0 {
				deduction = 0
			}
			total = (total - deduction) % mod
			if total < 0 {
				total += mod
			}
		}
		fmt.Fprintln(writer, total)
	}
}
