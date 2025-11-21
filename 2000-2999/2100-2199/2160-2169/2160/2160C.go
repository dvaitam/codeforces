package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func canObtain(n int) bool {
	if n == 0 {
		return true
	}
	k := bits.Len(uint(n))
	for p := 0; p < k; p++ {
		mask := (1 << uint(p)) - 1
		if n&mask != 0 {
			continue
		}
		m := n >> uint(p)
		length := k - p
		ok := true
		for i := 0; i < length/2; i++ {
			left := (m >> uint(i)) & 1
			right := (m >> uint(length-1-i)) & 1
			if left != right {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if length%2 == 1 {
			mid := (m >> uint(length/2)) & 1
			if mid == 1 {
				continue
			}
		}
		return true
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		if canObtain(n) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
