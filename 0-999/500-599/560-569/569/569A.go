package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var T, S, q int
	if _, err := fmt.Fscan(reader, &T, &S, &q); err != nil {
		return
	}
	cur := S
	count := 0
	for cur < T {
		cur *= q
		count++
	}
	fmt.Println(count)
}
