package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		var k int
		fmt.Fscan(in, &n, &k)
		ans := int64(0)
		i := int64(0)
		for i < n {
			if bits.OnesCount64(uint64(i)) > k {
				i++
				continue
			}
			j := i
			for j < n && bits.OnesCount64(uint64(j)) <= k {
				j++
			}
			length := j - i
			ans = (ans + (length*(length+1)/2)%mod) % mod
			i = j
		}
		fmt.Fprintln(out, ans%mod)
	}
}
