package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	uniq := make(map[int]struct{})
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x > 0 {
			uniq[x] = struct{}{}
		}
	}
	fmt.Println(len(uniq))
}
