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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
		p[i]--
	}

	b := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
		sum += b[i]
	}

	visited := make([]bool, n)
	cycles := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			cycles++
			for j := i; !visited[j]; j = p[j] {
				visited[j] = true
			}
		}
	}

	ans := cycles
	if sum%2 == 0 {
		ans++
	}
	if cycles == 1 {
		ans--
	}

	fmt.Fprintln(writer, ans)
}
