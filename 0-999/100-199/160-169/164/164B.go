package main

import (
	"bufio"
	"fmt"
	"os"
)

// fast integer reader
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
	la := nextInt(reader)
	lb := nextInt(reader)

	a := make([]int, la)
	for i := 0; i < la; i++ {
		a[i] = nextInt(reader)
	}

	index := make(map[int]int, lb)
	for i := 0; i < lb; i++ {
		v := nextInt(reader)
		index[v] = i
	}

	// map elements of a to positions in b (or -1 if absent)
	pos := make([]int, la)
	for i, v := range a {
		if p, ok := index[v]; ok {
			pos[i] = p
		} else {
			pos[i] = -1
		}
	}

	// handle rotation of a by doubling the array
	pos2 := make([]int, 2*la)
	copy(pos2, pos)
	copy(pos2[la:], pos)

	const INF int64 = -1
	arr := make([]int64, len(pos2)) // unwrapped positions in b
	var prev int64 = INF
	lb64 := int64(lb)
	for i, p := range pos2 {
		if p == -1 {
			arr[i] = INF
			prev = INF
			continue
		}
		cur := int64(p)
		if prev != INF && cur <= prev {
			cur += ((prev-cur)/lb64 + 1) * lb64
		}
		arr[i] = cur
		prev = cur
	}

	best := 0
	l := 0
	for r := 0; r < len(arr); r++ {
		if arr[r] == INF {
			l = r + 1
			continue
		}
		for l <= r && (arr[r]-arr[l] >= lb64 || r-l+1 > la) {
			l++
		}
		if r-l+1 > best {
			best = r - l + 1
		}
	}

	if best > la {
		best = la
	}
	if best > lb {
		best = lb
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, best)
	writer.Flush()
}
