package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]uint64, 60)
	for i := 0; i < n; i++ {
		var x uint64
		fmt.Fscan(in, &x)
		b := bits.Len64(x) - 1
		g[b] = append(g[b], x)
	}

	res := make([]uint64, 0, n)
	for bit := 59; bit >= 0; bit-- {
		arr := g[bit]
		if len(arr) == 0 {
			continue
		}
		idx := 0
		newRes := make([]uint64, 0, len(res)+len(arr))
		var p uint64
		j := 0
		for j < len(res) {
			for idx < len(arr) && ((p>>uint(bit))&1) == 0 {
				newRes = append(newRes, arr[idx])
				p ^= arr[idx]
				idx++
			}
			newRes = append(newRes, res[j])
			p ^= res[j]
			j++
		}
		for idx < len(arr) && ((p>>uint(bit))&1) == 0 {
			newRes = append(newRes, arr[idx])
			p ^= arr[idx]
			idx++
		}
		if idx != len(arr) {
			fmt.Fprintln(out, "No")
			return
		}
		res = newRes
	}

	fmt.Fprintln(out, "Yes")
	for i := 0; i < len(res); i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	fmt.Fprintln(out)
}
