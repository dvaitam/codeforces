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
	if n%2 == 1 {
		// Use one '7' for the extra segment
		fmt.Print("7")
		n -= 3
	}
	for i := 0; i < n/2; i++ {
		fmt.Print("1")
	}
}
