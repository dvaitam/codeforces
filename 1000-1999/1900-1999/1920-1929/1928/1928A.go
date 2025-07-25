package main

import (
	"bufio"
	"fmt"
	"os"
)

func canTransform(a, b int64) bool {
	// ensure a <= b for simplified logic
	if a > b {
		a, b = b, a
	}

	// both odd cannot form a different rectangle
	if a%2 == 1 && b%2 == 1 {
		return false
	}

	// if one side is 1
	if a == 1 {
		if b%2 == 0 && b > 2 {
			return true
		}
		return false
	}

	// one side odd, other even
	if a%2 == 1 || b%2 == 1 {
		// let odd=a (since a<=b but if b odd, swap?). Actually we just check even == 2*odd.
		var odd, even int64
		if a%2 == 1 {
			odd, even = a, b
		} else {
			odd, even = b, a
		}
		if even == 2*odd {
			return false
		}
		return true
	}

	// both even
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		if canTransform(a, b) {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
