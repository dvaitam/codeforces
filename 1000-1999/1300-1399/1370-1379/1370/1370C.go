package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func isPowerOfTwo(n int64) bool {
	return n&(n-1) == 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		if n == 1 {
			fmt.Fprintln(writer, "FastestFinger")
		} else if n == 2 {
			fmt.Fprintln(writer, "Ashishgup")
		} else if n%2 == 1 {
			fmt.Fprintln(writer, "Ashishgup")
		} else if isPowerOfTwo(n) {
			fmt.Fprintln(writer, "FastestFinger")
		} else if n%4 != 0 && isPrime(n/2) {
			fmt.Fprintln(writer, "FastestFinger")
		} else {
			fmt.Fprintln(writer, "Ashishgup")
		}
	}
}
