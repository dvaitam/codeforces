package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the actual algorithm for Codeforces problem 1331G as
// described in problemG.txt. The statement is missing in this repository,
// so this placeholder just reads all integers from standard input and
// prints 0 so the file compiles.
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			break
		}
	}
	fmt.Println(0)
}
