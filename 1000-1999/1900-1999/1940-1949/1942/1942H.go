package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	children []int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		if _, err := fmt.Fscan(in, &n, &q); err != nil {
			return
		}
		parent := make([]int, n)
		nodes := make([]Node, n)
		for i := 1; i < n; i++ {
			fmt.Fscan(in, &parent[i])
			parent[i]--
			nodes[parent[i]].children = append(nodes[parent[i]].children, i)
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		g := make([]int64, n)
		h := make([]int64, n)

		ops := make([][3]int64, q)
		for i := 0; i < q; i++ {
			var t, x, v int64
			fmt.Fscan(in, &t, &x, &v)
			ops[i] = [3]int64{t, x - 1, v}
		}

		for i := 0; i < q; i++ {
			op := ops[i]
			t := op[0]
			x := op[1]
			v := op[2]
			if t == 1 {
				g[x] += v
			} else {
				h[x] += v
			}
			if feasible(nodes, b, g, h) {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

func feasible(nodes []Node, b, g, h []int64) bool {
	var dfs func(int) (int64, int64)
	dfs = func(u int) (int64, int64) {
		need := int64(0)
		sumNeed := int64(0)
		childBonus := int64(0)
		for _, v := range nodes[u].children {
			n, bonus := dfs(v)
			sumNeed += n
			childBonus += bonus
		}
		supply := g[u]
		if supply >= sumNeed {
			supply -= sumNeed
		} else {
			need += sumNeed - supply
			supply = 0
		}
		required := b[u] + h[u]
		if childBonus >= required {
			childBonus -= required
		} else {
			required -= childBonus
			childBonus = 0
			if supply >= required {
				supply -= required
			} else {
				need += required - supply
				supply = 0
			}
		}
		bonus := supply
		return need, bonus
	}

	needRoot, _ := dfs(0)
	return needRoot == 0
}
