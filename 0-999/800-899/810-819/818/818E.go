package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	if k == 1 {
		// every subarray works
		total := int64(n) * int64(n+1) / 2
		fmt.Fprintln(out, total)
		return
	}

	// factorize k
	primes := []int{}
	exps := []int{}
	tmp := k
	for p := 2; p*p <= tmp; p++ {
		if tmp%p == 0 {
			cnt := 0
			for tmp%p == 0 {
				tmp /= p
				cnt++
			}
			primes = append(primes, p)
			exps = append(exps, cnt)
		}
	}
	if tmp > 1 {
		primes = append(primes, tmp)
		exps = append(exps, 1)
	}
	m := len(primes)

	cur := make([]int, m)
	ans := int64(0)
	left := 0

	for right := 0; right < n; right++ {
		// add a[right]
		val := a[right]
		for j := 0; j < m; j++ {
			c := 0
			p := primes[j]
			for val%p == 0 {
				val /= p
				c++
			}
			cur[j] += c
		}
		// shrink from left while window valid
		for left <= right {
			ok := true
			for j := 0; j < m; j++ {
				if cur[j] < exps[j] {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
			ans += int64(n - right)
			val2 := a[left]
			for j := 0; j < m; j++ {
				c := 0
				p := primes[j]
				for val2%p == 0 {
					val2 /= p
					c++
				}
				cur[j] -= c
			}
			left++
		}
	}

	fmt.Fprintln(out, ans)
}
