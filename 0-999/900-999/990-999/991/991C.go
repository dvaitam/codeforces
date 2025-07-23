package main

import (
	"bufio"
	"fmt"
	"os"
)

func good(n, k int64) bool {
	candies := n
	eaten := int64(0)
	for candies > 0 {
		if candies < k {
			eaten += candies
			break
		}
		eaten += k
		candies -= k
		petya := candies / 10
		candies -= petya
	}
	return eaten*2 >= n
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	fmt.Fscan(in, &n)
	l, r := int64(1), n
	for l < r {
		mid := (l + r) / 2
		if good(n, mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	fmt.Println(l)
}
