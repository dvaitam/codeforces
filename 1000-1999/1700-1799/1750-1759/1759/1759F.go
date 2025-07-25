package main

import (
	"bufio"
	"fmt"
	"os"
)

func solveOne(n int, p int, a []int) int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = a[n-1-i]
	}
	arr0 := arr[0]
	union := map[int]bool{}
	earliest := map[int]int{}
	for _, d := range arr {
		union[d] = true
		earliest[d] = 0
	}

	prefix := p - arr0
	if arr0 != 0 && prefix <= p-1 {
		carry := 1
		i := 1
		for {
			if i < n {
				newd := (arr[i] + carry) % p
				t := prefix
				jlsd := (newd - arr0) % p
				if jlsd < 0 {
					jlsd += p
				}
				if jlsd < t {
					t = jlsd
				}
				if old, ok := earliest[newd]; !ok || t < old {
					earliest[newd] = t
				}
				union[newd] = true
				if arr[i]+carry >= p {
					carry = 1
					i++
					continue
				} else {
					break
				}
			} else {
				newd := carry
				t := prefix
				jlsd := (newd - arr0) % p
				if jlsd < 0 {
					jlsd += p
				}
				if jlsd < t {
					t = jlsd
				}
				if old, ok := earliest[newd]; !ok || t < old {
					earliest[newd] = t
				}
				union[newd] = true
				break
			}
		}
	}

	jSet := map[int]bool{}
	for d := range union {
		j := (d - arr0) % p
		if j < 0 {
			j += p
		}
		jSet[j] = true
	}
	j := p - 1
	for ; j >= 0; j-- {
		if !jSet[j] {
			break
		}
	}
	tMissing := 0
	if j >= 0 {
		tMissing = j
	}
	tUnion := 0
	for _, v := range earliest {
		if v > tUnion {
			tUnion = v
		}
	}
	if tUnion > tMissing {
		return tUnion
	}
	return tMissing
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, p int
		fmt.Fscan(in, &n, &p)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		ans := solveOne(n, p, a)
		fmt.Fprintln(out, ans)
	}
}
