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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		board := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			board[i] = []byte(s)
		}
		total := n * m
		dist := make([]int, total)
		vis := make([]int, total)

		move := func(idx int) int {
			r := idx / m
			c := idx % m
			switch board[r][c] {
			case 'L':
				c--
			case 'R':
				c++
			case 'U':
				r--
			case 'D':
				r++
			}
			if r < 0 || r >= n || c < 0 || c >= m {
				return -1
			}
			return r*m + c
		}

		for i := 0; i < total; i++ {
			if dist[i] != 0 {
				continue
			}
			cur := i
			stack := []int{}
			for {
				if cur == -1 {
					base := 0
					for j := len(stack) - 1; j >= 0; j-- {
						base++
						dist[stack[j]] = base
					}
					break
				}
				if dist[cur] != 0 {
					base := dist[cur]
					for j := len(stack) - 1; j >= 0; j-- {
						base++
						dist[stack[j]] = base
					}
					break
				}
				if vis[cur] == i+1 {
					cycleLen := 1
					nxt := move(cur)
					for nxt != cur {
						cycleLen++
						nxt = move(nxt)
					}
					tmp := cur
					dist[tmp] = cycleLen
					nxt = move(tmp)
					for nxt != cur {
						dist[nxt] = cycleLen
						tmp = nxt
						nxt = move(tmp)
					}
					base := cycleLen
					for j := len(stack) - 1; j >= 0; j-- {
						node := stack[j]
						if dist[node] != 0 {
							base = dist[node]
						} else {
							base++
							dist[node] = base
						}
					}
					break
				}
				vis[cur] = i + 1
				stack = append(stack, cur)
				cur = move(cur)
			}
		}

		bestIdx := 0
		for i := 1; i < total; i++ {
			if dist[i] > dist[bestIdx] {
				bestIdx = i
			}
		}
		fmt.Fprintln(writer, bestIdx/m+1, bestIdx%m+1, dist[bestIdx])
	}
}
