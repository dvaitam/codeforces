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
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
	}
	ans := int64(0)
	for length := 2; length <= n; length++ {
		l := int64(length - 1)
		ans += (l + 1) / 2
	}
	fmt.Println(ans)
}
