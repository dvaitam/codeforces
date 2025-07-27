package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(indices []int, in *bufio.Reader, out *bufio.Writer) int {
	fmt.Fprintf(out, "? %d", len(indices))
	for _, v := range indices {
		fmt.Fprintf(out, " %d", v)
	}
	fmt.Fprintln(out)
	out.Flush()
	var resp int
	fmt.Fscan(in, &resp)
	return resp
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		subsets := make([][]int, k)
		for i := 0; i < k; i++ {
			var c int
			fmt.Fscan(in, &c)
			subsets[i] = make([]int, c)
			for j := 0; j < c; j++ {
				fmt.Fscan(in, &subsets[i][j])
			}
		}

		all := make([]int, n)
		for i := 0; i < n; i++ {
			all[i] = i + 1
		}
		globalMax := query(all, in, out)

		unionSet := []int{}
		for _, s := range subsets {
			unionSet = append(unionSet, s...)
		}
		unionMax := 0
		if len(unionSet) > 0 {
			unionMax = query(unionSet, in, out)
		}

		ans := make([]int, k)
		if unionMax < globalMax {
			for i := 0; i < k; i++ {
				ans[i] = globalMax
			}
		} else {
			left, right := 0, k-1
			for left < right {
				mid := (left + right) / 2
				q := []int{}
				for i := left; i <= mid; i++ {
					q = append(q, subsets[i]...)
				}
				if len(q) == 0 {
					left = mid + 1
					continue
				}
				if query(q, in, out) == globalMax {
					right = mid
				} else {
					left = mid + 1
				}
			}
			pos := left

			m := make(map[int]bool)
			for _, v := range subsets[pos] {
				m[v] = true
			}
			comp := []int{}
			for i := 1; i <= n; i++ {
				if !m[i] {
					comp = append(comp, i)
				}
			}
			val := 0
			if len(comp) > 0 {
				val = query(comp, in, out)
			}
			for i := 0; i < k; i++ {
				if i == pos {
					ans[i] = val
				} else {
					ans[i] = globalMax
				}
			}
		}

		fmt.Fprint(out, "!")
		for i := 0; i < k; i++ {
			fmt.Fprintf(out, " %d", ans[i])
		}
		fmt.Fprintln(out)
		out.Flush()

		var verdict string
		fmt.Fscan(in, &verdict)
		if verdict != "Correct" {
			return
		}
	}
}
