package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var h1, h2 int
	if _, err := fmt.Fscan(reader, &h1, &h2); err != nil {
		return
	}
	var a, b int
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}

	// Check if the caterpillar reaches the apple on the first day
	if h1+8*a >= h2 {
		fmt.Println(0)
		return
	}

	// If it can't make progress after the first day, it will never reach
	if a <= b {
		fmt.Println(-1)
		return
	}

	h := h1
	hours := 0
	for {
		// Evening: 2 pm to 10 pm (8 hours upward)
		hours += 8
		h += 8 * a
		if h >= h2 {
			fmt.Println(hours / 24)
			return
		}
		// Night: 10 pm to 10 am (12 hours downward)
		hours += 12
		h -= 12 * b
		// Morning: 10 am to 2 pm (4 hours upward)
		hours += 4
		h += 4 * a
		if h >= h2 {
			fmt.Println(hours / 24)
			return
		}
		// Continue with next day
	}
}
