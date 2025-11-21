package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func nextInt(r *bufio.Reader) int {
	sign, val := 1, 0
	c, _ := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = r.ReadByte()
	}
	return sign * val
}

func main() {
	in := bufio.NewReader(os.Stdin)
	n := nextInt(in)
	weights := make([]int, n)
	for i := range weights {
		weights[i] = nextInt(in)
	}
	sort.Ints(weights)

	used := make(map[int]bool, n*2)
	result := 0
	for _, w := range weights {
		switch {
		case w > 1 && !used[w-1]:
			used[w-1] = true
			result++
		case !used[w]:
			used[w] = true
			result++
		case !used[w+1]:
			used[w+1] = true
			result++
		}
	}

	fmt.Println(result)
}
