package main

import (
	"bufio"
	"fmt"
	"os"
)

// xorTo returns XOR of all integers from 0 to n inclusive.
func xorTo(n uint64) uint64 {
	switch n & 3 {
	case 0:
		return n
	case 1:
		return 1
	case 2:
		return n + 1
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var x, m uint64
	var totalXor uint64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &x, &m)
		start := x
		end := x + m - 1
		// XOR of range [start, end] is xorTo(end) ^ xorTo(start-1)
		totalXor ^= xorTo(end) ^ xorTo(start-1)
	}
	if totalXor != 0 {
		fmt.Println("tolik")
	} else {
		fmt.Println("bolik")
	}
}
