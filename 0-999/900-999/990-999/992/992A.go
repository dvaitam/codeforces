package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	seen := make(map[int]bool)
	for i := 0; i < n; i++ {
		var val int
		fmt.Fscan(reader, &val)
		if val != 0 {
			seen[val] = true
		}
	}
	fmt.Println(len(seen))
}
