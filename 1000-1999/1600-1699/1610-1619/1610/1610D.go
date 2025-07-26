package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	cnt := make([]int64, 31)
	for _, v := range arr {
		c := 0
		for v%2 == 0 {
			c++
			v >>= 1
		}
		cnt[c]++
	}
	prefix := make([]int64, 32)
	for i := 30; i >= 0; i-- {
		prefix[i] = cnt[i] + prefix[i+1]
	}
	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}
	res := (pow2[n] - pow2[n-int(cnt[0])] + mod) % mod
	for t := 1; t <= 30; t++ {
		if cnt[t] == 0 {
			continue
		}
		choose := (pow2[int(cnt[t])-1] - 1 + mod) % mod
		res = (res + choose*pow2[int(prefix[t+1])]) % mod
	}
	fmt.Fprintln(out, res)
}
