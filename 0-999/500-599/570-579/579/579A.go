package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var x int
	if _, err := fmt.Fscan(reader, &x); err != nil {
		return
	}
	count := 0
	for x > 0 {
		if x&1 == 1 {
			count++
		}
		x >>= 1
	}
	fmt.Println(count)
}
