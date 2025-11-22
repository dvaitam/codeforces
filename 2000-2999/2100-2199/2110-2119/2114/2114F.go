package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxV = 1_000_000
const inf = int(1e9)

func sieveSPF(n int) []int {
	spf := make([]int, n+1)
	for i := 0; i <= n; i++ {
		spf[i] = i
	}
	for i := 2; i*i <= n; i++ {
		if spf[i] == i {
			for j := i * i; j <= n; j += i {
				if spf[j] == j {
					spf[j] = i
				}
			}
		}
	}
	return spf
}

func divisors(n int, spf []int) []int {
	divs := []int{1}
	for n > 1 {
		p := spf[n]
		cnt := 0
		for n%p == 0 {
			n /= p
			cnt++
		}
		cur := make([]int, len(divs))
		copy(cur, divs)
		mult := 1
		for i := 0; i < cnt; i++ {
			mult *= p
			for _, v := range cur {
				divs = append(divs, v*mult)
			}
		}
	}
	return divs
}

func minOps(n, k int, spf []int) int {
	if n == 1 {
		return 0
	}
	if k == 1 {
		return -1
	}
	divs := divisors(n, spf)
	sort.Ints(divs)
	idx := make(map[int]int, len(divs))
	for i, v := range divs {
		idx[v] = i
	}
	lim := make([]int, 0, len(divs))
	for _, d := range divs {
		if d > 1 && d <= k {
			lim = append(lim, d)
		} else if d > k {
			break
		}
	}

	dp := make([]int, len(divs))
	for i := range dp {
		dp[i] = inf
	}
	dp[idx[1]] = 0

	for _, v := range divs[1:] {
		best := inf
		for _, d := range lim {
			if d > v {
				break
			}
			if v%d == 0 {
				u := v / d
				cur := dp[idx[u]] + 1
				if cur < best {
					best = cur
				}
			}
		}
		dp[idx[v]] = best
	}
	res := dp[idx[n]]
	if res >= inf {
		return -1
	}
	return res
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	in := newScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	spf := sieveSPF(maxV)

	t := in.nextInt()
	for ; t > 0; t-- {
		x := in.nextInt()
		y := in.nextInt()
		k := in.nextInt()

		g := gcd(x, y)
		stepsX := minOps(x/g, k, spf)
		if stepsX == -1 {
			fmt.Fprintln(out, -1)
			continue
		}
		stepsY := minOps(y/g, k, spf)
		if stepsY == -1 {
			fmt.Fprintln(out, -1)
			continue
		}
		fmt.Fprintln(out, stepsX+stepsY)
	}
}
