package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, a, b int64
	if _, err := fmt.Fscan(reader, &n, &m, &a, &b); err != nil {
		return
	}
	remainder := n % m
	if remainder == 0 {
		fmt.Println(0)
		return
	}
	addCost := (m - remainder) * a
	removeCost := remainder * b
	if addCost < removeCost {
		fmt.Println(addCost)
	} else {
		fmt.Println(removeCost)
	}
}
