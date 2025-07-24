package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, b, d int
	if _, err := fmt.Fscan(reader, &n, &b, &d); err != nil {
		return
	}
	waste := 0
	empties := 0
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(reader, &a)
		if a > b {
			continue
		}
		waste += a
		if waste > d {
			empties++
			waste = 0
		}
	}
	fmt.Println(empties)
}
