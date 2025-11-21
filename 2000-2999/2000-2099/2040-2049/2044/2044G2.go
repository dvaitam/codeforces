package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        r := make([]int, n)
        radj := make([][]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &r[i])
            r[i]--
            radj[r[i]] = append(radj[r[i]], i)
        }

        visited := make([]bool, n)
        onStack := make([]bool, n)
        depth := make([]int, n)

        var dfs func(int)
        dfs = func(u int) {
            visited[u] = true
            onStack[u] = true
            for _, v := range radj[u] {
                if onStack[v] {
                    depth[u] = max(depth[u], 1)
                } else {
                    if !visited[v] {
                        dfs(v)
                    }
                    if depth[v] != 0 {
                        depth[u] = max(depth[u], depth[v]+1)
                    }
                }
            }
            onStack[u] = false
        }

        for i := 0; i < n; i++ {
            if !visited[i] {
                dfs(i)
            }
        }

        answer := 2
        for i := 0; i < n; i++ {
            if depth[i] > answer {
                answer = depth[i]
            }
        }
        fmt.Fprintln(out, answer)
    }
}
