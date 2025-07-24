package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	id  int
	amt int64
}

type edge struct {
	u, v int
	w    int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)

	balance := make([]int64, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		var d int64
		fmt.Fscan(reader, &u, &v, &d)
		balance[u] -= d
		balance[v] += d
	}

	debtors := make([]pair, 0)
	creditors := make([]pair, 0)
	for i := 1; i <= n; i++ {
		if balance[i] < 0 {
			debtors = append(debtors, pair{i, -balance[i]})
		} else if balance[i] > 0 {
			creditors = append(creditors, pair{i, balance[i]})
		}
	}

	ans := make([]edge, 0)
	i, j := 0, 0
	for i < len(debtors) && j < len(creditors) {
		d := &debtors[i]
		c := &creditors[j]
		var take int64
		if d.amt < c.amt {
			take = d.amt
		} else {
			take = c.amt
		}
		ans = append(ans, edge{d.id, c.id, take})
		d.amt -= take
		c.amt -= take
		if d.amt == 0 {
			i++
		}
		if c.amt == 0 {
			j++
		}
	}

	fmt.Fprintln(writer, len(ans))
	for _, e := range ans {
		fmt.Fprintf(writer, "%d %d %d\n", e.u, e.v, e.w)
	}
}
