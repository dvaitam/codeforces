package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the actual algorithm for Codeforces problem 1331F as
// described in problemF.txt. The statement is missing in this repository,
// so this placeholder just reads the input string and prints "NO" so the
// file compiles and the repository builds.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(reader, &s)
	fmt.Println("NO")
}
