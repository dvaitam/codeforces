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
	var s string
	fmt.Fscan(reader, &s)
	a := []byte(s)
	ans := 0
	colors := []byte{'R', 'G', 'B'}
	for i := 1; i < n; i++ {
		if a[i] == a[i-1] {
			ans++
			for _, c := range colors {
				if c != a[i-1] && (i == n-1 || c != a[i+1]) {
					a[i] = c
					break
				}
			}
		}
	}
	fmt.Println(ans)
	fmt.Println(string(a))
}
