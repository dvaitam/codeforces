package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ p, c int }

const MAXA = 1000000

var spf [MAXA + 1]int

func initSieve() {
	for i := 2; i <= MAXA; i++ {
		if spf[i] == 0 {
			for j := i; j <= MAXA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func factorize(x int) []pair {
	var res []pair
	for x > 1 {
		p := spf[x]
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		res = append(res, pair{p, cnt})
	}
	return res
}

func main() {
	initSieve()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	var ops string
	fmt.Fscan(in, &ops)

	factors := make([][]pair, n)
	for i := 0; i < n; i++ {
		factors[i] = factorize(a[i])
	}

	ans := 0
	for l := 0; l < n; l++ {
		counts := make(map[int]int)
		valid := true
		for r := l; r < n; r++ {
			for _, pc := range factors[r] {
				if ops[r] == '*' {
					counts[pc.p] += pc.c
				} else {
					counts[pc.p] -= pc.c
					if counts[pc.p] < 0 {
						valid = false
					}
				}
			}
			if !valid {
				break
			}
			ans++
		}
	}

	fmt.Fprintln(out, ans)
}
