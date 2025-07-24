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
	sum := 0
	allZero := true
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		sum += a[i]
		if a[i] != 0 {
			allZero = false
		}
	}
	if allZero {
		fmt.Println("NO")
		return
	}
	fmt.Println("YES")
	if sum != 0 {
		fmt.Println(1)
		fmt.Printf("1 %d\n", n)
		return
	}
	prefix := 0
	split := -1
	for i := 0; i < n-1; i++ {
		prefix += a[i]
		if prefix != 0 {
			split = i + 1
			break
		}
	}
	if split == -1 {
		fmt.Println("NO")
		return
	}
	fmt.Println(2)
	fmt.Printf("1 %d\n", split)
	fmt.Printf("%d %d\n", split+1, n)
}
