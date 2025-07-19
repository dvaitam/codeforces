package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for T > 0 {
		T--
		var n int
		fmt.Fscan(in, &n)
		b := make([]int, n+2)
		// v arrays of positions
		v := make([][]int, n+2)
		var k int
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &b[i])
			if b[i] > i {
				k++
			}
			if b[i] >= 0 && b[i] <= n+1 {
				v[b[i]] = append(v[b[i]], i)
			}
		}
		// build sequence
		a := make([]int, 0, n)
		cur := 0
		if len(v[n+1]) > 0 {
			cur = n + 1
		}
		cnt := 0
		for cnt < n {
			// count and prepare next cur
			cnt += len(v[cur])
			// find last element in v[cur] that has children
			last := len(v[cur]) - 1
			good := -1
			for j := 0; j <= last; j++ {
				nxt := v[cur][j]
				if nxt >= 0 && nxt < len(v) && len(v[nxt]) > 0 {
					good = j
				}
			}
			if good != -1 && good != last {
				v[cur][good], v[cur][last] = v[cur][last], v[cur][good]
			}
			// append all
			a = append(a, v[cur]...)
			// move cur to last element
			if len(v[cur]) > 0 {
				cur = v[cur][len(v[cur])-1]
			} else {
				// no more, break to avoid infinite
				break
			}
		}
		// output
		fmt.Fprintln(out, k)
		for i, x := range a {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, x)
		}
		fmt.Fprintln(out)
	}
}
