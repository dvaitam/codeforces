package main

import (
	"bufio"
	"fmt"
	"os"
)

type Operation struct {
	typ  int
	v    int
	x    int64
	node int
}

type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) Add(idx int, val int64) {
	for i := idx; i <= f.n; i += i & -i {
		f.bit[i] += val
	}
}

func (f *Fenwick) Prefix(idx int) int64 {
	res := int64(0)
	for i := idx; i > 0; i -= i & -i {
		res += f.bit[i]
	}
	return res
}

func (f *Fenwick) RangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	f.Add(l, val)
	if r+1 <= f.n {
		f.Add(r+1, -val)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var q int
		fmt.Fscan(in, &q)

		ops := make([]Operation, q+1) // 1-indexed
		parent := make([]int, q+2)
		createTime := make([]int, q+2)
		nodeID := 1
		parent[1] = 0
		createTime[1] = 0

		for i := 1; i <= q; i++ {
			var t int
			fmt.Fscan(in, &t)
			if t == 1 {
				var v int
				fmt.Fscan(in, &v)
				nodeID++
				parent[nodeID] = v
				createTime[nodeID] = i
				ops[i] = Operation{typ: 1, v: v, node: nodeID}
			} else {
				var v int
				var x int64
				fmt.Fscan(in, &v, &x)
				ops[i] = Operation{typ: 2, v: v, x: x}
			}
		}
		n := nodeID
		children := make([][]int, n+1)
		for i := 2; i <= n; i++ {
			children[parent[i]] = append(children[parent[i]], i)
		}
		tin := make([]int, n+1)
		tout := make([]int, n+1)
		timer := 0
		var dfs func(int)
		dfs = func(v int) {
			timer++
			tin[v] = timer
			for _, to := range children[v] {
				dfs(to)
			}
			tout[v] = timer
		}
		dfs(1)

		bit := NewFenwick(n)
		initial := make([]int64, n+1)
		for i := 1; i <= q; i++ {
			op := ops[i]
			if op.typ == 1 {
				id := op.node
				initial[id] = bit.Prefix(tin[id])
			} else {
				l := tin[op.v]
				r := tout[op.v]
				bit.RangeAdd(l, r, op.x)
			}
		}

		ans := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			ans[i] = bit.Prefix(tin[i]) - initial[i]
		}
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
