package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		c := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &c[i])
		}
		neigh := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			neigh[v] = append(neigh[v], u)
		}
		mp := make(map[string]int64)
		for v := 1; v <= n; v++ {
			if len(neigh[v]) == 0 {
				continue
			}
			sort.Ints(neigh[v])
			var sb strings.Builder
			for _, u := range neigh[v] {
				sb.WriteString(strconv.Itoa(u))
				sb.WriteByte(',')
			}
			key := sb.String()
			mp[key] += c[v]
		}
		var g int64
		for _, val := range mp {
			if g == 0 {
				g = val
			} else {
				g = gcd(g, val)
			}
		}
		fmt.Fprintln(writer, g)
	}
}
