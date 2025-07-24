package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := range arr {
		fmt.Fscan(in, &arr[i])
	}
	// Placeholder solution. Proper computation of the expected value
	// is nontrivial and has not been implemented yet.
	fmt.Println(0)
}
