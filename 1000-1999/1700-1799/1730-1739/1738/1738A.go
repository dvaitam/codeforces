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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		tp := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &tp[i])
		}
		dmg := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &dmg[i])
		}
		var fire, frost []int64
		for i := 0; i < n; i++ {
			if tp[i] == 0 {
				fire = append(fire, dmg[i])
			} else {
				frost = append(frost, dmg[i])
			}
		}
		sort.Slice(fire, func(i, j int) bool { return fire[i] > fire[j] })
		sort.Slice(frost, func(i, j int) bool { return frost[i] > frost[j] })
		var sumFire, sumFrost int64
		for _, x := range fire {
			sumFire += x
		}
		for _, x := range frost {
			sumFrost += x
		}
		if len(fire) == len(frost) {
			if len(fire) == 0 {
				fmt.Fprintln(out, 0)
				continue
			}
			minVal := fire[len(fire)-1]
			if frost[len(frost)-1] < minVal {
				minVal = frost[len(frost)-1]
			}
			total := 2*(sumFire+sumFrost) - minVal
			fmt.Fprintln(out, total)
		} else {
			if len(fire) < len(frost) {
				fire, frost = frost, fire
				sumFire, sumFrost = sumFrost, sumFire
			}
			m := len(frost)
			var extra int64
			for i := 0; i < m; i++ {
				extra += fire[i] + frost[i]
			}
			total := sumFire + sumFrost + extra
			fmt.Fprintln(out, total)
		}
	}
}
