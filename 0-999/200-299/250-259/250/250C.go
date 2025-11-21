package main

import (
	"bufio"
	"fmt"
	"os"
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
	reader := bufio.NewReader(os.Stdin)
	n := nextInt(reader)
	k := nextInt(reader)

	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt(reader)
	}

	totalDiff := 0
	removed := make([]int, k+1)
	newEdges := make([]int, k+1)

	for i := 0; i+1 < n; i++ {
		if a[i] != a[i+1] {
			totalDiff++
			removed[a[i]]++
			removed[a[i+1]]++
		}
	}

	for i := 0; i < n; {
		val := a[i]
		j := i
		for j < n && a[j] == val {
			j++
		}
		if i > 0 && j < n && a[i-1] != a[j] {
			newEdges[val]++
		}
		i = j
	}

	bestGenre := 1
	bestStress := totalDiff - removed[1] + newEdges[1]

	for g := 2; g <= k; g++ {
		stress := totalDiff - removed[g] + newEdges[g]
		if stress < bestStress {
			bestStress = stress
			bestGenre = g
		}
	}

	fmt.Println(bestGenre)
}
