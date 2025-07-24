package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, l, r int
	if _, err := fmt.Fscan(reader, &n, &l, &r); err != nil {
		return
	}
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	for i := 1; i < l; i++ {
		if a[i] != b[i] {
			fmt.Println("LIE")
			return
		}
	}
	for i := r + 1; i <= n; i++ {
		if a[i] != b[i] {
			fmt.Println("LIE")
			return
		}
	}
	freq := make([]int, n+1)
	for i := l; i <= r; i++ {
		freq[a[i]]++
		freq[b[i]]--
	}
	for i := 1; i <= n; i++ {
		if freq[i] != 0 {
			fmt.Println("LIE")
			return
		}
	}
	fmt.Println("TRUTH")
}
