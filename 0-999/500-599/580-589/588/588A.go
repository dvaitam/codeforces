package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var minPrice int64 = 1 << 60
	var total int64
	for i := 0; i < n; i++ {
		var a, p int64
		fmt.Fscan(in, &a, &p)
		if p < minPrice {
			minPrice = p
		}
		total += a * minPrice
	}
	fmt.Println(total)
}
