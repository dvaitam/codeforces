package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	var x int
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &x)

	if (s == "ABC" && x < 2000) || (s == "ARC" && x < 2800) || (s == "AGC" && x >= 1200) {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
