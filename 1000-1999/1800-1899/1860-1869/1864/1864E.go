package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod int = 998244353

func msb(x int) int {
	if x == 0 {
		return -1
	}
	return bits.Len(uint(x)) - 1
}

func turns(a, b int) int {
	if a == 0 {
		if b == 0 {
			return 1
		}
		return 1
	}
	ha := msb(a)
	hb := msb(b)
	if ha < hb {
		return 1
	}
	if ha > hb {
		return 2
	}
	if a < b {
		return 3
	}
	if a == b {
		pc := bits.OnesCount(uint(a)) + 1
		if pc < 4 {
			return pc
		}
		return 4
	}
	thresh := 3 << (ha - 1)
	if hb == ha && b >= thresh {
		return 4
	}
	return 2
}

func modInv(a int) int {
	return modPow(a, mod-2)
}

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		var sum int
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				sum += turns(arr[i], arr[j])
			}
		}
		inv := modInv(n * n % mod)
		ans := sum % mod * inv % mod
		fmt.Fprintln(out, ans)
	}
}
