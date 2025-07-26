package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b int
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	fmt.Println(a + b)
}
