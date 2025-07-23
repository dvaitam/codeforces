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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	i := 0
	for i+1 < n && a[i] < a[i+1] {
		i++
	}
	for i+1 < n && a[i] == a[i+1] {
		i++
	}
	for i+1 < n && a[i] > a[i+1] {
		i++
	}
	if i == n-1 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
