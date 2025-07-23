package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxVal = 200000

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	dist := make([]int, maxVal+1)
	cnt := make([]int, maxVal+1)

	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		visited := make(map[int]bool)

		val := v
		j := 0
		for val > 0 {
			x := val
			k := 0
			for x <= maxVal {
				if !visited[x] {
					dist[x] += j + k
					cnt[x]++
					visited[x] = true
				}
				x <<= 1
				k++
			}
			val >>= 1
			j++
		}
		// handle zero
		if !visited[0] {
			dist[0] += j
			cnt[0]++
		}
	}

	ans := int(^uint(0) >> 1)
	for x := 0; x <= maxVal; x++ {
		if cnt[x] == n && dist[x] < ans {
			ans = dist[x]
		}
	}
	fmt.Fprintln(writer, ans)
}
