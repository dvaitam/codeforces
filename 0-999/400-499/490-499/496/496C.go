package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	s := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}
	// sorted[i] indicates whether row i < row i+1 is already settled
	sorted := make([]bool, n-1)
	ans := 0
	for j := 0; j < m; j++ {
		// check if this column must be deleted
		del := false
		for i := 0; i < n-1; i++ {
			if !sorted[i] && s[i][j] > s[i+1][j] {
				del = true
				break
			}
		}
		if del {
			ans++
			continue
		}
		// mark settled pairs
		for i := 0; i < n-1; i++ {
			if !sorted[i] && s[i][j] < s[i+1][j] {
				sorted[i] = true
			}
		}
	}
	fmt.Println(ans)
}
