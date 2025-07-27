package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func divisors(n int) []int {
	ds := []int{}
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			ds = append(ds, i)
			if i != n/i {
				ds = append(ds, n/i)
			}
		}
	}
	sort.Ints(ds)
	return ds
}

func solve(n int, p []int, c []int) int {
	visited := make([]bool, n)
	ans := n

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		var cycle []int
		j := i
		for !visited[j] {
			visited[j] = true
			cycle = append(cycle, j)
			j = p[j]
		}
		l := len(cycle)
		ds := divisors(l)
		for _, d := range ds {
			if d >= ans {
				continue
			}
			for start := 0; start < d; start++ {
				color := c[cycle[start]]
				good := true
				for pos := start; pos < l; pos += d {
					if c[cycle[pos]] != color {
						good = false
						break
					}
				}
				if good {
					if d < ans {
						ans = d
					}
					goto NextCycle
				}
			}
		}
	NextCycle:
	}

	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
			p[i]--
		}
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &c[i])
		}
		res := solve(n, p, c)
		fmt.Fprintln(writer, res)
	}
}
