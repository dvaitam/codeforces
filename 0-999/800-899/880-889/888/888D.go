package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	der := []int64{1, 0, 1, 2, 9}
	if k > 4 {
		k = 4
	}
	if n < k {
		k = n
	}

	var res int64
	comb := int64(1)
	for i := 0; i <= k; i++ {
		if i > 0 {
			comb = comb * int64(n-i+1) / int64(i)
		}
		res += comb * der[i]
	}

	fmt.Println(res)
}
