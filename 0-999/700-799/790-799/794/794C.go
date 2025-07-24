package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s, t string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	n := len(s)
	a := []byte(s)
	b := []byte(t)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })

	a = a[:(n+1)/2]
	b = b[:n/2]

	aL, aR := 0, len(a)-1
	bL, bR := 0, len(b)-1
	res := make([]byte, n)
	l, r := 0, n-1

	for i := 0; i < n; i++ {
		if i%2 == 0 {
			if bL > bR || (aL <= aR && a[aL] < b[bL]) {
				res[l] = a[aL]
				aL++
				l++
			} else {
				res[r] = a[aR]
				aR--
				r--
			}
		} else {
			if aL > aR || (bL <= bR && a[aL] < b[bL]) {
				res[l] = b[bL]
				bL++
				l++
			} else {
				res[r] = b[bR]
				bR--
				r--
			}
		}
	}

	fmt.Println(string(res))
}
