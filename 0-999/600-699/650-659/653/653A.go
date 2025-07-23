package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Ints(arr)
	uniq := make([]int, 0, n)
	for i, v := range arr {
		if i == 0 || v != arr[i-1] {
			uniq = append(uniq, v)
		}
	}
	ok := false
	for i := 0; i+2 < len(uniq); i++ {
		if uniq[i+2]-uniq[i] <= 2 {
			ok = true
			break
		}
	}
	if ok {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
