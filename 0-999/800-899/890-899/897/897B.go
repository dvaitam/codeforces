package main

import (
	"bufio"
	"fmt"
	"os"
)

// makePal constructs an even-length palindromic number by
// appending the reverse of the digits of x to x itself.
func makePal(x int) int64 {
	res := int64(x)
	y := x
	for y > 0 {
		res = res*10 + int64(y%10)
		y /= 10
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k int
	var p int64
	if _, err := fmt.Fscan(in, &k, &p); err != nil {
		return
	}
	var sum int64
	for i := 1; i <= k; i++ {
		pal := makePal(i)
		sum = (sum + pal) % p
	}
	fmt.Fprintln(out, sum%p)
}
