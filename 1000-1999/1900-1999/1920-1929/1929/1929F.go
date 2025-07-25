package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAXN int = 500005

var fact [MAXN]int64
var invFact [MAXN]int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i < MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[MAXN-1] = modPow(fact[MAXN-1], MOD-2)
	for i := MAXN - 1; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
}

func combLargeN(N int64, K int) int64 {
	if K < 0 || N < int64(K) {
		return 0
	}
	res := int64(1)
	for i := 0; i < K; i++ {
		res = res * ((N - int64(i)) % MOD) % MOD
	}
	res = res * invFact[K] % MOD
	return res
}

func inorder(left, right []int, root int) []int {
	order := make([]int, 0, len(left)-1)
	stack := make([]int, 0)
	cur := root
	for cur != -1 || len(stack) > 0 {
		for cur != -1 {
			stack = append(stack, cur)
			cur = left[cur]
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, cur)
		cur = right[cur]
	}
	return order
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var C int64
		fmt.Fscan(reader, &n, &C)

		left := make([]int, n+1)
		right := make([]int, n+1)
		val := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			var L, R int
			var V int64
			fmt.Fscan(reader, &L, &R, &V)
			left[i] = L
			right[i] = R
			val[i] = V
		}

		order := inorder(left, right, 1)
		knownPos := make([]int, 0)
		knownVal := make([]int64, 0)
		valid := true
		lastVal := int64(-1)
		for idx, node := range order {
			v := val[node]
			if v != -1 {
				if v < 1 || v > C {
					valid = false
					break
				}
				if lastVal != -1 && v < lastVal {
					valid = false
					break
				}
				lastVal = v
				knownPos = append(knownPos, idx)
				knownVal = append(knownVal, v)
			}
		}

		if !valid {
			fmt.Fprintln(writer, 0)
			continue
		}

		ans := int64(1)
		lastIdx := -1
		prevVal := int64(1)
		for i := 0; i < len(knownPos); i++ {
			pos := knownPos[i]
			v := knownVal[i]
			segLen := pos - (lastIdx + 1)
			if v < prevVal {
				ans = 0
				break
			}
			ans = ans * combLargeN(v-prevVal+int64(segLen), segLen) % MOD
			lastIdx = pos
			prevVal = v
		}
		if ans != 0 {
			segLen := n - (lastIdx + 1)
			if prevVal > C {
				ans = 0
			} else {
				ans = ans * combLargeN(C-prevVal+int64(segLen), segLen) % MOD
			}
		}
		fmt.Fprintln(writer, ans%MOD)
	}
}
