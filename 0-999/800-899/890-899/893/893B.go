package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)

	beautiful := []int{}
	for k := 0; ; k++ {
		val := (1<<(k+1) - 1) * (1 << k)
		if val > 100000 {
			break
		}
		beautiful = append(beautiful, val)
	}

	ans := 1
	for i := len(beautiful) - 1; i >= 0; i-- {
		if n%beautiful[i] == 0 {
			ans = beautiful[i]
			break
		}
	}
	fmt.Println(ans)
}
