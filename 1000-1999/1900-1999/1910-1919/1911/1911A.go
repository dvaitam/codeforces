package main

import (
	"bufio"
	"fmt"
	"os"
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
	freq := make(map[int]int)
	for _, v := range arr {
		freq[v]++
	}
	var unique int
	for v, c := range freq {
		if c == 1 {
			unique = v
			break
		}
	}
	for i, v := range arr {
		if v == unique {
			fmt.Println(i + 1)
			return
		}
	}
}
