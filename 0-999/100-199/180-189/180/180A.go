package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, m int
	a    [][]int
	b    []int
	c    []int
	a1   []int
	a2   []int
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &m)
	a = make([][]int, m+1)
	b = make([]int, n+1)
	c = make([]int, 0, n)
	a1 = make([]int, 0, 2*n)
	a2 = make([]int, 0, 2*n)
	for i := 1; i <= m; i++ {
		var cnt int
		fmt.Fscan(reader, &cnt)
		a[i] = make([]int, cnt)
		for j := 0; j < cnt; j++ {
			fmt.Fscan(reader, &a[i][j])
			b[a[i][j]] = i
		}
	}
	for i := 1; i <= n; i++ {
		if b[i] == 0 {
			c = append(c, i)
		}
	}
	k := 1
	for i := 1; i <= m; i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] == k {
				k++
				continue
			}
			if b[k] == 0 {
				old := a[i][j]
				b[k] = i
				b[old] = 0
				c = append(c, old)
				a1 = append(a1, old)
				a2 = append(a2, k)
				a[i][j] = k
			} else {
				old := a[i][j]
				x := c[len(c)-1]
				c = c[:len(c)-1]
				c = append(c, old)
				bi := b[k]
				for idx, val := range a[bi] {
					if val == k {
						a[bi][idx] = x
						break
					}
				}
				b[x] = b[k]
				b[k] = i
				b[old] = 0
				a1 = append(a1, k)
				a2 = append(a2, x)
				a1 = append(a1, old)
				a2 = append(a2, k)
				a[i][j] = k
			}
			k++
		}
	}
	fmt.Fprintln(writer, len(a1))
	for idx := range a1 {
		fmt.Fprintln(writer, a1[idx], a2[idx])
	}
}
