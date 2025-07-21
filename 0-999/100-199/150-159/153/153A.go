package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b int
	// Read two integers (possibly with leading zeros) and sum them
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	fmt.Println(a + b)
}
