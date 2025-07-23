package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var h, k int
	if _, err := fmt.Fscan(reader, &n, &h, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	cur := 0
	ans := 0
	for _, x := range a {
		if cur+x <= h {
			cur += x
		} else {
			need := cur + x - h
			t := (need + k - 1) / k
			ans += t
			cur -= t * k
			if cur < 0 {
				cur = 0
			}
			cur += x
		}
	}
	if cur > 0 {
		ans += (cur + k - 1) / k
	}
	fmt.Println(ans)
}
