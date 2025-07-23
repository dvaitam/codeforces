package main

import (
	"bufio"
	"fmt"
	"os"
)

func digitsNeeded(x int) int {
	d := 1
	for x >= 7 {
		x /= 7
		d++
	}
	return d
}

func toBase7(x, d int) []int {
	res := make([]int, d)
	for i := d - 1; i >= 0; i-- {
		res[i] = x % 7
		x /= 7
	}
	return res
}

func distinct(hd, md []int) bool {
	var used [7]bool
	for _, v := range hd {
		if used[v] {
			return false
		}
		used[v] = true
	}
	for _, v := range md {
		if used[v] {
			return false
		}
		used[v] = true
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)

	d1 := digitsNeeded(n - 1)
	d2 := digitsNeeded(m - 1)

	if d1+d2 > 7 {
		fmt.Println(0)
		return
	}

	ans := 0
	for h := 0; h < n; h++ {
		hd := toBase7(h, d1)
		for mm := 0; mm < m; mm++ {
			md := toBase7(mm, d2)
			if distinct(hd, md) {
				ans++
			}
		}
	}
	fmt.Println(ans)
}
