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

	var n, m, k, t int
	fmt.Fscan(in, &n, &m, &k, &t)

	type flight struct {
		u, v int
	}

	flights := make([]flight, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &flights[i].u, &flights[i].v)
	}

	degLeft := make([][]int, n)
	degRight := make([][]int, m)
	for i := range degLeft {
		degLeft[i] = make([]int, t)
	}
	for i := range degRight {
		degRight[i] = make([]int, t)
	}

	assignment := make([]int, k)
	for i := 0; i < k; i++ {
		u := flights[i].u - 1
		v := flights[i].v - 1

		bestCompany := 0
		bestDelta := 1 << 30
		for comp := 0; comp < t; comp++ {
			delta := (degLeft[u][comp]+1)*(degLeft[u][comp]+1) - degLeft[u][comp]*degLeft[u][comp]
			delta += (degRight[v][comp]+1)*(degRight[v][comp]+1) - degRight[v][comp]*degRight[v][comp]
			if delta < bestDelta {
				bestDelta = delta
				bestCompany = comp
			}
		}

		assignment[i] = bestCompany + 1
		degLeft[u][bestCompany]++
		degRight[v][bestCompany]++
	}

	unevenness := 0
	for i := 0; i < n; i++ {
		for comp := 0; comp < t; comp++ {
			unevenness += degLeft[i][comp] * degLeft[i][comp]
		}
	}
	for i := 0; i < m; i++ {
		for comp := 0; comp < t; comp++ {
			unevenness += degRight[i][comp] * degRight[i][comp]
		}
	}

	fmt.Fprintln(out, unevenness)
	for i, val := range assignment {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, val)
	}
	fmt.Fprintln(out)
}
