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
		var n int
		fmt.Fscan(reader, &n)
		q := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &q[i])
		}

		used := make([]bool, n+1)
		p := make([]int, n)
		possible := true

		if q[0] < 1 || q[0] > n {
			possible = false
		} else {
			p[0] = q[0]
			used[q[0]] = true
		}
		next := 1
		for next <= n && used[next] {
			next++
		}
		for i := 1; i < n && possible; i++ {
			if q[i] < i+1 || q[i] < q[i-1] || q[i] > n {
				possible = false
				break
			}
			if q[i] > q[i-1] {
				if used[q[i]] {
					possible = false
					break
				}
				p[i] = q[i]
				used[q[i]] = true
			} else { // q[i] == q[i-1]
				for next <= n && used[next] {
					next++
				}
				if next >= q[i] {
					possible = false
					break
				}
				p[i] = next
				used[next] = true
			}
		}

		if !possible {
			fmt.Fprintln(writer, -1)
		} else {
			for i := 0; i < n; i++ {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, p[i])
			}
			writer.WriteByte('\n')
		}
	}
}
