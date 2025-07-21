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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int, n)
	copy(b, a)
	sort.Ints(b)
	diff := 0
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			diff++
		}
	}
	if diff == 0 || diff == 2 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
