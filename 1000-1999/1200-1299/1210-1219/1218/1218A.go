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
	for i := 0; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	fmt.Println(int64(n) * int64(n+1) / 2)
}
