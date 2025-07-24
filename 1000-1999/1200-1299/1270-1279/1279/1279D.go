package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func readInt(r *bufio.Reader) int {
	sign := 1
	num := 0
	c, err := r.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		num = num*10 + int(c-'0')
		c, err = r.ReadByte()
		if err != nil {
			break
		}
	}
	return num * sign
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

var invCache = make(map[int]int64)

func modInv(x int) int64 {
	if v, ok := invCache[x]; ok {
		return v
	}
	v := modPow(int64(x), mod-2)
	invCache[x] = v
	return v
}

func main() {
	rd := bufio.NewReader(os.Stdin)
	n := readInt(rd)
	lists := make([][]int, n)
	cnt := make([]int, 1000001)
	for i := 0; i < n; i++ {
		k := readInt(rd)
		lists[i] = make([]int, k)
		for j := 0; j < k; j++ {
			x := readInt(rd)
			lists[i][j] = x
			cnt[x]++
		}
	}

	invN := modInv(n)
	var ans int64
	for i := 0; i < n; i++ {
		k := len(lists[i])
		invK := modInv(k)
		for _, item := range lists[i] {
			ans = (ans + invK*int64(cnt[item])%mod) % mod
		}
	}
	ans = ans * invN % mod
	ans = ans * invN % mod

	wr := bufio.NewWriter(os.Stdout)
	defer wr.Flush()
	fmt.Fprintln(wr, ans)
}
