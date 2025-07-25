package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	g int
	l int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	const BIG = 1000000007
	segs := make([]pair, 0)
	ans := 0
	res := make([]int, n)

	for i, x := range arr {
		newSegs := []pair{{x, i}}
		for _, p := range segs {
			g := gcd(p.g, x)
			last := &newSegs[len(newSegs)-1]
			if last.g == g {
				if p.l < last.l {
					last.l = p.l
				}
			} else {
				newSegs = append(newSegs, pair{g, p.l})
			}
		}
		bad := false
		for _, p := range newSegs {
			if p.g == i-p.l+1 {
				bad = true
				break
			}
		}
		if bad {
			ans++
			segs = []pair{{BIG, i}}
		} else {
			segs = newSegs
		}
		res[i] = ans
	}

	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
