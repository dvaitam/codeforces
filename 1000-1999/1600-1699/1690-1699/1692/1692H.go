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

	var tt int
	fmt.Fscan(reader, &tt)
	for tt > 0 {
		tt--
		var n int
		fmt.Fscan(reader, &n)
		v := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &v[i])
		}
		mp := make(map[int64]int64)
		occ := make(map[int64]int64)
		var mx, a, l, r int64
		a, l, r = -1, -1, -1
		for i := 0; i < n; i++ {
			val := v[i]
			if _, ok := mp[val]; ok {
				mp[val]++
			} else {
				mp[val] = 1
				occ[val] = int64(i)
			}
			if mp[val] > mx {
				a = val
				l = occ[val] + 1
				r = int64(i) + 1
				mx = mp[val]
			}
			var toDelete []int64
			for k, cnt := range mp {
				if k != val {
					cnt--
					mp[k] = cnt
					if cnt == 0 {
						toDelete = append(toDelete, k)
					}
				}
			}
			for _, k := range toDelete {
				delete(mp, k)
				delete(occ, k)
			}
		}
		fmt.Fprintf(writer, "%d %d %d\n", a, l, r)
	}
}
