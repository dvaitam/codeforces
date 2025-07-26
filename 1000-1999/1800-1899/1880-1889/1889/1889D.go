package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	stacks := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		var k int
		fmt.Fscan(in, &k)
		arr := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &arr[j])
		}
		stacks[i] = arr
	}

	visited := make([]int, n+1)
	used := make([]int, 0)

	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		pos := i
		used = used[:0]
		for {
			idx := len(stacks[pos]) - visited[pos]
			if idx <= 0 {
				break
			}
			visited[pos]++
			if visited[pos] == 1 {
				used = append(used, pos)
			}
			pos = stacks[pos][idx-1]
		}
		for _, v := range used {
			visited[v] = 0
		}
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, pos)
	}
	out.WriteByte('\n')
	out.Flush()
}
