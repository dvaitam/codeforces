package main

import (
	"bufio"
	"fmt"
	"os"
)

func countPairs(a []int) int64 {
	n := len(a)
	triples := make([][3]int, n-2)
	for i := 0; i < n-2; i++ {
		triples[i] = [3]int{a[i], a[i+1], a[i+2]}
	}
	var ans int64
	type key struct{ x, y int }
	ab := make(map[key]int)
	bc := make(map[key]int)
	ac := make(map[key]int)
	abc := make(map[[3]int]int)
	for _, t := range triples {
		k1 := key{t[0], t[1]}
		k2 := key{t[1], t[2]}
		k3 := key{t[0], t[2]}
		ans += int64(ab[k1] - abc[t])
		ans += int64(bc[k2] - abc[t])
		ans += int64(ac[k3] - abc[t])
		ab[k1]++
		bc[k2]++
		ac[k3]++
		abc[t]++
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		fmt.Fprintln(out, countPairs(arr))
	}
}
