package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	l, r := 0, n-1
	var ans int64
	for i := 0; i < n/2; i++ {
		rem := r - l + 1
		if rem != 3 {
			sum := a[l] + a[r]
			ans += sum * sum
			l++
			r--
		} else {
			sum := a[l] + a[l+1] + a[l+2]
			ans += sum * sum
			break
		}
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprint(out, ans)
}
