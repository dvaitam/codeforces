package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var m, n int
	if _, err := fmt.Fscan(reader, &m, &n); err != nil {
		return
	}
	// Maximum number of 2x1 dominoes that can be placed is floor(m*n/2)
	fmt.Println((m * n) / 2)
}
