package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, A, B int
	if _, err := fmt.Fscan(reader, &n, &A, &B); err != nil {
		return
	}
	s := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &s[i])
	}

	total := 0
	for _, v := range s {
		total += v
	}

	others := make([]int, n-1)
	copy(others, s[1:])
	sort.Slice(others, func(i, j int) bool { return others[i] > others[j] })

	S := int64(total)
	target := int64(s[0]) * int64(A)

	if target >= int64(B)*S {
		fmt.Println(0)
		return
	}

	blocked := 0
	for _, v := range others {
		S -= int64(v)
		blocked++
		if target >= int64(B)*S {
			fmt.Println(blocked)
			return
		}
	}

	fmt.Println(blocked)
}
