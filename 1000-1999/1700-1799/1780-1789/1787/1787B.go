package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type factor struct {
	p int64
	e int64
}

func factorize(n int64) []factor {
	res := make([]factor, 0)
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			cnt := int64(0)
			for n%i == 0 {
				n /= i
				cnt++
			}
			res = append(res, factor{p: i, e: cnt})
		}
	}
	if n > 1 {
		res = append(res, factor{p: n, e: 1})
	}
	return res
}

func maxSum(n int64) int64 {
	fac := factorize(n)
	sort.Slice(fac, func(i, j int) bool { return fac[i].e < fac[j].e })
	prod := int64(1)
	for _, f := range fac {
		prod *= f.p
	}
	ans := int64(0)
	prev := int64(0)
	for _, f := range fac {
		diff := f.e - prev
		if diff > 0 {
			ans += prod * diff
		}
		prod /= f.p
		prev = f.e
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		fmt.Fprintln(writer, maxSum(n))
	}
}
