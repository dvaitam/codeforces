package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(reader, &n, &k)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	sort.Ints(arr)
	var x int
	if k == 0 {
		x = arr[0] - 1
	} else {
		x = arr[k-1]
	}
	cnt := 0
	for _, v := range arr {
		if v <= x {
			cnt++
		}
	}
	if cnt != k || x < 1 || x > 1000000000 {
		fmt.Println(-1)
	} else {
		fmt.Println(x)
	}
}
