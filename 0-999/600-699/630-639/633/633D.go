package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	counts := make(map[int]int)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		counts[arr[i]]++
	}

	ans := 0
	if counts[0] > ans {
		ans = counts[0]
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			a, b := arr[i], arr[j]
			if a == 0 && b == 0 {
				continue
			}
			used := make(map[int]int)
			used[a]++
			used[b]++
			if used[a] > counts[a] || used[b] > counts[b] {
				continue
			}
			length := 2
			x, y := a, b
			for {
				next := x + y
				if counts[next]-used[next] <= 0 {
					break
				}
				used[next]++
				length++
				x, y = y, next
			}
			if length > ans {
				ans = length
			}
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, ans)
}
