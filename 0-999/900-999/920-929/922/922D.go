package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type strInfo struct {
	s     string
	cntS  int64
	cntH  int64
	noise int64
}

func calcInfo(t string) strInfo {
	var cntS, cntH, noise int64
	for _, ch := range t {
		if ch == 's' {
			cntS++
		} else {
			noise += cntS
			cntH++
		}
	}
	return strInfo{s: t, cntS: cntS, cntH: cntH, noise: noise}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]strInfo, n)
	for i := 0; i < n; i++ {
		var t string
		fmt.Fscan(in, &t)
		arr[i] = calcInfo(t)
	}
	sort.Slice(arr, func(i, j int) bool {
		a := arr[i]
		b := arr[j]
		return a.cntS*b.cntH > b.cntS*a.cntH
	})
	var res, sCount int64
	for _, it := range arr {
		res += it.noise
	}
	for _, it := range arr {
		res += sCount * it.cntH
		sCount += it.cntS
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, res)
	out.Flush()
}
