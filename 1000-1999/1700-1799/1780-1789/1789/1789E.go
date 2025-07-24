package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

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
		sn := arr[n-1]
		freq := make([]byte, sn+1)
		for _, v := range arr {
			freq[v] = 1
		}
		prefix := make([]int, sn+1)
		for i := 1; i <= sn; i++ {
			prefix[i] = prefix[i-1] + int(freq[i])
		}
		modCount := make([]int, sn+1)
		for d := 1; d <= sn; d++ {
			c := 0
			for k := 0; k*d <= sn; k++ {
				left := k * d
				rmax := k
				if rmax > d-1 {
					rmax = d - 1
				}
				right := left + rmax
				if right > sn {
					right = sn
				}
				if right >= left {
					if left == 0 {
						c += prefix[right]
					} else {
						c += prefix[right] - prefix[left-1]
					}
				}
			}
			modCount[d] = c
		}
		divisors := make([]int, 0)
		for i := 1; i*i <= sn; i++ {
			if sn%i == 0 {
				divisors = append(divisors, i)
				if i*i != sn {
					divisors = append(divisors, sn/i)
				}
			}
		}
		divCount := make(map[int]int)
		for _, d := range divisors {
			c := 0
			for j := d; j <= sn; j += d {
				if freq[j] == 1 {
					c++
				}
			}
			divCount[d] = c
		}
		ans := int64(0)
		for d := 1; d <= sn; d++ {
			L := sn/(d+1) + 1
			R := sn / d
			if L > R {
				continue
			}
			base := modCount[d]
			lenRange := R - L + 1
			sumX := int64(L+R) * int64(lenRange) / 2 % MOD
			ans = (ans + int64(base)%MOD*sumX) % MOD
			if sn%d == 0 {
				x0 := R
				diff := divCount[d] - base
				ans = (ans + int64(diff)%MOD*int64(x0)%MOD) % MOD
			}
		}
		if ans < 0 {
			ans += MOD
		}
		fmt.Fprintln(out, ans%MOD)
	}
}
