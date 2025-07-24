package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var x, k int64
	if _, err := fmt.Fscan(in, &n, &x, &k); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	ans := int64(0)
	for _, val := range arr {
		t := val / x
		b := t - k
		L := b*x + 1
		R := (b + 1) * x
		if R > val {
			R = val
		}
		if L < 1 {
			L = 1
		}
		if R >= L {
			l := sort.Search(len(arr), func(i int) bool { return arr[i] >= L })
			r := sort.Search(len(arr), func(i int) bool { return arr[i] > R })
			ans += int64(r - l)
		}
	}
	fmt.Fprintln(out, ans)
}
