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
	minDays := (n / 7) * 2
	maxDays := (n / 7) * 2
	rem := n % 7
	if rem > 2 {
		maxDays += 2
	} else {
		maxDays += rem
	}
	if rem > 5 {
		minDays += rem - 5
	}
	fmt.Printf("%d %d\n", minDays, maxDays)
}
