package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxValue = 10000000

var (
	spf       []int32
	prevPrime []int32
)

func buildSieve(limit int) {
	spf = make([]int32, limit+1)
	prevPrime = make([]int32, limit+1)
	primes := make([]int, 0)
	for i := 2; i <= limit; i++ {
		if spf[i] == 0 {
			spf[i] = int32(i)
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > int(spf[i]) || i*p > limit {
				break
			}
			spf[i*p] = int32(p)
		}
	}
	var last int32
	for i := 2; i <= limit; i++ {
		if spf[i] == int32(i) {
			last = int32(i)
		}
		prevPrime[i] = last
	}
}

func addFactors(x int, freq map[int]int) {
	for x > 1 {
		p := int(spf[x])
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		freq[p] += cnt
	}
}

func vpFactorial(n, p int, cache map[int]int) int {
	if val, ok := cache[p]; ok {
		return val
	}
	res := 0
	tmp := n
	for tmp > 0 {
		tmp /= p
		res += tmp
	}
	cache[p] = res
	return res
}

func main() {
	buildSieve(maxValue)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	const inf = int(1 << 60)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		if n < 2 {
			fmt.Fprintln(out, 0)
			continue
		}
		P := int(prevPrime[n])
		if P == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		total := int64(0)
		if P <= n-1 {
			freq := make(map[int]int)
			vpCache := make(map[int]int)
			for x := n - 1; x >= P; x-- {
				addFactors(x+1, freq)
				best := inf
				for p, delta := range freq {
					vn := vpFactorial(n, p, vpCache)
					vx := vn - delta
					if vx < 0 {
						vx = 0
					}
					cur := int64(p)
					for e := 1; cur <= int64(m) && e <= vn; e++ {
						qx := vx / e
						qn := vn / e
						if qx != qn && qx < best {
							best = qx
							if best == 0 {
								break
							}
						}
						if cur > int64(m)/int64(p) {
							break
						}
						cur *= int64(p)
					}
					if best == 0 {
						break
					}
				}
				if best == inf {
					best = 0
				}
				total += int64(best)
			}
		}
		fmt.Fprintln(out, total)
	}
}
