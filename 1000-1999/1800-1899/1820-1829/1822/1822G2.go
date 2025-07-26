package main

import (
	"bufio"
	"fmt"
	"os"
)

var primes []int

func initPrimes(limit int) {
	sieve := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		if !sieve[i] {
			primes = append(primes, i)
			for j := i * i; j <= limit; j += i {
				sieve[j] = true
			}
		}
	}
}

func divisors(n int) []int {
	res := []int{1}
	x := n
	for _, p := range primes {
		if p*p > x {
			break
		}
		if x%p == 0 {
			cnt := 0
			for x%p == 0 {
				x /= p
				cnt++
			}
			base := res
			res = make([]int, 0, len(base)*(cnt+1))
			pow := 1
			for i := 0; i <= cnt; i++ {
				for _, d := range base {
					res = append(res, d*pow)
				}
				pow *= p
			}
		}
	}
	if x > 1 {
		base := res
		res = make([]int, 0, len(base)*2)
		for _, d := range base {
			res = append(res, d)
			res = append(res, d*x)
		}
	}
	return res
}

func main() {
	initPrimes(31623)
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int)
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &vals[i])
			freq[vals[i]]++
		}
		ans := int64(0)
		for _, c := range freq {
			if c >= 3 {
				ans += int64(c) * int64(c-1) * int64(c-2)
			}
		}
		for y, cy := range freq {
			ds := divisors(y)
			for _, i := range ds {
				if i == y {
					continue
				}
				k64 := int64(y) * int64(y) / int64(i)
				if k64 > 1000000000 {
					continue
				}
				k := int(k64)
				ci, ok1 := freq[i]
				ck, ok2 := freq[k]
				if ok1 && ok2 {
					ans += int64(ci) * int64(cy) * int64(ck)
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
