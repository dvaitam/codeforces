package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var p, q, r int64
	if _, err := fmt.Fscan(reader, &n, &p, &q, &r); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	pref := make([]int64, n)
	pref[0] = p * a[0]
	for i := 1; i < n; i++ {
		val := p * a[i]
		if pref[i-1] > val {
			pref[i] = pref[i-1]
		} else {
			pref[i] = val
		}
	}

	suff := make([]int64, n)
	suff[n-1] = r * a[n-1]
	for i := n - 2; i >= 0; i-- {
		val := r * a[i]
		if suff[i+1] > val {
			suff[i] = suff[i+1]
		} else {
			suff[i] = val
		}
	}

	ans := pref[0] + q*a[0] + suff[0]
	for i := 1; i < n; i++ {
		val := pref[i] + q*a[i] + suff[i]
		if val > ans {
			ans = val
		}
	}

	fmt.Fprintln(writer, ans)
}
