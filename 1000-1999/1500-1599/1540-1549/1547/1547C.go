package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var k, n, m int
		fmt.Fscan(reader, &k, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}
		res := make([]int, 0, n+m)
		i, j := 0, 0
		lines := k
		possible := true
		for i < n || j < m {
			if i < n && a[i] == 0 {
				res = append(res, 0)
				lines++
				i++
				continue
			}
			if j < m && b[j] == 0 {
				res = append(res, 0)
				lines++
				j++
				continue
			}
			if i < n && a[i] <= lines {
				res = append(res, a[i])
				i++
				continue
			}
			if j < m && b[j] <= lines {
				res = append(res, b[j])
				j++
				continue
			}
			possible = false
			break
		}
		if !possible {
			fmt.Fprintln(writer, -1)
		} else {
			for idx, v := range res {
				if idx > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, v)
			}
			writer.WriteByte('\n')
		}
	}
}
