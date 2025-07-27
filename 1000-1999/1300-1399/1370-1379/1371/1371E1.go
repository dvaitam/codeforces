package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, p int
	if _, err := fmt.Fscan(reader, &n, &p); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)

	// prefix counts of numbers <= x
	const maxVal = 4000
	pref := make([]int, maxVal+1)
	for _, v := range a {
		if v > maxVal {
			pref[maxVal]++
		} else {
			pref[v]++
		}
	}
	for i := 1; i <= maxVal; i++ {
		pref[i] += pref[i-1]
	}

	start := 1
	for i, v := range a {
		if t := v - i; t > start {
			start = t
		}
	}
	maxX := a[n-1] - 1
	res := make([]int, 0)
	for x := start; x <= maxX; x++ {
		prod := 1 % p
		for i := 0; i < n; i++ {
			idx := x + i
			if idx > maxVal {
				idx = maxVal
			}
			cnt := pref[idx] - i
			if cnt <= 0 {
				prod = 0
				break
			}
			prod = (prod * (cnt % p)) % p
		}
		if prod%p != 0 {
			res = append(res, x)
		}
	}

	fmt.Fprintln(writer, len(res))
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	if len(res) > 0 {
		fmt.Fprintln(writer)
	}
}
