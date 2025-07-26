package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxVal = 1000000

var spf [maxVal + 1]int
var primes []int

func initSieve() {
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || p*i > maxVal {
				break
			}
			spf[p*i] = p
		}
	}
}

func factorSmall(x int) map[int]int {
	m := make(map[int]int)
	for x > 1 {
		p := spf[x]
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		m[p] += cnt
	}
	return m
}

func factorBig(n int64) map[int]int {
	m := make(map[int]int)
	nn := n
	for _, p := range primes {
		pp := int64(p)
		if pp*pp > nn {
			break
		}
		if nn%pp == 0 {
			c := 0
			for nn%pp == 0 {
				nn /= pp
				c++
			}
			m[p] = c
		}
	}
	if nn > 1 {
		m[int(nn)]++
	}
	return m
}

func copyMap(src map[int]int) map[int]int {
	dst := make(map[int]int, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func main() {
	initSieve()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n0, q int
		fmt.Fscan(in, &n0, &q)

		cur := factorSmall(n0)
		dVal := int64(1)
		for _, e := range cur {
			dVal *= int64(e + 1)
		}

		orig := copyMap(cur)
		origD := dVal

		for i := 0; i < q; i++ {
			var typ int
			fmt.Fscan(in, &typ)
			if typ == 1 {
				var x int
				fmt.Fscan(in, &x)
				fac := factorSmall(x)
				for p, c := range fac {
					old := cur[p]
					dVal = dVal / int64(old+1) * int64(old+c+1)
					cur[p] = old + c
				}
				df := factorBig(dVal)
				ok := true
				for p, need := range df {
					if cur[p] < need {
						ok = false
						break
					}
				}
				if ok {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			} else {
				cur = copyMap(orig)
				dVal = origD
			}
		}
	}
}
