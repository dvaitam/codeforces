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
	var c, d int
	if _, err := fmt.Fscan(reader, &c, &d); err != nil {
		return
	}

	// iterate over Rick's scream times and check when Morty matches
	// upper bound large enough for given constraints
	for t := b; t <= 100000; t += a {
		if t >= d && (t-d)%c == 0 {
			fmt.Println(t)
			return
		}
	}
	fmt.Println(-1)
}
