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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		d := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &d[i])
			sum += d[i]
		}
		if k == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		sort.Ints(d)
		since := 0
		watched := 0
		ads := 0
		for _, length := range d {
			since += length
			watched++
			need := false
			if watched == 3 {
				need = true
			}
			if since >= k {
				need = true
			}
			if need {
				ads++
				since = 0
				watched = 0
			}
		}
		fmt.Fprintln(out, ads)
	}
}
