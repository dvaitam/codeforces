package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	p  [105]int
	c  [105]int
	vd [105]bool
	lk int
)

func clc(n int) int {
	r := 0
	for i := 0; i < n; i++ {
		vd[i] = false
	}
	for i := 0; i < n; i++ {
		if !vd[i] {
			j := i
			for !vd[j] {
				vd[j] = true
				j = p[j]
			}
			r++
		}
	}
	return r
}

func op(x [][]int) {
	fmt.Println(len(x))
	for _, arr := range x {
		n := len(arr)
		for j := 0; j < n; j++ {
			p[j] = arr[j]
		}
		if clc(n) != lk {
			panic("cycle mismatch")
		}
		for _, j := range arr {
			fmt.Print(j+1, " ")
		}
		fmt.Println()
	}
}

func scs(n, k int) {
	lk = k
	ans := [][]int{}
	v := make([]int, n)

	if n <= 20 {
		for i := 0; i < (1 << (n - 1)); i++ {
			t := 0
			for j := 0; j < n-1; j++ {
				if i&(1<<j) != 0 {
					p[t] = j
					t++
				}
			}
			p[t] = n - 1
			t++
			for j := n - 2; j >= 0; j-- {
				if i&(1<<j) == 0 {
					p[t] = j
					t++
				}
			}
			if clc(n) == k {
				for j := 0; j < n; j++ {
					v[j] = p[j]
				}
				ans = append(ans, append([]int(nil), v...))
				if len(ans) >= 2000 {
					break
				}
			}
		}
		op(ans)
		return
	}

	if n-20 <= k-1 {
		g := 20
		k2 := k - (n - g)
		for i := 0; i < (1 << (g - 1)); i++ {
			t := 0
			for j := 0; j < g-1; j++ {
				if i&(1<<j) != 0 {
					p[t] = j
					t++
				}
			}
			p[t] = g - 1
			t++
			for j := g - 2; j >= 0; j-- {
				if i&(1<<j) == 0 {
					p[t] = j
					t++
				}
			}
			if clc(g) == k2 {
				for j := 0; j < n; j++ {
					if j >= n-g {
						v[j] = p[j-(n-g)] + (n - g)
					} else {
						v[j] = j
					}
				}
				ans = append(ans, append([]int(nil), v...))
				if len(ans) >= 2000 {
					break
				}
			}
		}
		op(ans)
		return
	}

	dc := k - 1
	for i := 0; i < (1 << 19); i += 2 {
		t := 0
		for j := 0; j < 19; j++ {
			if i&(1<<j) != 0 {
				p[t] = j
				t++
			}
		}
		p[t] = 19
		t++
		for j := 18; j >= 0; j-- {
			if i&(1<<j) == 0 {
				p[t] = j
				t++
			}
		}
		if clc(20) == 1 {
			for j := 0; j < n; j++ {
				if j >= n-20 {
					v[j] = p[j-(n-20)] + (n - 20)
				} else {
					if j >= dc {
						v[j] = j + 1
					} else {
						v[j] = j
					}
				}
			}
			v[n-1] = dc
			ans = append(ans, append([]int(nil), v...))
			if len(ans) >= 2000 {
				break
			}
		}
	}
	op(ans)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(in, &n, &k)
	scs(n, k)
}

