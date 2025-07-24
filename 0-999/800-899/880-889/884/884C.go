package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
		p[i]--
	}
	visited := make([]bool, n)
	var cycles []int
	for i := 0; i < n; i++ {
		if !visited[i] {
			cnt := 0
			j := i
			for !visited[j] {
				visited[j] = true
				j = p[j]
				cnt++
			}
			cycles = append(cycles, cnt)
		}
	}
	sort.Ints(cycles)
	if len(cycles) >= 2 {
		a := cycles[len(cycles)-1]
		b := cycles[len(cycles)-2]
		cycles = cycles[:len(cycles)-2]
		cycles = append(cycles, a+b)
	}
	var ans int64
	for _, c := range cycles {
		ans += int64(c) * int64(c)
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, ans)
}
