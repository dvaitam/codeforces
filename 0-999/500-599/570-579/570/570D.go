package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(reader, &n, &m)

	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(reader, &p)
		children[p] = append(children[p], i)
	}
	var letters string
	fmt.Fscan(reader, &letters)
	letters = " " + letters

	tin := make([]int, n+1)
	tout := make([]int, n+1)
	depth := make([]int, n+1)

	depthTins := make([][]int, n+2)
	depthPrefix := make([][]int, n+2)

	type item struct{ v, idx int }
	stack := []item{{1, 0}}
	depth[1] = 1
	time := 0
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		v := top.v
		if top.idx == 0 {
			time++
			tin[v] = time
			d := depth[v]
			if depthPrefix[d] == nil {
				depthPrefix[d] = []int{0}
			}
			bit := 1 << (letters[v] - 'a')
			depthTins[d] = append(depthTins[d], tin[v])
			pref := depthPrefix[d]
			depthPrefix[d] = append(pref, pref[len(pref)-1]^bit)
		}
		if top.idx < len(children[v]) {
			child := children[v][top.idx]
			top.idx++
			depth[child] = depth[v] + 1
			stack = append(stack, item{child, 0})
		} else {
			tout[v] = time
			stack = stack[:len(stack)-1]
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; m > 0; m-- {
		var v, h int
		fmt.Fscan(reader, &v, &h)
		if h >= len(depthTins) || len(depthTins[h]) == 0 {
			fmt.Fprintln(writer, "Yes")
			continue
		}
		tins := depthTins[h]
		l := sort.SearchInts(tins, tin[v])
		r := sort.Search(len(tins), func(i int) bool { return tins[i] > tout[v] })
		if l >= r {
			fmt.Fprintln(writer, "Yes")
			continue
		}
		pref := depthPrefix[h]
		mask := pref[r] ^ pref[l]
		if bits.OnesCount(uint(mask)) <= 1 {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
