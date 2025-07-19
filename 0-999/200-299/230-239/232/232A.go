package main

import "fmt"

func main() {
	var k int
	if _, err := fmt.Scan(&k); err != nil {
		return
	}
	var n int
	var adj [101][101]bool
	for n = 1; k > 0; n++ {
		m := n
		for m*(m-1)/2 > k {
			m--
		}
		for i := 0; i < m; i++ {
			adj[i][n] = true
			adj[n][i] = true
		}
		k -= m * (m - 1) / 2
	}
	// output
	fmt.Println(n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if adj[i][j] {
				fmt.Print(1)
			} else {
				fmt.Print(0)
			}
		}
		fmt.Println()
	}
}
